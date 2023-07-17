package audit

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

// getDefaultOptions returns AuditOptions with their default values.
func getDefaultOptions() AuditOptions {
	return AuditOptions{
		Options: event.Options{
			WithNow: time.Now(),
		},
		withFacility:    "AUTH",
		withTag:         "vault",
		withSocketType:  "tcp",
		withMaxDuration: 2 * time.Second,
	}
}

// getOpts applies all the supplied AuditOption and returns configured AuditOptions.
// Each AuditOption is applied in the order it appears in the argument list, so it is
// possible to supply the same AuditOption numerous times and the 'last write wins'.
func getOpts(opt ...AuditOption) (AuditOptions, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return AuditOptions{}, err
		}
	}
	return opts, nil
}

// WithID provides an optional ID.
func WithID(id string) AuditOption {
	return func(o *AuditOptions) error {
		var err error

		id := strings.TrimSpace(id)
		switch {
		case id == "":
			err = errors.New("id cannot be empty")
		default:
			o.WithID = id
		}

		return err
	}
}

// WithNow provides an option to represent 'now'.
func WithNow(now time.Time) AuditOption {
	return func(o *AuditOptions) error {
		var err error

		switch {
		case now.IsZero():
			err = errors.New("cannot specify 'now' to be the zero time instant")
		default:
			o.WithNow = now
		}

		return err
	}
}

// WithSubtype provides an option to represent the subtype.
func WithSubtype(subtype string) AuditOption {
	return func(o *AuditOptions) error {
		s := strings.TrimSpace(subtype)
		if s == "" {
			return errors.New("subtype cannot be empty")
		}

		parsed := auditSubtype(s)
		err := parsed.validate()
		if err != nil {
			return err
		}

		o.withSubtype = parsed
		return nil
	}
}

// WithFormat provides an option to represent event format.
func WithFormat(format string) AuditOption {
	return func(o *AuditOptions) error {
		f := strings.TrimSpace(format)
		if f == "" {
			return errors.New("format cannot be empty")
		}

		parsed := auditFormat(f)
		err := parsed.validate()
		if err != nil {
			return err
		}

		o.withFormat = parsed
		return nil
	}
}

// WithFileMode provides an option to represent a file mode for a file sink.
// Supplying an empty string or whitespace will prevent this option from being
// applied, but it will not return an error in those circumstances.
func WithFileMode(mode string) AuditOption {
	return func(o *AuditOptions) error {
		// If supplied file mode is empty, just return early without setting anything.
		// We can assume that this option was called by something that didn't
		// parse the incoming value, perhaps from a config map etc.
		mode = strings.TrimSpace(mode)
		if mode == "" {
			return nil
		}

		// By now we believe we have something that the caller really intended to
		// be parsed into a file mode.
		raw, err := strconv.ParseUint(mode, 8, 32)

		switch {
		case err != nil:
			return fmt.Errorf("unable to parse file mode: %w", err)
		default:
			m := os.FileMode(raw)
			o.withFileMode = &m
		}

		return nil
	}
}

// WithPrefix provides an option to represent a prefix for a file sink.
func WithPrefix(prefix string) AuditOption {
	return func(o *AuditOptions) error {
		o.withPrefix = prefix
		return nil
	}
}

// WithFacility provides an option to represent a 'facility' for a syslog sink.
func WithFacility(facility string) AuditOption {
	return func(o *AuditOptions) error {
		facility = strings.TrimSpace(facility)

		if facility != "" {
			o.withFacility = facility
		}

		return nil
	}
}

// WithTag provides an option to represent a 'tag' for a syslog sink.
func WithTag(tag string) AuditOption {
	return func(o *AuditOptions) error {
		tag = strings.TrimSpace(tag)

		if tag != "" {
			o.withTag = tag
		}

		return nil
	}
}

// WithSocketType provides an option to represent the socket type for a socket sink.
func WithSocketType(socketType string) AuditOption {
	return func(o *AuditOptions) error {
		socketType = strings.TrimSpace(socketType)

		if socketType != "" {
			o.withSocketType = socketType
		}

		return nil
	}
}

// WithMaxDuration provides an option to represent the max duration for writing to a socket sink.
func WithMaxDuration(duration string) AuditOption {
	return func(o *AuditOptions) error {
		duration = strings.TrimSpace(duration)

		if duration == "" {
			return nil
		}

		parsed, err := parseutil.ParseDurationSecond(duration)
		if err != nil {
			return err
		}

		o.withMaxDuration = parsed

		return nil
	}
}
