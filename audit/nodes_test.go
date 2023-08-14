// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"

	"github.com/hashicorp/go-uuid"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/eventlogger"

	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/stretchr/testify/require"
)

// TestProcessManual_NilData tests ProcessManual when nil data is supplied.
func TestProcessManual_NilData(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	ids = append(ids, filterId)
	nodes[filterId] = filterNode

	// Sink node
	sinkId, sinkNode := newSinkNode(t)
	ids = append(ids, sinkId)
	nodes[sinkId] = sinkNode

	err := ProcessManual(namespace.RootContext(context.Background()), nil, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "data cannot be nil")
}

// TestProcessManual_NoIds tests ProcessManual when no IDs are supplied
func TestProcessManual_NoIds(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	nodes[filterId] = filterNode

	// Sink node
	sinkId, sinkNode := newSinkNode(t)
	nodes[sinkId] = sinkNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "ids are required")
}

// TestProcessManual_NoNodes tests ProcessManual when no nodes are supplied.
func TestProcessManual_NoNodes(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, _ := newFilterNode(t)
	ids = append(ids, filterId)

	// Sink node
	sinkId, _ := newSinkNode(t)
	ids = append(ids, sinkId)

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "nodes are required")
}

// TestProcessManual_IdNodeMismatch tests ProcessManual when IDs don't match with
// the nodes in the supplied map.
func TestProcessManual_IdNodeMismatch(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	ids = append(ids, filterId)
	nodes[filterId] = filterNode

	// Sink node
	sinkId, _ := newSinkNode(t)
	ids = append(ids, sinkId)

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.ErrorContains(t, err, "node not found: ")
}

// TestProcessManual_NotEnoughNodes tests ProcessManual when there is only one
// node provided.
func TestProcessManual_NotEnoughNodes(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	ids = append(ids, filterId)
	nodes[filterId] = filterNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "not enough nodes")
}

// TestProcessManual_LastNodeNotSink tests ProcessManual when the last node is
// not a Sink node.
func TestProcessManual_LastNodeNotSink(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	ids = append(ids, filterId)
	nodes[filterId] = filterNode

	formatterId, formatterNode := newFormatterNode(t)
	ids = append(ids, formatterId)
	nodes[formatterId] = formatterNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.Error(t, err)
	require.EqualError(t, err, "last node must be a sink")
}

// TestProcessManual ensures that the manual processing of a test message works
// as expected with proper inputs.
func TestProcessManual(t *testing.T) {
	t.Parallel()

	var ids []eventlogger.NodeID
	nodes := make(map[eventlogger.NodeID]eventlogger.Node)

	// Filter node
	filterId, filterNode := newFilterNode(t)
	ids = append(ids, filterId)
	nodes[filterId] = filterNode

	// Sink node
	sinkId, sinkNode := newSinkNode(t)
	ids = append(ids, sinkId)
	nodes[sinkId] = sinkNode

	// Data
	requestId, err := uuid.GenerateUUID()
	require.NoError(t, err)
	data := newData(requestId)

	err = ProcessManual(namespace.RootContext(context.Background()), data, ids, nodes)
	require.NoError(t, err)
}

// newFilterNode creates a new UUID and EntryFormatter (filter node).
func newFilterNode(t *testing.T) (eventlogger.NodeID, *EntryFormatter) {
	t.Helper()

	filterId, err := event.GenerateNodeID()
	require.NoError(t, err)
	cfg, err := NewFormatterConfig()
	require.NoError(t, err)
	filterNode, err := NewEntryFormatter(cfg, newStaticSalt(t))
	require.NoError(t, err)

	return filterId, filterNode
}

// newSinkNode creates a new UUID and NoopSink (sink node).
func newSinkNode(t *testing.T) (eventlogger.NodeID, *event.NoopSink) {
	t.Helper()

	sinkId, err := event.GenerateNodeID()
	require.NoError(t, err)
	sinkNode := event.NewNoopSink()

	return sinkId, sinkNode
}

// TestFormatter is a trivial implementation of the eventlogger.Node interface
// used as a place-holder for Formatter nodes in tests.
type TestFormatter struct{}

// Process trivially formats the event by storing "test" as a byte slice under
// the test format type.
func (f *TestFormatter) Process(_ context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	e.FormattedAs("test", []byte("test"))

	return e, nil
}

// Reopen does nothing.
func (f *TestFormatter) Reopen() error {
	return nil
}

// Type returns the eventlogger.NodeTypeFormatter type.
func (f *TestFormatter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFormatter
}

// newFormatterNode creates a new TestFormatter (formatter node).
func newFormatterNode(t *testing.T) (eventlogger.NodeID, *TestFormatter) {
	nodeId, err := event.GenerateNodeID()
	require.NoError(t, err)
	node := &TestFormatter{}

	return nodeId, node
}

// newData creates a sample logical.LogInput to be used as data for tests.
func newData(id string) *logical.LogInput {
	return &logical.LogInput{
		Type: "request",
		Auth: nil,
		Request: &logical.Request{
			ID:        id,
			Operation: "update",
			Path:      "sys/audit/test",
		},
		Response: nil,
		OuterErr: nil,
	}
}
