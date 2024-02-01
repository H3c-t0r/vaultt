// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/helper/backoff"
	"github.com/hashicorp/vault/sdk/logical"
)

type Factory func(context.Context) (SubscriptionPlugin, error)

// SubscriptionPlugin is the interface implemented by plugins that can subscribe to and receive events.
type SubscriptionPlugin interface {
	// Send is used to set up a new connection and send events to that connection.
	// The first call should have .Subscribe populated with the configuration.
	// Other calls should populate .Event with the event being sent.
	Send(ctx context.Context, request *Request) error
	// PluginName returns the name for the particular event subscription plugin.
	// This type name is usually set as a constant the backend, e.g., "sqs" for the
	// AWS SQS backend.
	PluginName() string
	PluginVersion() logical.PluginVersion
	Close(ctx context.Context) error
}

type Request struct {
	Subscribe   *SubscribeRequest
	Unsubscribe *UnsubscribeRequest
	Event       *SendEventRequest
}

type SubscribeRequest struct {
	SubscriptionID   string
	Config           map[string]interface{}
	VerifyConnection bool
}

type UnsubscribeRequest struct {
	SubscriptionID string
}

type SendEventRequest struct {
	SubscriptionID string
	EventJSON      string
}

// SubscribeConfigDefaults defines configuration map keys for common default options.
// Embed this in your own config struct to pick up these default options.
type SubscribeConfigDefaults struct {
	Retries         *int           `mapstructure:"retries"`
	RetryMinBackoff *time.Duration `mapstructure:"retry_min_backoff"`
	RetryMaxBackoff *time.Duration `mapstructure:"retry_max_backoff"`
}

// default values for common configuration keys
const (
	DefaultRetries         = 3
	DefaultRetryMinBackoff = 100 * time.Millisecond
	DefaultRetryMaxBackoff = 5 * time.Second
)

func (c *SubscribeConfigDefaults) GetRetries() int {
	if c.Retries == nil {
		return DefaultRetries
	}
	return *c.Retries
}

func (c *SubscribeConfigDefaults) GetRetryMinBackoff() time.Duration {
	if c.RetryMinBackoff == nil {
		return DefaultRetryMinBackoff
	}
	return *c.RetryMinBackoff
}

func (c *SubscribeConfigDefaults) GetRetryMaxBackoff() time.Duration {
	if c.RetryMaxBackoff == nil {
		return DefaultRetryMaxBackoff
	}
	return *c.RetryMaxBackoff
}

func (c *SubscribeConfigDefaults) NewRetryBackoff() *backoff.Backoff {
	return backoff.NewBackoff(c.GetRetries(), c.GetRetryMinBackoff(), c.GetRetryMaxBackoff())
}
