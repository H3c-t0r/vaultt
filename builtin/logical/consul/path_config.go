package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"address": {
				Type:        framework.TypeString,
				Description: "Consul server address",
			},

			"scheme": {
				Type:        framework.TypeString,
				Description: "URI scheme for the Consul address",

				// https would be a better default but Consul on its own
				// defaults to HTTP access, and when HTTPS is enabled it
				// disables HTTP, so there isn't really any harm done here.
				Default: "http",
			},

			"token": {
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},

			"ca_cert": {
				Type: framework.TypeString,
				Description: `CA certificate to use when verifying Consul server certificate,
must be x509 PEM encoded.`,
			},

			"client_cert": {
				Type: framework.TypeString,
				Description: `Client certificate used for Consul's TLS communication,
must be x509 PEM encoded and if this is set you need to also set client_key.`,
			},

			"client_key": {
				Type: framework.TypeString,
				Description: `Client key used for Consul's TLS communication,
must be x509 PEM encoded and if this is set you need to also set client_cert.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigAccessRead,
			logical.UpdateOperation: b.pathConfigAccessWrite,
		},
	}
}

func (b *backend) readConfigAccess(ctx context.Context, storage logical.Storage) (*accessConfig, error, error) {
	entry, err := storage.Get(ctx, "config/access")
	if err != nil {
		return nil, nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("access credentials for the backend itself haven't been configured; please configure them at the '/config/access' endpoint"), nil
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, nil, errwrap.Wrapf("error reading consul access configuration: {{err}}", err)
	}

	return conf, nil, nil
}

func (b *backend) pathConfigAccessRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, userErr, intErr := b.readConfigAccess(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}
	if conf == nil {
		return nil, fmt.Errorf("no user error reported but consul access configuration not found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address": conf.Address,
			"scheme":  conf.Scheme,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/access", accessConfig{
		Address:    data.Get("address").(string),
		Scheme:     data.Get("scheme").(string),
		Token:      data.Get("token").(string),
		CACert:     data.Get("ca_cert").(string),
		ClientCert: data.Get("client_cert").(string),
		ClientKey:  data.Get("client_key").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type accessConfig struct {
	Address    string `json:"address"`
	Scheme     string `json:"scheme"`
	Token      string `json:"token"`
	CACert     string `json:"ca_cert"`
	ClientCert string `json:"client_cert"`
	ClientKey  string `json:"client_key"`
}
