package monitor

import (
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/hashicorp/go-hclog"
)

// Monitor provides a mechanism to stream logs using go-hclog
// InterceptLogger and SinkAdapter. It allows streaming of logs
// at a different log level than what is set on the logger.
type Monitor interface {
	// Start returns a channel of log messages which are sent
	// ever time a log message occurs
	Start() <-chan []byte

	// Stop de-registers the sink from the InterceptLogger
	// and closes the log channels
	Stop()
}

// monitor implements the Monitor interface
type monitor struct {
	sink log.SinkAdapter

	// logger is the logger we will be monitoring
	logger log.InterceptLogger

	// logCh is a buffered chan where we send logs when streaming
	logCh chan []byte

	// doneCh coordinates the shutdown of logCh
	doneCh chan struct{}

	// droppedCount is the current count of messages
	// that were dropped from the logCh buffer.
	droppedCount uint64
	bufSize      int
	// dropCheckInterval is the amount of time we should
	// wait to check for dropped messages. Defaults
	// to 3 seconds
	dropCheckInterval time.Duration
}

// New creates a new Monitor. Start must be called in order to actually start
// streaming logs
func New(buf int, logger log.InterceptLogger, opts *log.LoggerOptions) Monitor {
	return new(buf, logger, opts)
}

func new(buf int, logger log.InterceptLogger, opts *log.LoggerOptions) *monitor {
	sw := &monitor{
		logger:          logger,
		logCh:           make(chan []byte, buf),
		doneCh:          make(chan struct{}, 1),
		bufSize:         buf,
		dropCheckInterval: 3 * time.Second,
	}

	opts.Output = sw
	sink := log.NewSinkAdapter(opts)
	sw.sink = sink

	return sw
}

// Stop deregisters the sink and stops the monitoring process
func (d *monitor) Stop() {
	d.logger.DeregisterSink(d.sink)
	close(d.doneCh)
}

// Start registers a sink on the monitor's logger and starts sending
// received log messages over the returned channel.
func (d *monitor) Start() <-chan []byte {
	// register our sink with the logger
	d.logger.RegisterSink(d.sink)

	streamCh := make(chan []byte, d.bufSize)

	// run a go routine that listens for streamed
	// log messages and sends them to streamCh
	go func() {
		defer close(streamCh)

		for {
			select {
			case log := <-d.logCh:
				select {
				case <-d.doneCh:
					return
				case streamCh <- log:
				}
			case <-d.doneCh:
				return
			}
		}
	}()

	// run a go routine that periodically checks for
	// dropped messages and makes room on the logCh
	// to add a dropped message count warning
	go func() {
		// loop and check for dropped messages
		for {
			select {
			case <-d.doneCh:
				return
			case <-time.After(d.dropCheckInterval):
				// Check if there have been any dropped messages.
				dc := atomic.LoadUint64(&d.droppedCount)

				if dc > 0 {
					dropped := fmt.Sprintf("[WARN] Monitor dropped %d logs during monitor request\n", dc)
					select {
					case <-d.doneCh:
						return
					// Try sending dropped message count to logCh in case
					// there is room in the buffer now.
					case d.logCh <- []byte(dropped):
					default:
						// Drop a log message to make room for "Monitor dropped.." message
						select {
						case <-d.logCh:
							atomic.AddUint64(&d.droppedCount, 1)
							dc := atomic.LoadUint64(&d.droppedCount)
							dropped = fmt.Sprintf("[WARN] Monitor dropped %d logs during monitor request\n", dc)
						default:
						}
						d.logCh <- []byte(dropped)
					}
					atomic.SwapUint64(&d.droppedCount, 0)
				}
			}
		}
	}()

	return streamCh
}

// Write attempts to send latest log to logCh
// it drops the log if channel is unavailable to receive
func (d *monitor) Write(p []byte) (n int, err error) {
	// ensure logCh is still open
	select {
	case <-d.doneCh:
		return
	default:
	}

	bytes := make([]byte, len(p))
	copy(bytes, p)

	select {
	case d.logCh <- bytes:
	default:
		atomic.AddUint64(&d.droppedCount, 1)
	}

	return len(p), nil
}
