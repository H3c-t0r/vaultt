package azureservicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigLease(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/lease",
		Fields: map[string]*framework.FieldSchema{
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Default lease time for the SAS token",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigLeaseRead,
			logical.UpdateOperation: b.pathConfigLeaseWrite,
		},

		HelpSynopsis:    pathConfigLeaseHelpSyn,
		HelpDescription: pathConfigLeaseHelpDesc,
	}
}

func (b *backend) pathConfigLeaseWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ttlRaw := d.Get("ttl").(string)

	ttl, err := time.ParseDuration(ttlRaw)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Invalid lease time: %s", err)), nil
	}

	// Store it
	entry, err := logical.StorageEntryJSON("config/lease", &configLease{
		TTL: ttl,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigLeaseRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	leaseConfig, err := b.LeaseConfig(req.Storage)

	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"ttl": leaseConfig.TTL.String(),
		},
	}, nil
}

type configLease struct {
	TTL time.Duration
}

const pathConfigLeaseHelpSyn = `
Configure the default lease time for generated SAS token.
`

const pathConfigLeaseHelpDesc = `
This configures the default lease time used for SAS tokens
generated by this backend. The ttl specifies the duration that a
token will be valid for. SAS tokens will not renew after lease expires.

The format is "1h" or integer and then unit. The longest unit is hour.
`
