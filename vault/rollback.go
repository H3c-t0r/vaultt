package vault

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// rollbackPeriod is how often we attempt rollbacks for all the backends
	rollbackPeriod = time.Minute
)

// RollbackManager is responsible for performing rollbacks of partial
// secrets within logical backends.
//
// During normal operations, it is possible for logical backends to
// error partially through an operation. These are called "partial secrets":
// they are never sent back to a user, but they do need to be cleaned up.
// This manager handles that by periodically (on a timer) requesting that the
// backends clean up.
//
// The RollbackManager periodically initiates a logical.RollbackOperation
// on every mounted logical backend. It ensures that only one rollback operation
// is in-flight at any given time within a single seal/unseal phase.
type RollbackManager struct {
	logger log.Logger

	// This gives the current mount table of both logical and credential backends,
	// plus a RWMutex that is locked for reading. It is up to the caller to RUnlock
	// it when done with the mount table.
	backends func() []*MountEntry

	router *Router
	period time.Duration

	inflightAll  sync.WaitGroup
	inflight     map[string]*rollbackState
	inflightLock sync.RWMutex

	doneCh       chan struct{}
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
	quitContext  context.Context

	core *Core

	// First run worker pool
	firstRun bool
}

// rollbackState is used to track the state of a single rollback attempt
type rollbackState struct {
	lastError error
	sync.WaitGroup
	cancelLockGrabCtx       context.Context
	cancelLockGrabCtxCancel context.CancelFunc
}

// NewRollbackManager is used to create a new rollback manager
func NewRollbackManager(ctx context.Context, logger log.Logger, backendsFunc func() []*MountEntry, router *Router, core *Core) *RollbackManager {
	r := &RollbackManager{
		logger:      logger,
		backends:    backendsFunc,
		router:      router,
		period:      rollbackPeriod,
		inflight:    make(map[string]*rollbackState),
		doneCh:      make(chan struct{}),
		shutdownCh:  make(chan struct{}),
		quitContext: ctx,
		core:        core,
		firstRun:    true,
	}
	return r
}

// Start starts the rollback manager
func (m *RollbackManager) Start() {
	go m.run()
}

// Stop stops the running manager. This will wait for any in-flight
// rollbacks to complete.
func (m *RollbackManager) Stop() {
	m.shutdownLock.Lock()
	defer m.shutdownLock.Unlock()
	if !m.shutdown {
		m.shutdown = true
		close(m.shutdownCh)
		<-m.doneCh
	}
	m.inflightAll.Wait()
}

// run is a long running routine to periodically invoke rollback
func (m *RollbackManager) run() {
	m.logger.Info("starting rollback manager")
	tick := time.NewTicker(m.period)

	defer tick.Stop()
	defer close(m.doneCh)
	for {
		select {
		case <-tick.C:
			m.triggerRollbacks()

		case <-m.shutdownCh:
			m.logger.Info("stopping rollback manager")
			return
		}
	}
}

// triggerRollbacks is used to trigger the rollbacks across all the backends
func (m *RollbackManager) triggerRollbacks() {
	backends := m.backends()

	var jobs chan func()
	var wg sync.WaitGroup
	firstRun := m.firstRun
	m.firstRun = false
	if firstRun {
		// Use a small temporary worker pool to run the very first rollbacks in parallel, as they will trigger
		// backend initialization
		numWorkers := runtime.NumCPU() * 4 // For existing releases, don't modify the current behavior without an env override
		if v := os.Getenv("VAULT_INITIAL_ROLLBACK_CONCURRENCY"); v != "" {
			pv, err := strconv.Atoi(v)
			if err != nil || pv < 1 {
				m.logger.Warn("invalid value for VAULT_INITIAL_ROLLBACK_CONCURRENCY, must be a positive integer", "error", err, "value", pv)
			} else {
				numWorkers = pv
			}
		}
		jobs = make(chan func())
		// Start 'em
		for i := 0; i < numWorkers; i++ {
			go func() {
				for v := range jobs {
					v()
					wg.Done()
				}
			}()
		}

		// Randomize the jobs on first run to allow progress on migrations, etc to be made if there are restarts/reunseals
		// This is safe as the backendsFunc returns a fresh slice each time
		sort.Slice(backends, func(i, j int) bool {
			return strings.Compare(backends[i].UUID, backends[j].UUID) < 0
		})
	}

	midpoint := rand.Intn(len(backends))
	for _, e := range append(backends[midpoint:], backends[:midpoint]...) {
		path := e.Path
		if e.Table == credentialTableType {
			path = credentialRoutePrefix + path
		}

		// When the mount is filtered, the backend will be nil
		ctx := namespace.ContextWithNamespace(m.quitContext, e.namespace)
		backend := m.router.MatchingBackend(ctx, path)
		if backend == nil {
			continue
		}
		fullPath := e.namespace.Path + path

		// Start a rollback if necessary
		if firstRun {
			wg.Add(1)
			m.startOrLookupRollback(ctx, fullPath, true, jobs)
		} else {
			m.startOrLookupRollback(ctx, fullPath, true, nil)
		}
	}
	if firstRun {
		close(jobs)
		wg.Wait()
	}
}

