package pki

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy$",
		Fields: map[string]*framework.FieldSchema{
			"tidy_cert_store": {
				Type: framework.TypeBool,
				Description: `Set to true to enable tidying up
the certificate store`,
			},

			"tidy_revocation_list": {
				Type:        framework.TypeBool,
				Description: `Deprecated; synonym for 'tidy_revoked_certs`,
			},

			"tidy_revoked_certs": {
				Type: framework.TypeBool,
				Description: `Set to true to expire all revoked
and expired certificates, removing them both from the CRL and from storage. The
CRL will be rotated if this causes any values to be removed.`,
			},

			"safety_buffer": {
				Type: framework.TypeDurationSecond,
				Description: `The amount of extra time that must have passed
beyond certificate expiration before it is removed
from the backend storage and/or revocation list.
Defaults to 72 hours.`,
				Default: 259200, // 72h, but TypeDurationSecond currently requires defaults to be int
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTidyWrite,
		},

		HelpSynopsis:    pathTidyHelpSyn,
		HelpDescription: pathTidyHelpDesc,
	}
}

func pathTidyStatus(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy-status$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTidyStatusRead,
		},
		HelpSynopsis:    pathTidyStatusHelpSyn,
		HelpDescription: pathTidyStatusHelpDesc,
	}
}

