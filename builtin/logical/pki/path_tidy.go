package pki

import (
	"context"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy",
		Fields: map[string]*framework.FieldSchema{
			"tidy_cert_store": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Set to true to enable tidying up
the certificate store`,
				Default: false,
			},

			"tidy_revocation_list": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Set to true to enable tidying up
the revocation list`,
				Default: false,
			},

			"safety_buffer": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `The amount of extra time that must have passed
beyond certificate expiration before it is removed
from the backend storage and/or revocation list.
Defaults to 72 hours.`,
				Default: 259200, //72h, but TypeDurationSecond currently requires defaults to be int
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTidyWrite,
		},

		HelpSynopsis:    pathTidyHelpSyn,
		HelpDescription: pathTidyHelpDesc,
	}
}

func (b *backend) pathTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	safetyBuffer := d.Get("safety_buffer").(int)
	tidyCertStore := d.Get("tidy_cert_store").(bool)
	tidyRevocationList := d.Get("tidy_revocation_list").(bool)

	if safetyBuffer < 1 {
		return logical.ErrorResponse("safety_buffer must be greater than zero"), nil
	}

	bufferDuration := time.Duration(safetyBuffer) * time.Second

	var resp *logical.Response

	if tidyCertStore {
		serials, err := req.Storage.List(ctx, "certs/")
		if err != nil {
			return nil, errwrap.Wrapf("error fetching list of certs: {{err}}", err)
		}

		for _, serial := range serials {
			certEntry, err := req.Storage.Get(ctx, "certs/"+serial)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("error fetching certificate %q: {{err}}", serial), err)
			}

			if certEntry == nil {
				if resp == nil {
					resp = &logical.Response{}
				}
				resp.AddWarning(fmt.Sprintf("Certificate entry for serial %s is nil; tidying up since it is no longer useful for any server operations", serial))
				if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error deleting nil entry with serial %s: {{err}}", serial), err)
				}
			}

			if certEntry.Value == nil || len(certEntry.Value) == 0 {
				if resp == nil {
					resp = &logical.Response{}
				}
				resp.AddWarning(fmt.Sprintf("Certificate entry for serial %s is nil; tidying up since it is no longer useful for any server operations", serial))
				if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error deleting entry with nil value with serial %s: {{err}}", serial), err)
				}
			}

			cert, err := x509.ParseCertificate(certEntry.Value)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("unable to parse stored certificate with serial %q: {{err}}", serial), err)
			}

			if time.Now().After(cert.NotAfter.Add(bufferDuration)) {
				if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error deleting serial %q from storage: {{err}}", serial), err)
				}
			}
		}
	}

	if tidyRevocationList {
		b.revokeStorageLock.Lock()
		defer b.revokeStorageLock.Unlock()

		tidiedRevoked := false

		revokedSerials, err := req.Storage.List(ctx, "revoked/")
		if err != nil {
			return nil, errwrap.Wrapf("error fetching list of revoked certs: {{err}}", err)
		}

		var revInfo revocationInfo
		for _, serial := range revokedSerials {
			revokedEntry, err := req.Storage.Get(ctx, "revoked/"+serial)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("unable to fetch revoked cert with serial %q: {{err}}", serial), err)
			}

			if revokedEntry == nil {
				if resp == nil {
					resp = &logical.Response{}
				}
				resp.AddWarning(fmt.Sprintf("Revoked entry for serial %s is nil; tidying up since it is no longer useful for any server operations", serial))
				if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error deleting nil revoked entry with serial %s: {{err}}", serial), err)
				}
			}

			if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
				if resp == nil {
					resp = &logical.Response{}
				}
				resp.AddWarning(fmt.Sprintf("Revoked entry for serial %s has nil value; tidying up since it is no longer useful for any server operations", serial))
				if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error deleting revoked entry with nil value with serial %s: {{err}}", serial), err)
				}
			}

			err = revokedEntry.DecodeJSON(&revInfo)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("error decoding revocation entry for serial %q: {{err}}", serial), err)
			}

			revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("unable to parse stored revoked certificate with serial %q: {{err}}", serial), err)
			}

			if time.Now().After(revokedCert.NotAfter.Add(bufferDuration)) {
				if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error deleting serial %q from revoked list: {{err}}", serial), err)
				}
				tidiedRevoked = true
			}
		}

		if tidiedRevoked {
			if err := buildCRL(ctx, b, req); err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}

const pathTidyHelpSyn = `
Tidy up the backend by removing expired certificates, revocation information,
or both.
`

const pathTidyHelpDesc = `
This endpoint allows expired certificates and/or revocation information to be
removed from the backend, freeing up storage and shortening CRLs.

For safety, this function is a noop if called without parameters; cleanup from
normal certificate storage must be enabled with 'tidy_cert_store' and cleanup
from revocation information must be enabled with 'tidy_revocation_list'.

The 'safety_buffer' parameter is useful to ensure that clock skew amongst your
hosts cannot lead to a certificate being removed from the CRL while it is still
considered valid by other hosts (for instance, if their clocks are a few
minutes behind). The 'safety_buffer' parameter can be an integer number of
seconds or a string duration like "72h".

All certificates and/or revocation information currently stored in the backend
will be checked when this endpoint is hit. The expiration of the
certificate/revocation information of each certificate being held in
certificate storage or in revocation information will then be checked. If the
current time, minus the value of 'safety_buffer', is greater than the
expiration, it will be removed.
`
