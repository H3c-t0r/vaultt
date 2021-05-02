package diagnose

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
)

const (
	success   string = "success"
	secretKey string = "diagnose"
	secretVal string = "diagnoseSecret"

	timeOutErr        string = "storage call timed out after 20 seconds: "
	wrongRWValsPrefix string = "Storage get and put gave wrong values: "
)

// StorageEndToEndLatencyCheck calls Write, Read, and Delete on a secret in the root
// directory of the backend.
// Note: Just checking read, write, and delete for root. It's a very basic check,
// but I don't think we can necessarily do any more than that. We could check list,
// but I don't think List is ever going to break in isolation.
func StorageEndToEndLatencyCheck(ctx context.Context, b physical.Backend) error {

	c2 := make(chan error)
	go func() {
		err := b.Put(context.Background(), &physical.Entry{Key: secretKey, Value: []byte(secretVal)})
		c2 <- err
	}()
	select {
	case errOut := <-c2:
		if errOut != nil {
			return errOut
		}
	case <-time.After(20 * time.Second):
		return fmt.Errorf(timeOutErr + "operation: Put")
	}

	c3 := make(chan *physical.Entry)
	c4 := make(chan error)
	go func() {
		val, err := b.Get(context.Background(), "diagnose")
		if err != nil {
			c4 <- err
		} else {
			c3 <- val
		}
	}()
	select {
	case err := <-c4:
		return err
	case val := <-c3:
		if val.Key != "diagnose" && string(val.Value) != "diagnose" {
			return fmt.Errorf(wrongRWValsPrefix+"expecting diagnose, but got %s, %s", val.Key, val.Value)
		}
	case <-time.After(20 * time.Second):
		return fmt.Errorf(timeOutErr + "operation: Get")
	}

	c5 := make(chan error)
	go func() {
		err := b.Delete(context.Background(), "diagnose")
		c5 <- err
	}()
	select {
	case errOut := <-c5:
		if errOut != nil {
			return errOut
		}
	case <-time.After(20 * time.Second):
		return fmt.Errorf(timeOutErr + "operation: Delete")
	}
	return nil
}
