package jwt

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

type jwtMethod struct {
	logger                hclog.Logger
	path                  string
	mountPath             string
	role                  string
	removeJWTAfterReading bool
	credsFound            chan struct{}
	watchCh               chan string
	stopCh                chan struct{}
	doneCh                chan struct{}
	credSuccessGate       chan struct{}
	ticker                *time.Ticker
	once                  *sync.Once
	latestToken           *atomic.Value
}

// NewJWTAuthMethod returns an implementation of Agent's auth.AuthMethod
// interface for JWT auth.
func NewJWTAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	j := &jwtMethod{
		logger:                conf.Logger,
		mountPath:             conf.MountPath,
		removeJWTAfterReading: true,
		credsFound:            make(chan struct{}),
		watchCh:               make(chan string),
		stopCh:                make(chan struct{}),
		doneCh:                make(chan struct{}),
		credSuccessGate:       make(chan struct{}),
		once:                  new(sync.Once),
		latestToken:           new(atomic.Value),
	}
	j.latestToken.Store("")

	pathRaw, ok := conf.Config["path"]
	if !ok {
		return nil, errors.New("missing 'path' value")
	}
	j.path, ok = pathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'path' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	j.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	if removeJWTAfterReadingRaw, ok := conf.Config["remove_jwt_after_reading"]; ok {
		removeJWTAfterReading, err := parseutil.ParseBool(removeJWTAfterReadingRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'remove_jwt_after_reading' value: %w", err)
		}
		j.removeJWTAfterReading = removeJWTAfterReading
	}

	switch {
	case j.path == "":
		return nil, errors.New("'path' value is empty")
	case j.role == "":
		return nil, errors.New("'role' value is empty")
	}

	// If we don't delete the JWT after reading, use a slower reload period,
	// otherwise we would re-read the whole file every 500ms, instead of just
	// doing a stat on the file every 500ms.
	readPeriod := 1 * time.Minute
	if j.removeJWTAfterReading {
		readPeriod = 500 * time.Millisecond
	}

	// This is test-only config, so that tests don't need to wait for minutes
	// in the case that removeJWTAfterReading != true
	if testOverrideReadPeriodRaw, ok := conf.Config["test_override_read_period"]; ok {
		testOverrideReadPeriod, err := parseutil.ParseBool(testOverrideReadPeriodRaw)
		if err == nil {
			if testOverrideReadPeriod {
				readPeriod = 500 * time.Millisecond
			}
		}
		// We do nothing if err is non-nil, as this is meant to be a test-only config.
	}

	j.ticker = time.NewTicker(readPeriod)

	go j.runWatcher()

	j.logger.Info("jwt auth method created", "path", j.path)

	return j, nil
}

func (j *jwtMethod) Authenticate(_ context.Context, _ *api.Client) (string, http.Header, map[string]interface{}, error) {
	j.logger.Trace("beginning authentication")

	j.ingressToken()

	latestToken := j.latestToken.Load().(string)
	if latestToken == "" {
		return "", nil, nil, errors.New("latest known jwt is empty, cannot authenticate")
	}

	return fmt.Sprintf("%s/login", j.mountPath), nil, map[string]interface{}{
		"role": j.role,
		"jwt":  latestToken,
	}, nil
}

func (j *jwtMethod) NewCreds() chan struct{} {
	return j.credsFound
}

func (j *jwtMethod) CredSuccess() {
	j.once.Do(func() {
		close(j.credSuccessGate)
	})
}

func (j *jwtMethod) Shutdown() {
	j.ticker.Stop()
	close(j.stopCh)
	<-j.doneCh
}

func (j *jwtMethod) runWatcher() {
	defer close(j.doneCh)

	select {
	case <-j.stopCh:
		return

	case <-j.credSuccessGate:
		// We only start the next loop once we're initially successful,
		// since at startup Authenticate will be called, and we don't want
		// to end up immediately re-authenticating by having found a new
		// value
	}

	for {
		select {
		case <-j.stopCh:
			return

		case <-j.ticker.C:
			latestToken := j.latestToken.Load().(string)
			j.ingressToken()
			newToken := j.latestToken.Load().(string)
			if newToken != latestToken {
				j.logger.Debug("new jwt file found")
				j.credsFound <- struct{}{}
			}
		}
	}
}

func (j *jwtMethod) ingressToken() {
	fi, err := os.Lstat(j.path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		j.logger.Error("error encountered stat'ing jwt file", "error", err)
		return
	}

	// Check that the path refers to a file.
	// If it's a symlink, it could still be a symlink to a directory,
	// but os.ReadFile below will return a descriptive error.
	var symlink bool
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		// regular file
	case mode&fs.ModeSymlink != 0:
		symlink = true
	default:
		j.logger.Error("jwt file is not a regular file or symlink")
		return
	}

	evalSymlinkPath := j.path
	// If our file path is a symlink, we should also return early (like above) without error
	// if the file that is linked to is not present, otherwise we will error when trying
	// to read that file by following the link in the os.ReadFile call.
	if symlink {
		evalSymlinkPath, err = filepath.EvalSymlinks(j.path)
		if err != nil {
			j.logger.Error("error encountered evaluating symlinks", "error", err)
			return
		}
		_, err := os.Lstat(evalSymlinkPath)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			j.logger.Error("error encountered stat'ing jwt file after evaluating symlinks", "error", err)
			return
		}
	}

	token, err := os.ReadFile(j.path)
	if err != nil {
		j.logger.Error("failed to read jwt file", "error", err)
		return
	}

	switch len(token) {
	case 0:
		j.logger.Warn("empty jwt file read")

	default:
		j.latestToken.Store(string(token))
	}

	if j.removeJWTAfterReading {
		// We remove the evalSymlink'd path here so that we actually delete the jwt
		// not just the symlink that links to the jwt
		if err := os.Remove(evalSymlinkPath); err != nil {
			j.logger.Error("error removing jwt file", "error", err)
		}
	}
}
