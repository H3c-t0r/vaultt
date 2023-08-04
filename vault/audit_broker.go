// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/eventlogger"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/logical"
)

type backendEntry struct {
	backend audit.Backend
	local   bool
}

// AuditBroker is used to provide a single ingest interface to auditable
// events given that multiple backends may be configured.
type AuditBroker struct {
	sync.RWMutex
	backends map[string]backendEntry
	logger   log.Logger

	broker *eventlogger.Broker
}

// NewAuditBroker creates a new audit broker
func NewAuditBroker(log log.Logger, useEventLogger bool) (*AuditBroker, error) {
	var eventBroker *eventlogger.Broker
	var err error

	if useEventLogger {
		// Ignoring the second error return value since an error will only occur
		// if an unrecognized eventlogger.RegistrationPolicy is provided to an
		// eventlogger.Option function.
		eventBroker, err = eventlogger.NewBroker(eventlogger.WithNodeRegistrationPolicy(eventlogger.DenyOverwrite), eventlogger.WithPipelineRegistrationPolicy(eventlogger.DenyOverwrite))
		if err != nil {
			return nil, fmt.Errorf("error creating event broker for audit events: %w", err)
		}
	}

	b := &AuditBroker{
		backends: make(map[string]backendEntry),
		logger:   log,
		broker:   eventBroker,
	}
	return b, nil
}

// Register is used to add new audit backend to the broker
func (a *AuditBroker) Register(name string, b audit.Backend, local bool) error {
	a.Lock()
	defer a.Unlock()

	if a.broker != nil {
		err := b.RegisterNodesAndPipeline(a.broker, name)
		if err != nil {
			return err
		}
	}

	a.backends[name] = backendEntry{
		backend: b,
		local:   local,
	}

	return nil
}

// Deregister is used to remove an audit backend from the broker
func (a *AuditBroker) Deregister(ctx context.Context, name string) error {
	a.Lock()
	defer a.Unlock()

	// Remove the Backend from the map first, so that if an error occurs while
	// removing the pipeline and nodes, we can quickly exit this method with
	// the error.
	delete(a.backends, name)

	if a.broker != nil {
		// The first return value, a bool, indicates whether
		// RemovePipelineAndNodes encountered the error while evaluating
		// pre-conditions (false) or once it started removing the pipeline and
		// the nodes (true). This code doesn't care either way.
		_, err := a.broker.RemovePipelineAndNodes(ctx, eventlogger.EventType("audit"), eventlogger.PipelineID(name))
		if err != nil {
			return err
		}
	}

	return nil
}

// IsRegistered is used to check if a given audit backend is registered
func (a *AuditBroker) IsRegistered(name string) bool {
	a.RLock()
	defer a.RUnlock()

	_, ok := a.backends[name]
	return ok
}

// IsLocal is used to check if a given audit backend is registered
func (a *AuditBroker) IsLocal(name string) (bool, error) {
	a.RLock()
	defer a.RUnlock()
	be, ok := a.backends[name]
	if ok {
		return be.local, nil
	}
	return false, fmt.Errorf("unknown audit backend %q", name)
}

// GetHash returns a hash using the salt of the given backend
func (a *AuditBroker) GetHash(ctx context.Context, name string, input string) (string, error) {
	a.RLock()
	defer a.RUnlock()
	be, ok := a.backends[name]
	if !ok {
		return "", fmt.Errorf("unknown audit backend %q", name)
	}

	return audit.HashString(ctx, be.backend, input)
}

