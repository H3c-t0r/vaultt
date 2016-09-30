package awsiam

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathConfigClient(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b := Backend()
	_, err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// make sure we start with empty roles, which gives us confidence that the read later
	// actually is the two roles we created
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	// at this point, resp == nil is valid as no client config exists
	// if resp != nil, then resp.Data must have EndPoint and HeaderValue as nil
	if resp != nil {
		if resp.IsError() {
			t.Fatalf("failed to read client config entry")
		} else if resp.Data["endpoint"] != nil || resp.Data["vault_header_value"] != nil {
			t.Fatalf("Returned endpoint or vault_header_value non-nil")
		}
	}

	data := map[string]interface{}{
		"endpoint":           "https://my-custom-sts-endpoint.example.com",
		"vault_header_value": "vault_server_identification_314159",
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Data:      data,
		Storage:   storage,
	})

	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("failed to create the client config entry")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("failed to read the client config entry")
	}
	if resp.Data["vault_header_value"] != data["vault_header_value"] {
		t.Fatalf("Expected vault_header_value: '%#v'; returned vault_header_value: '%#v'",
			data["vault_header_value"], resp.Data["vault_header_value"])
	}
}
