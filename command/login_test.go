package command

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/cli"

	"github.com/hashicorp/vault/api"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/vault"
)

// minTokenLengthExternal is the minimum size of SSC
// tokens we are currently handing out to end users, without any
// namespace information
const minTokenLengthExternal = 91

func testLoginCommand(tb testing.TB) (*cli.MockUi, *LoginCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &LoginCommand{
		BaseCommand: &BaseCommand{
			UI: ui,

			// Override to our own token helper
			tokenHelper: token.NewTestingTokenHelper(),
		},
		Handlers: map[string]LoginHandler{
			"token":    &credToken.CLIHandler{},
			"userpass": &credUserpass.CLIHandler{},
		},
	}
}

func TestLoginCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("custom_path", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("my-auth", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/my-auth/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		// Emulate an unknown token format present in ~/.vault-token, for example
		client.SetToken("a.a")

		code := cmd.Run([]string{
			"-method", "userpass",
			"-path", "my-auth",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! You are now authenticated."
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to be %q", combined, expected)
		}

		storedToken, err := tokenHelper.Get()
		if err != nil {
			t.Fatal(err)
		}

		if l, exp := len(storedToken), minTokenLengthExternal+vault.TokenPrefixLength; l < exp {
			t.Errorf("expected token to be %d characters, was %d: %q", exp, l, storedToken)
		}
	})

	t.Run("no_store", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			TTL:      "30m",
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		_, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		// Ensure we have no token to start
		if storedToken, err := tokenHelper.Get(); err != nil || storedToken != "" {
			t.Errorf("expected token helper to be empty: %s: %q", err, storedToken)
		}

		code := cmd.Run([]string{
			"-no-store",
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		storedToken, err := tokenHelper.Get()
		if err != nil {
			t.Fatal(err)
		}

		if exp := ""; storedToken != exp {
			t.Errorf("expected %q to be %q", storedToken, exp)
		}
	})

	t.Run("stores", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
			TTL:      "30m",
		})
		if err != nil {
			t.Fatal(err)
		}
		token := secret.Auth.ClientToken

		_, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			token,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		storedToken, err := tokenHelper.Get()
		if err != nil {
			t.Fatal(err)
		}

		if storedToken != token {
			t.Errorf("expected %q to be %q", storedToken, token)
		}
	})

	t.Run("token_only", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"-token-only",
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Verify only the token was printed
		token := ui.OutputWriter.String()
		if l, exp := len(token), minTokenLengthExternal+vault.TokenPrefixLength; l != exp {
			t.Errorf("expected token to be %d characters, was %d: %q", exp, l, token)
		}

		// Verify the token was not stored
		if storedToken, err := tokenHelper.Get(); err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}
	})

	t.Run("failure_no_store", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}

		code := cmd.Run([]string{
			"not-a-real-token",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error authenticating: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		if storedToken, err := tokenHelper.Get(); err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}
	})

	t.Run("wrap_auto_unwrap", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		_, cmd := testLoginCommand(t)
		cmd.client = client

		// Set the wrapping ttl to 5s. We can't set this via the flag because we
		// override the client object before that particular flag is parsed.
		client.SetWrappingLookupFunc(func(string, string) string { return "5m" })

		code := cmd.Run([]string{
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Unset the wrapping
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		token, err := tokenHelper.Get()
		if err != nil || token == "" {
			t.Fatalf("expected token from helper: %s: %q", err, token)
		}
		client.SetToken(token)

		// Ensure the resulting token is unwrapped
		secret, err := client.Auth().Token().LookupSelf()
		if err != nil {
			t.Error(err)
		}
		if secret == nil {
			t.Fatal("secret was nil")
		}

		if secret.WrapInfo != nil {
			t.Errorf("expected to be unwrapped: %#v", secret)
		}
	})

	t.Run("wrap_token_only", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		// Set the wrapping ttl to 5s. We can't set this via the flag because we
		// override the client object before that particular flag is parsed.
		client.SetWrappingLookupFunc(func(string, string) string { return "5m" })

		code := cmd.Run([]string{
			"-token-only",
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Unset the wrapping
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		storedToken, err := tokenHelper.Get()
		if err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}

		token := strings.TrimSpace(ui.OutputWriter.String())
		if token == "" {
			t.Errorf("expected %q to not be %q", token, "")
		}

		// Ensure the resulting token is, in fact, still wrapped.
		client.SetToken(token)
		secret, err := client.Logical().Unwrap("")
		if err != nil {
			t.Error(err)
		}
		if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
			t.Fatalf("expected secret to have auth: %#v", secret)
		}
	})

	t.Run("wrap_no_store", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
			"password": "test",
			"policies": "default",
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		// Set the wrapping ttl to 5s. We can't set this via the flag because we
		// override the client object before that particular flag is parsed.
		client.SetWrappingLookupFunc(func(string, string) string { return "5m" })

		code := cmd.Run([]string{
			"-no-store",
			"-method", "userpass",
			"username=test",
			"password=test",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		// Unset the wrapping
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		storedToken, err := tokenHelper.Get()
		if err != nil || storedToken != "" {
			t.Fatalf("expected token to not be stored: %s: %q", err, storedToken)
		}

		expected := "wrapping_token"
		output := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(output, expected) {
			t.Errorf("expected %q to contain %q", output, expected)
		}
	})

	t.Run("login_mfa_single_phase", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		userclient, entityID, methodID, methodName := testhelpers.SetupLoginMFATOTP(t, client)

		enginePath := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

		runCommand := func(methodIdentifier string) {
			time.Sleep(21 * time.Second)
			totpCode := testhelpers.GetTOTPCodeFromEngine(t, client, enginePath)
			ui, cmd := testLoginCommand(t)
			cmd.client = userclient

			// login command bails early for test clients, so we have to explicitly set this
			cmd.client.SetMFACreds([]string{methodIdentifier + ":" + totpCode})
			code := cmd.Run([]string{
				"-method", "userpass",
				"username=testuser1",
				"password=testpassword",
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d", code, exp)
			}

			tokenHelper, err := cmd.TokenHelper()
			if err != nil {
				t.Fatal(err)
			}
			storedToken, err := tokenHelper.Get()
			if err != nil {
				t.Fatal(err)
			}
			output = ui.OutputWriter.String() + ui.ErrorWriter.String()
			t.Logf("\n%+v", output)
			if !strings.Contains(output, storedToken) {
				t.Fatalf("expected stored token: %q, got: %q", storedToken, output)
			}
		}
		runCommand(methodID)
		runCommand(methodName)
	})

	t.Run("login_mfa_two_phase", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testLoginCommand(t)

		userclient, entityID, methodID, _ := testhelpers.SetupLoginMFATOTP(t, client)
		cmd.client = userclient

		_ = testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

		// clear the MFA creds just to be sure
		cmd.client.SetMFACreds([]string{})

		code := cmd.Run([]string{
			"-method", "userpass",
			"username=testuser1",
			"password=testpassword",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := methodID
		output = ui.OutputWriter.String() + ui.ErrorWriter.String()
		t.Logf("\n%+v", output)
		if !strings.Contains(output, expected) {
			t.Fatalf("expected stored token: %q, got: %q", expected, output)
		}

		tokenHelper, err := cmd.TokenHelper()
		if err != nil {
			t.Fatal(err)
		}
		storedToken, err := tokenHelper.Get()
		if storedToken != "" {
			t.Fatal("expected empty stored token")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("login_mfa_two_phase_non_interactive_method_name", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testLoginCommand(t)

		userclient, entityID, methodID, methodName := testhelpers.SetupLoginMFATOTP(t, client)
		cmd.client = userclient

		engineName := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

		// clear the MFA creds just to be sure
		cmd.client.SetMFACreds([]string{})

		code := cmd.Run([]string{
			"-method", "userpass",
			"-non-interactive",
			"username=testuser1",
			"password=testpassword",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		output = ui.OutputWriter.String() + ui.ErrorWriter.String()

		reqIdReg, err := regexp.Compile(`mfa_request_id\s+(?P<name>[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})?\s+mfa_constraint.*`)
		if err != nil {
			t.Fatalf("failed to compile regex")
		}

		reqIDRaw := reqIdReg.FindAllStringSubmatch(output, -1)
		mfaReqID := reqIDRaw[0][1]

		validateFunc := func(methodIdentifier string) {
			// the time required for the totp engine to generate a new code
			time.Sleep(22 * time.Second)
			totpPasscode1 := "passcode=" + testhelpers.GetTOTPCodeFromEngine(t, client, engineName)

			secret, err := cmd.client.Logical().WriteWithContext(context.Background(), "sys/mfa/validate", map[string]interface{}{
				"mfa_request_id": mfaReqID,
				"mfa_payload": map[string][]string{
					methodIdentifier: {totpPasscode1},
				},
			})
			if err != nil {
				t.Fatalf("mfa validation failed: %v", err)
			}

			if secret.Auth == nil || secret.Auth.ClientToken == "" {
				t.Fatalf("mfa validation did not return a client token")
			}
		}

		validateFunc(methodName)
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testLoginCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"token",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error authenticating: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testLoginCommand(t)
		assertNoTabs(t, cmd)
	})
}
