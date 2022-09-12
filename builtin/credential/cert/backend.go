package cert

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/ocsp"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	bConf, err := b.Config(ctx, conf.StorageView)
	if err != nil {
		return nil, err
	}
	if conf != nil {
		b.initOCSPClient(bConf.OcspCacheSize)
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathConfig(&b),
			pathLogin(&b),
			pathListCerts(&b),
			pathCerts(&b),
			pathCRLs(&b),
		},
		AuthRenew:   b.pathLoginRenew,
		Invalidate:  b.invalidate,
		BackendType: logical.TypeCredential,
	}

	b.crlUpdateMutex = &sync.RWMutex{}
	return &b
}

type backend struct {
	*framework.Backend
	MapCertId *framework.PathMap

	crls            map[string]CRLInfo
	ocspDisabled    bool
	crlUpdateMutex  *sync.RWMutex
	ocspClientMutex sync.RWMutex
	ocspClient      *ocsp.Client
}

func (b *backend) invalidate(_ context.Context, key string) {
	switch {
	case strings.HasPrefix(key, "crls/"):
		b.crlUpdateMutex.Lock()
		defer b.crlUpdateMutex.Unlock()
		b.crls = nil
	}
}

func (b *backend) initOCSPClient(cacheSize int) {
	b.ocspClientMutex.Lock()
	defer b.ocspClientMutex.Unlock()
	b.ocspClient = ocsp.New(func() hclog.Logger {
		return b.Logger()
	}, cacheSize)
}

const backendHelp = `
The "cert" credential provider allows authentication using
TLS client certificates. A client connects to Vault and uses
the "login" endpoint to generate a client token.

Trusted certificates are configured using the "certs/" endpoint
by a user with root access. A certificate authority can be trusted,
which permits all keys signed by it. Alternatively, self-signed
certificates can be trusted avoiding the need for a CA.
`