// LogRequest is used to ensure all the audit backends have an opportunity to
// log the given request and that *at least one* succeeds.
func (a *AuditBroker) LogRequest(ctx context.Context, in *logical.LogInput, headersConfig *AuditedHeadersConfig) (ret error) {
	defer metrics.MeasureSince([]string{"audit", "log_request"}, time.Now())

	a.RLock()
	defer a.RUnlock()

	if in.Request.InboundSSCToken != "" {
		if in.Auth != nil {
			reqAuthToken := in.Auth.ClientToken
			in.Auth.ClientToken = in.Request.InboundSSCToken
			defer func() {
				in.Auth.ClientToken = reqAuthToken
			}()
		}
	}

	var retErr *multierror.Error

	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("panic during logging", "request_path", in.Request.Path, "error", r, "stacktrace", string(debug.Stack()))
			retErr = multierror.Append(retErr, fmt.Errorf("panic generating audit log"))
		}

		ret = retErr.ErrorOrNil()
		failure := float32(0.0)
		if ret != nil {
			failure = 1.0
		}
		metrics.IncrCounter([]string{"audit", "log_request_failure"}, failure)
	}()

	// All logged requests must have an identifier
	//if req.ID == "" {
	//	a.logger.Error("missing identifier in request object", "request_path", req.Path)
	//	retErr = multierror.Append(retErr, fmt.Errorf("missing identifier in request object: %s", req.Path))
	//	return
	//}

	headers := in.Request.Headers
	defer func() {
		in.Request.Headers = headers
	}()

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		in.Request.Headers = nil
		transHeaders, thErr := headersConfig.ApplyConfig(ctx, headers, be.backend)
		if thErr != nil {
			a.logger.Error("backend failed to include headers", "backend", name, "error", thErr)
			continue
		}
		in.Request.Headers = transHeaders

		start := time.Now()
		lrErr := be.backend.LogRequest(ctx, in)
		metrics.MeasureSinceWithLabels([]string{"audit", "device", "log_request"}, start, []metrics.Label{
			{"device_name", name},
		})
		if lrErr != nil {
			a.logger.Error("backend failed to log request", "backend", name, "error", lrErr)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		retErr = multierror.Append(retErr, fmt.Errorf("no audit backend succeeded in logging the request"))
	}

	return retErr.ErrorOrNil()
}

// LogResponse is used to ensure all the audit backends have an opportunity to
// log the given response and that *at least one* succeeds.
func (a *AuditBroker) LogResponse(ctx context.Context, in *logical.LogInput, headersConfig *AuditedHeadersConfig) (ret error) {
	defer metrics.MeasureSince([]string{"audit", "log_response"}, time.Now())
	a.RLock()
	defer a.RUnlock()
	if in.Request.InboundSSCToken != "" {
		if in.Auth != nil {
			reqAuthToken := in.Auth.ClientToken
			in.Auth.ClientToken = in.Request.InboundSSCToken
			defer func() {
				in.Auth.ClientToken = reqAuthToken
			}()
		}
	}

	var retErr *multierror.Error

	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("panic during logging", "request_path", in.Request.Path, "error", r, "stacktrace", string(debug.Stack()))
			retErr = multierror.Append(retErr, fmt.Errorf("panic generating audit log"))
		}

		ret = retErr.ErrorOrNil()

		failure := float32(0.0)
		if ret != nil {
			failure = 1.0
		}
		metrics.IncrCounter([]string{"audit", "log_response_failure"}, failure)
	}()

	headers := in.Request.Headers
	defer func() {
		in.Request.Headers = headers
	}()

	// Ensure at least one backend logs
	anyLogged := false
	for name, be := range a.backends {
		in.Request.Headers = nil
		transHeaders, thErr := headersConfig.ApplyConfig(ctx, headers, be.backend)
		if thErr != nil {
			a.logger.Error("backend failed to include headers", "backend", name, "error", thErr)
			continue
		}
		in.Request.Headers = transHeaders

		start := time.Now()
		lrErr := be.backend.LogResponse(ctx, in)
		metrics.MeasureSinceWithLabels([]string{"audit", "device", "log_response"}, start, []metrics.Label{
			{"device_name", name},
		})
		if lrErr != nil {
			a.logger.Error("backend failed to log response", "backend", name, "error", lrErr)
		} else {
			anyLogged = true
		}
	}
	if !anyLogged && len(a.backends) > 0 {
		retErr = multierror.Append(retErr, fmt.Errorf("no audit backend succeeded in logging the response"))
	}

	return retErr.ErrorOrNil()
}

func (a *AuditBroker) Invalidate(ctx context.Context, key string) {
	// For now we ignore the key as this would only apply to salts. We just
	// sort of brute force it on each one.
	a.Lock()
	defer a.Unlock()
	for _, be := range a.backends {
		be.backend.Invalidate(ctx)
	}
}
