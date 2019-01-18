package jwtauth

import (
	"context"
	"sync"

	oidc "github.com/coreos/go-oidc"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	configPath string = "config"
	rolePrefix string = "role/"
)

// Factory is used by framework
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

type jwtAuthBackend struct {
	*framework.Backend

	l            sync.RWMutex
	provider     *oidc.Provider
	cachedConfig *jwtConfig

	providerCtx       context.Context
	providerCtxCancel context.CancelFunc
}

func backend(c *logical.BackendConfig) *jwtAuthBackend {
	b := new(jwtAuthBackend)
	b.providerCtx, b.providerCtxCancel = context.WithCancel(context.Background())

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Invalidate:  b.invalidate,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathLogin(b),
				pathRoleList(b),
				pathRole(b),
				pathConfig(b),
			},
		),
		Clean: b.cleanup,
	}

	return b
}

func (b *jwtAuthBackend) cleanup(_ context.Context) {
	b.l.Lock()
	if b.providerCtxCancel != nil {
		b.providerCtxCancel()
	}
	b.l.Unlock()
}

func (b *jwtAuthBackend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.reset()
	}
}

func (b *jwtAuthBackend) reset() {
	b.l.Lock()
	b.provider = nil
	b.cachedConfig = nil
	b.l.Unlock()
}

func (b *jwtAuthBackend) getProvider(ctx context.Context, config *jwtConfig) (*oidc.Provider, error) {
	b.l.RLock()
	unlockFunc := b.l.RUnlock
	defer func() { unlockFunc() }()

	if b.provider != nil {
		return b.provider, nil
	}

	b.l.RUnlock()
	b.l.Lock()
	unlockFunc = b.l.Unlock

	if b.provider != nil {
		return b.provider, nil
	}

	provider, err := b.createProvider(config)
	if err != nil {
		return nil, err
	}

	b.provider = provider
	return provider, nil
}

const (
	backendHelp = `
The JWT backend plugin allows authentication using JWTs (including OIDC).
`
)
