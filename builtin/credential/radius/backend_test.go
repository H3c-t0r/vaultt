package radius

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

const (
	testSysTTL    = time.Hour * 10
	testSysMaxTTL = time.Hour * 20
)

func TestBackend_Config(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	config_data_basic := map[string]interface{}{
		"host":   "test.radius.hostname.com",
		"secret": "test-secret",
	}

	config_data_missingrequired := map[string]interface{}{
		"host": "test.radius.hostname.com",
	}

	config_data_invalidport := map[string]interface{}{
		"host":   "test.radius.hostname.com",
		"port":   "notnumeric",
		"secret": "test-secret",
	}

	config_data_invalidbool := map[string]interface{}{
		"host":                  "test.radius.hostname.com",
		"secret":                "test-secret",
		"enable_default_policy": "test",
	}

	config_data_emptyport := map[string]interface{}{
		"host":   "test.radius.hostname.com",
		"port":   "",
		"secret": "test-secret",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		// PreCheck:       func() { testAccPreCheck(t) },
		Backend: b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, config_data_basic, false),
			testConfigWrite(t, config_data_emptyport, true),
			testConfigWrite(t, config_data_invalidport, true),
			testConfigWrite(t, config_data_invalidbool, true),
			testConfigWrite(t, config_data_missingrequired, true),
		},
	})
}

func TestBackend_users(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testStepUpdateUser(t, "web", "foo"),
			testStepUpdateUser(t, "web2", "foo"),
			testStepUpdateUser(t, "web3", "foo"),
			testStepUserList(t, []string{"web", "web2", "web3"}),
		},
	})
}

func TestBackend_acceptance(t *testing.T) {
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: testSysTTL,
			MaxLeaseTTLVal:     testSysMaxTTL,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	config_data_acceptance_defpol := map[string]interface{}{
		"host":                  os.Getenv("RADIUS_HOST"),
		"port":                  os.Getenv("RADIUS_PORT"),
		"secret":                os.Getenv("RADIUS_SECRET"),
		"enable_default_policy": "true",
	}

	if config_data_acceptance_defpol["port"] == "" {
		config_data_acceptance_defpol["port"] = "1812"
	}

	config_data_acceptance_nodefpol := map[string]interface{}{
		"host":                  os.Getenv("RADIUS_HOST"),
		"port":                  os.Getenv("RADIUS_PORT"),
		"secret":                os.Getenv("RADIUS_SECRET"),
		"enable_default_policy": "false",
	}

	data_realpassword := map[string]interface{}{
		"password": os.Getenv("RADIUS_USERPASS"),
	}

	data_wrongpassword := map[string]interface{}{
		"password": "wrongpassword",
	}

	username := os.Getenv("RADIUS_USERNAME")

	if config_data_acceptance_nodefpol["port"] == "" {
		config_data_acceptance_nodefpol["port"] = "1812"
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend:        b,
		PreCheck:       func() { testAccPreCheck(t) },
		AcceptanceTest: true,
		Steps: []logicaltest.TestStep{
			// Login with valid but unknown user will fail since enable_default_policy is false
			testConfigWrite(t, config_data_acceptance_nodefpol, false),
			testAccUserLogin(t, username, data_realpassword, true),
			// Once the user is registered auth will succeed
			testStepUpdateUser(t, username, ""),
			testAccUserLoginPolicy(t, username, data_realpassword, []string{"default"}, false),

			testStepUpdateUser(t, username, "foopolicy"),
			testAccUserLoginPolicy(t, username, data_realpassword, []string{"default", "foopolicy"}, false),
			testAccStepDeleteUser(t, username),

			// When using enable_default_policy, an unknown user will be allowed to authenticate and given the default policy
			testConfigWrite(t, config_data_acceptance_defpol, false),
			testAccUserLoginPolicy(t, username, data_realpassword, []string{"default"}, false),

			// More tests
			testAccUserLogin(t, "nonexistinguser", data_realpassword, true),
			testAccUserLogin(t, username, data_wrongpassword, true),
			testStepUpdateUser(t, username, "foopolicy"),
			testAccUserLoginPolicy(t, username, data_realpassword, []string{"default", "foopolicy"}, false),
			testStepUpdateUser(t, username, "foopolicy, secondpolicy"),
			testAccUserLoginPolicy(t, username, data_realpassword, []string{"default", "foopolicy", "secondpolicy"}, false),
			testAccUserLoginPolicy(t, username, data_realpassword, []string{"default", "foopolicy", "secondpolicy", "thirdpolicy"}, true),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("RADIUS_HOST"); v == "" {
		t.Fatal("RADIUS_HOST must be set for acceptance tests")
	}

	if v := os.Getenv("RADIUS_USERNAME"); v == "" {
		t.Fatal("RADIUS_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("RADIUS_USERPASS"); v == "" {
		t.Fatal("RADIUS_USERPASS must be set for acceptance tests")
	}

	if v := os.Getenv("RADIUS_SECRET"); v == "" {
		t.Fatal("RADIUS_SECRET must be set for acceptance tests")
	}
}

func testConfigWrite(t *testing.T, d map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
		ErrorOk:   expectError,
	}
}

func testAccStepDeleteUser(t *testing.T, n string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "users/" + n,
	}
}

func testStepUserList(t *testing.T, users []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "users",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return fmt.Errorf("Got error response: %#v", *resp)
			}

			if !reflect.DeepEqual(users, resp.Data["keys"].([]string)) {
				return fmt.Errorf("expected:\n%#v\ngot:\n%#v\n", users, resp.Data["keys"])
			}
			return nil
		},
	}
}

func testStepUpdateUser(
	t *testing.T, name string, policies string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + name,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccUserLogin(t *testing.T, user string, data map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login/" + user,
		Data:            data,
		ErrorOk:         expectError,
		Unauthenticated: true,
	}
}

func testAccUserLoginPolicy(t *testing.T, user string, data map[string]interface{}, policies []string, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.UpdateOperation,
		Path:            "login/" + user,
		Data:            data,
		ErrorOk:         false,
		Unauthenticated: true,
		//Check:           logicaltest.TestCheckAuth(policies),
		Check: func(resp *logical.Response) error {
			res := logicaltest.TestCheckAuth(policies)(resp)
			if res != nil && expectError {
				return nil
			}
			return res
		},
	}
}