func (b *backend) pathTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// If we are a performance standby forward the request to the active node
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	safetyBuffer := d.Get("safety_buffer").(int)
	tidyCertStore := d.Get("tidy_cert_store").(bool)
	tidyRevokedCerts := d.Get("tidy_revoked_certs").(bool)
	tidyRevocationList := d.Get("tidy_revocation_list").(bool)

	if safetyBuffer < 1 {
		return logical.ErrorResponse("safety_buffer must be greater than zero"), nil
	}

	bufferDuration := time.Duration(safetyBuffer) * time.Second

	if !atomic.CompareAndSwapUint32(b.tidyCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	// Tests using framework will screw up the storage so make a locally
	// scoped req to hold a reference
	req = &logical.Request{
		Storage: req.Storage,
	}

	go func() {
		b.tidyStatusStart(safetyBuffer, tidyCertStore, tidyRevokedCerts || tidyRevocationList)
		defer b.tidyStatusStop(nil)

		defer atomic.StoreUint32(b.tidyCASGuard, 0)

		// Don't cancel when the original client request goes away
		ctx = context.Background()

		logger := b.Logger().Named("tidy")

		doTidy := func() error {
			if tidyCertStore {
				serials, err := req.Storage.List(ctx, "certs/")
				if err != nil {
					return fmt.Errorf("error fetching list of certs: %w", err)
				}

				for i, serial := range serials {
					b.tidyStatusMessage(fmt.Sprintf("Tidying certificate store: checking entry %d of %d", i, len(serials)))
					certEntry, err := req.Storage.Get(ctx, "certs/"+serial)
					if err != nil {
						return fmt.Errorf("error fetching certificate %q: %w", serial, err)
					}

					if certEntry == nil {
						logger.Warn("certificate entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
						if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
							return fmt.Errorf("error deleting nil entry with serial %s: %w", serial, err)
						}
						b.tidyStatusIncCertStoreCount()
						continue
					}

					if certEntry.Value == nil || len(certEntry.Value) == 0 {
						logger.Warn("certificate entry has no value; tidying up since it is no longer useful for any server operations", "serial", serial)
						if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
							return fmt.Errorf("error deleting entry with nil value with serial %s: %w", serial, err)
						}
						b.tidyStatusIncCertStoreCount()
						continue
					}

					cert, err := x509.ParseCertificate(certEntry.Value)
					if err != nil {
						return fmt.Errorf("unable to parse stored certificate with serial %q: %w", serial, err)
					}

					if time.Now().After(cert.NotAfter.Add(bufferDuration)) {
						if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
							return fmt.Errorf("error deleting serial %q from storage: %w", serial, err)
						}
						b.tidyStatusIncCertStoreCount()
					}
				}
			}

			if tidyRevokedCerts || tidyRevocationList {
				b.revokeStorageLock.Lock()
				defer b.revokeStorageLock.Unlock()

				rebuildCRL := false

				revokedSerials, err := req.Storage.List(ctx, "revoked/")
				if err != nil {
					return fmt.Errorf("error fetching list of revoked certs: %w", err)
				}

				var revInfo revocationInfo
				for i, serial := range revokedSerials {
					b.tidyStatusMessage(fmt.Sprintf("Tidying revoked certificates: checking certificate %d of %d", i, len(revokedSerials)))
					revokedEntry, err := req.Storage.Get(ctx, "revoked/"+serial)
					if err != nil {
						return fmt.Errorf("unable to fetch revoked cert with serial %q: %w", serial, err)
					}

					if revokedEntry == nil {
						logger.Warn("revoked entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
						if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
							return fmt.Errorf("error deleting nil revoked entry with serial %s: %w", serial, err)
						}
						b.tidyStatusIncRevokedCertCount()
						continue
					}

					if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
						logger.Warn("revoked entry has nil value; tidying up since it is no longer useful for any server operations", "serial", serial)
						if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
							return fmt.Errorf("error deleting revoked entry with nil value with serial %s: %w", serial, err)
						}
						b.tidyStatusIncRevokedCertCount()
						continue
					}

					err = revokedEntry.DecodeJSON(&revInfo)
					if err != nil {
						return fmt.Errorf("error decoding revocation entry for serial %q: %w", serial, err)
					}

					revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
					if err != nil {
						return fmt.Errorf("unable to parse stored revoked certificate with serial %q: %w", serial, err)
					}

					// Only remove the entries from revoked/ and certs/ if we're
					// past its NotAfter value. This is because we use the
					// information on revoked/ to build the CRL and the
					// information on certs/ for lookup.
					if time.Now().After(revokedCert.NotAfter.Add(bufferDuration)) {
						if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
							return fmt.Errorf("error deleting serial %q from revoked list: %w", serial, err)
						}
						if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
							return fmt.Errorf("error deleting serial %q from store when tidying revoked: %w", serial, err)
						}
						rebuildCRL = true
						b.tidyStatusIncRevokedCertCount()
					}
				}

				if rebuildCRL {
					if err := buildCRL(ctx, b, req, false); err != nil {
						return err
					}
				}
			}

			return nil
		}

		if err := doTidy(); err != nil {
			logger.Error("error running tidy", "error", err)
			b.tidyStatusStop(err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

func (b *backend) pathTidyStatusRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.tidyStatusLock.RLock()
	defer b.tidyStatusLock.RUnlock()

	resp := &logical.Response{
		Data: map[string]interface{}{
			"safety_buffer":      nil,
			"tidy_cert_store":    nil,
			"tidy_revoked_certs": nil,

			"state":              "Inactive",
			"error":              nil,
			"time_started":       nil,
			"time_finished":      nil,
			"message":            nil,
			"cert_store_count":   nil,
			"revoked_cert_count": nil,
		},
	}

	if b.tidyStatus.state == tidyStatusInactive {
		return resp, nil
	}

	resp.Data["safety_buffer"] = b.tidyStatus.safetyBuffer
	resp.Data["tidy_cert_store"] = b.tidyStatus.tidyCertStore
	resp.Data["tidy_revoked_certs"] = b.tidyStatus.tidyRevokedCerts
	resp.Data["time_started"] = b.tidyStatus.timeStarted
	resp.Data["message"] = b.tidyStatus.message
	resp.Data["cert_store_count"] = b.tidyStatus.certStoreCount
	resp.Data["revoked_cert_count"] = b.tidyStatus.revokedCertCount

	if b.tidyStatus.state == tidyStatusStarted {
		resp.Data["state"] = "Running"
		return resp, nil
	}

	if b.tidyStatus.err == nil {
		resp.Data["state"] = "Finished"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
		resp.Data["message"] = nil
	} else {
		resp.Data["state"] = "Finished with error"
		resp.Data["error"] = b.tidyStatus.err.Error()
	}

	return resp, nil
}

func (b *backend) tidyStatusStart(safetyBuffer int, tidyCertStore, tidyRevokedCerts bool) bool {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus = &tidyStatus{
		safetyBuffer:     safetyBuffer,
		tidyCertStore:    tidyCertStore,
		tidyRevokedCerts: tidyRevokedCerts,
		state:            tidyStatusStarted,
		timeStarted:      time.Now(),
	}

	return true
}

func (b *backend) tidyStatusStop(err error) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.state = tidyStatusFinished
	b.tidyStatus.timeFinished = time.Now()
	if err != nil {
		b.tidyStatus.err = err
	}
}

func (b *backend) tidyStatusMessage(msg string) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.message = msg
}

func (b *backend) tidyStatusIncCertStoreCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.certStoreCount++
}

func (b *backend) tidyStatusIncRevokedCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.revokedCertCount++
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

const pathTidyStatusHelpSyn = `
Returns the status of the tidy operation.
`

const pathTidyStatusHelpDesc = `
This is a read only endpoint that returns information about the current tidy
operation, or the most recent if none is currently running.

The result includes the following fields:
* 'safety_buffer': the value of this parameter when initiating the tidy operation
* 'tidy_cert_store': the value of this parameter when initiating the tidy operation
* 'tidy_revoked_certs': the value of this or the tidy_revocation_list parameter
  when initiating the tidy operation
* 'state': one of "Inactive", "Running", "Finished", "Finished with error"
* 'error': the error message, if the operation ran into an error
* 'time_started': the time the operation started
* 'time_finished': the time the operation finished
* 'message': One of "Tidying certificate store: checking entry N of TOTAL" or
  "Tidying revoked certificates: checking certificate N of TOTAL"
* 'cert_store_count': The number of certificate storage entries deleted
* 'revoked_cert_count': The number of revoked certificate entries deleted
`
