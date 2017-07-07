package physical

import (
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestPhysicalView_impl(t *testing.T) {
	var _ Backend = new(View)
}

func newInmemTestBackend() *InmemBackend {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	return NewInmem(logger)
}

func TestPhysicalView_BadKeysKeys(t *testing.T) {
	backend := newInmemTestBackend()
	view := NewView(backend, "foo/")

	_, err := view.List("../")
	if err == nil {
		t.Fatalf("expected error")
	}

	_, err = view.Get("../")
	if err == nil {
		t.Fatalf("expected error")
	}

	err = view.Delete("../foo")
	if err == nil {
		t.Fatalf("expected error")
	}

	le := &Entry{
		Key:   "../foo",
		Value: []byte("test"),
	}
	err = view.Put(le)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPhysicalView(t *testing.T) {
	backend := newInmemTestBackend()
	view := NewView(backend, "foo/")

	// Write a key outside of foo/
	entry := &Entry{Key: "test", Value: []byte("test")}
	if err := backend.Put(entry); err != nil {
		t.Fatalf("bad: %v", err)
	}

	// List should have no visibility
	keys, err := view.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", err)
	}

	// Get should have no visibility
	out, err := view.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Try to put the same entry via the view
	if err := view.Put(entry); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check it is nested
	entry, err = backend.Get("foo/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("missing nested foo/test")
	}

	// Delete nested
	if err := view.Delete("test"); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the nested key
	entry, err = backend.Get("foo/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry != nil {
		t.Fatalf("nested foo/test should be gone")
	}

	// Check the non-nested key
	entry, err = backend.Get("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("root test missing")
	}
}
