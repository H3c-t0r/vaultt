package logical

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

// TestRequest is a helper to create a purely in-memory Request struct.
func TestRequest(t *testing.T, op Operation, path string) *Request {
	return &Request{
		Operation: op,
		Path:      path,
		Data:      make(map[string]interface{}),
		Storage:   new(InmemStorage),
	}
}

// TestStorage is a helper that can be used from unit tests to verify
// the behavior of a Storage impl.
func TestStorage(t *testing.T, s Storage) {
	keys, err := s.List("")
	if err != nil {
		t.Fatalf("list error: %s", err)
	}
	if len(keys) > 0 {
		t.Fatalf("should have no keys to start: %#v", keys)
	}

	entry := &StorageEntry{Key: "foo", Value: []byte("bar")}
	if err := s.Put(entry); err != nil {
		t.Fatalf("put error: %s", err)
	}

	actual, err := s.Get("foo")
	if err != nil {
		t.Fatalf("get error: %s", err)
	}
	if !reflect.DeepEqual(actual, entry) {
		t.Fatalf("wrong value. Expected: %#v\nGot: %#v", entry, actual)
	}

	keys, err = s.List("")
	if err != nil {
		t.Fatalf("list error: %s", err)
	}
	if !reflect.DeepEqual(keys, []string{"foo"}) {
		t.Fatalf("bad keys: %#v", keys)
	}

	if err := s.Delete("foo"); err != nil {
		t.Fatalf("put error: %s", err)
	}

	keys, err = s.List("")
	if err != nil {
		t.Fatalf("list error: %s", err)
	}
	if len(keys) > 0 {
		t.Fatalf("should have no keys to start: %#v", keys)
	}
}

func TestSystemView() *StaticSystemView {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 2
	return &StaticSystemView{
		DefaultLeaseTTLVal: defaultLeaseTTLVal,
		MaxLeaseTTLVal:     maxLeaseTTLVal,
	}
}

func TestBackendConfig() *BackendConfig {
	bc := &BackendConfig{
		Logger: logformat.NewVaultLogger(log.LevelTrace),
		System: TestSystemView(),
	}
	bc.Logger.SetLevel(log.LevelTrace)

	return bc
}
