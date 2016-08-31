package physical

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/helper/jsonutil"
)

// FileBackend is a physical backend that stores data on disk
// at a given file path. It can be used for durable single server
// situations, or to develop locally where durability is not critical.
//
// WARNING: the file backend implementation is currently extremely unsafe
// and non-performant. It is meant mostly for local testing and development.
// It can be improved in the future.
type FileBackend struct {
	Path   string
	l      sync.Mutex
	logger log.Logger
}

// newFileBackend constructs a Filebackend using the given directory
func newFileBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	return &FileBackend{
		Path:   path,
		logger: logger,
	}, nil
}

func (b *FileBackend) Delete(k string) error {
	b.l.Lock()
	defer b.l.Unlock()

	path, key := b.path(k)
	fullPath := filepath.Join(path, key)

	// If the path doesn't exist return success; this is in line with Vault's
	// expected behavior and we don't want to check for an empty directory if
	// we couldn't even find the path in the first place.
	err := os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	// Check for the directory being empty and remove if so, with another
	// additional guard for the path not existing
	dir, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	list, err := dir.Readdir(1)
	dir.Close()
	if err != nil && err != io.EOF {
		return err
	}

	// If we have no entries, it's an empty directory; remove it
	if err == io.EOF || list == nil || len(list) == 0 {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *FileBackend) Get(k string) (*Entry, error) {
	b.l.Lock()
	defer b.l.Unlock()

	path, key := b.path(k)
	path = filepath.Join(path, key)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	var entry Entry
	if err := jsonutil.DecodeJSONFromReader(f, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func (b *FileBackend) Put(entry *Entry) error {
	path, key := b.path(entry.Key)

	b.l.Lock()
	defer b.l.Unlock()

	// Make the parent tree
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// JSON encode the entry and write it
	f, err := os.OpenFile(
		filepath.Join(path, key),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(entry)
}

func (b *FileBackend) List(prefix string) ([]string, error) {
	b.l.Lock()
	defer b.l.Unlock()

	path := b.Path
	if prefix != "" {
		path = filepath.Join(path, prefix)
	}

	// Read the directory contents
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for i, name := range names {
		if name[0] == '_' {
			names[i] = name[1:]
		} else {
			names[i] = name + "/"
		}
	}

	return names, nil
}

func (b *FileBackend) path(k string) (string, string) {
	path := filepath.Join(b.Path, k)
	key := filepath.Base(path)
	path = filepath.Dir(path)
	return path, "_" + key
}
