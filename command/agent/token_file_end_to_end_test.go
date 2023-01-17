package agent

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/auth"
	token_file "github.com/hashicorp/vault/command/agent/auth/token-file"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

func TestTokenFileEndToEnd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		removeTokenFile bool
		expectToken     bool
	}{
		// default behaviour => token expected
		{false, true},
		{true, true},
	}

	for _, tc := range testCases {
		secretFileAction := "preserve"
		if tc.removeTokenFile {
			secretFileAction = "remove"
		}
		tc := tc // capture range variable
		t.Run(fmt.Sprintf("%s_removeTokenFile, expectToken=%v", secretFileAction, tc.expectToken), func(t *testing.T) {
			t.Parallel()
			testTokenFileEndToEnd(t, tc.removeTokenFile, tc.expectToken)
		})
	}
}

func testTokenFileEndToEnd(t *testing.T, removeTokenFile bool, expectToken bool) {
	var err error
	logger := logging.NewVaultLogger(log.Trace)
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	secret, err := client.Auth().Token().Create(nil)
	if err != nil || secret == nil {
		t.Fatal(err)
	}

	tokenFile, err := os.CreateTemp("", "token_file")
	if err != nil {
		t.Fatal(err)
	}
	tokenFileName := tokenFile.Name()
	tokenFile.Close() // WriteFile doesn't need it open
	os.WriteFile(tokenFileName, []byte(secret.Auth.ClientToken), 0o666)
	defer os.Remove(tokenFileName)

	t.Logf("input token_file_path: %s", tokenFileName)

	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	am, err := token_file.NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path":                 tokenFileName,
			"remove_token_file_after_reading": removeTokenFile,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ah := auth.NewAuthHandler(ahConfig)
	errCh := make(chan error)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// We close these right away because we're just basically testing
	// permissions and finding a usable file name
	sinkFile, err := os.CreateTemp("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	tokenSinkFileName := sinkFile.Name()
	sinkFile.Close()
	os.Remove(tokenSinkFileName)
	t.Logf("output: %s", tokenSinkFileName)

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": tokenSinkFileName,
		},
		WrapTTL: 10 * time.Second,
	}

	fs, err := file.NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})
	go func() {
		errCh <- ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config})
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// This has to be after the other defers, so it happens first. It allows
	// successful test runs to immediately cancel all of the runner goroutines
	// and unblock any of the blocking defer calls by the runner's DoneCh that
	// comes before this and avoid successful tests from taking the entire
	// timeout duration.
	defer cancel()

	if stat, err := os.Lstat(tokenSinkFileName); err == nil {
		t.Fatalf("expected err but got %s", stat)
	} else if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	if expectToken {
		// Wait 2 seconds for the env variables to be detected and an auth to be generated.
		time.Sleep(time.Second * 2)

		token, err := readToken(tokenSinkFileName)
		if err != nil {
			t.Fatal(err)
		}

		if token.Token == "" {
			t.Fatal("expected token but didn't receive it")
		}
	}

	if !removeTokenFile {
		_, err := os.Stat(tokenFileName)
		if err != nil {
			t.Fatal("Token file removed despite remove token file being set to false")
		}
	} else {
		_, err := os.Stat(tokenFileName)
		if err == nil {
			t.Fatal("no error returned from stat, indicating the file is still present")
		}
		if !os.IsNotExist(err) {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