// startOrLookupRollback is used to start an async rollback attempt.
// This must be called with the inflightLock held.
func (m *RollbackManager) startOrLookupRollback(ctx context.Context, fullPath string, grabStatelock bool, workerChan chan func()) *rollbackState {
	m.inflightLock.Lock()
	rsInflight, ok := m.inflight[fullPath]
	if ok {
		m.inflightLock.Unlock()
		return rsInflight
	}
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	rs := &rollbackState{
		cancelLockGrabCtx:       cancelCtx,
		cancelLockGrabCtxCancel: cancelFunc,
	}

	// If no inflight rollback is already running, kick one off
	m.inflight[fullPath] = rs
	rs.Add(1)
	m.inflightAll.Add(1)

	if workerChan != nil {
		m.inflightLock.Unlock()
		workerChan <- func() {
			m.attemptRollback(ctx, fullPath, rs, grabStatelock)
		}
	} else {
		defer m.inflightLock.Unlock()
		go m.attemptRollback(ctx, fullPath, rs, grabStatelock)
	}
	return rs
}

// attemptRollback invokes a RollbackOperation for the given path
func (m *RollbackManager) attemptRollback(ctx context.Context, fullPath string, rs *rollbackState, grabStatelock bool) (err error) {
	defer metrics.MeasureSince([]string{"rollback", "attempt", strings.Replace(fullPath, "/", "-", -1)}, time.Now())

	defer func() {
		rs.lastError = err
		rs.Done()
		m.inflightAll.Done()
		m.inflightLock.Lock()
		delete(m.inflight, fullPath)
		m.inflightLock.Unlock()
	}()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		m.logger.Error("rollback failed to derive namespace from context", "path", fullPath)
		return err
	}
	if ns == nil {
		m.logger.Error("rollback found no namespace", "path", fullPath)
		return namespace.ErrNoNamespace
	}

	// Invoke a RollbackOperation
	req := &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      ns.TrimmedPath(fullPath),
	}

	releaseLock := true
	if grabStatelock {
		doneCh := make(chan struct{})
		defer close(doneCh)

		stopCh := make(chan struct{})
		go func() {
			defer close(stopCh)

			select {
			case <-m.shutdownCh:
			case <-rs.cancelLockGrabCtx.Done():
			case <-doneCh:
			}
		}()

		// Grab the statelock or stop
		if stopped := grabLockOrStop(m.core.stateLock.RLock, m.core.stateLock.RUnlock, stopCh); stopped {
			// If we stopped due to shutdown, return. Otherwise another thread
			// is holding the lock for us, continue on.
			select {
			case <-m.shutdownCh:
				return errors.New("rollback shutting down")
			default:
				releaseLock = false
			}
		}
	}

	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithTimeout(ctx, DefaultMaxRequestDuration)
	resp, err := m.router.Route(ctx, req)
	if grabStatelock && releaseLock {
		m.core.stateLock.RUnlock()
	}
	cancelFunc()

	// If the error is an unsupported operation, then it doesn't
	// matter, the backend doesn't support it.
	if err == logical.ErrUnsupportedOperation {
		err = nil
	}
	// If we failed due to read-only storage, we can't do anything; ignore
	if (err != nil && strings.Contains(err.Error(), logical.ErrReadOnly.Error())) ||
		(resp.IsError() && strings.Contains(resp.Error().Error(), logical.ErrReadOnly.Error())) {
		err = nil
	}
	if err != nil {
		m.logger.Error("error rolling back", "path", fullPath, "error", err)
	}
	return
}

// Rollback is used to trigger an immediate rollback of the path,
// or to join an existing rollback operation if in flight. Caller should have
// core's statelock held (write OR read). If an already inflight rollback is
// happening this function will simply wait for it to complete
func (m *RollbackManager) Rollback(ctx context.Context, path string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}
	fullPath := ns.Path + path

	// Check for an existing attempt or start one if none
	rs := m.startOrLookupRollback(ctx, fullPath, false, nil)

	// Since we have the statelock held, tell any inflight rollback to give up
	// trying to acquire it. This will prevent deadlocks in the case where we
	// have the write lock. In the case where it was waiting to grab
	// a read lock it will then simply continue with the rollback
	// operation under the protection of our write lock.
	rs.cancelLockGrabCtxCancel()

	// It's safe to do this, since the other thread either already has the lock
	// held, or we just canceled it above.
	rs.Wait()

	// Return the last error
	return rs.lastError
}

// The methods below are the hooks from core that are called pre/post seal.

// startRollback is used to start the rollback manager after unsealing
func (c *Core) startRollback() error {
	backendsFunc := func() []*MountEntry {
		ret := []*MountEntry{}
		c.mountsLock.RLock()
		defer c.mountsLock.RUnlock()
		// During teardown/setup after a leader change or unseal there could be
		// something racy here so make sure the table isn't nil
		if c.mounts != nil {
			for _, entry := range c.mounts.Entries {
				ret = append(ret, entry)
			}
		}
		c.authLock.RLock()
		defer c.authLock.RUnlock()
		// During teardown/setup after a leader change or unseal there could be
		// something racy here so make sure the table isn't nil
		if c.auth != nil {
			for _, entry := range c.auth.Entries {
				ret = append(ret, entry)
			}
		}
		return ret
	}
	rollbackLogger := c.baseLogger.Named("rollback")
	c.AddLogger(rollbackLogger)
	c.rollback = NewRollbackManager(c.activeContext, rollbackLogger, backendsFunc, c.router, c)
	c.rollback.Start()
	return nil
}

// stopRollback is used to stop running the rollback manager before sealing
func (c *Core) stopRollback() error {
	if c.rollback != nil {
		c.rollback.Stop()
		c.rollback = nil
	}
	return nil
}
