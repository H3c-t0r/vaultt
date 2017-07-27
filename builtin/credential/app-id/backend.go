package appId

import (
	"sync"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(conf *logical.BackendConfig) (*backend, error) {
	var b backend
	b.MapAppId = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "app-id",
			Schema: map[string]*framework.FieldSchema{
				"display_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "A name to map to this app ID for logs.",
				},

				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Policies for the app ID.",
				},
			},
		},
		DefaultKey: "default",
	}

	b.MapUserId = &framework.PathMap{
		Name: "user-id",
		Schema: map[string]*framework.FieldSchema{
			"cidr_block": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "If not blank, restricts auth by this CIDR block",
			},

			"value": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "App IDs that this user associates with.",
			},
		},
	}

	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
				"login/*",
			},
		},

		Paths: framework.PathAppend([]*framework.Path{
			pathLogin(&b),
			pathLoginWithAppIDPath(&b),
		},
			b.MapAppId.Paths(),
			b.MapUserId.Paths(),
		),

		AuthRenew: b.pathLoginRenew,

		Invalidate: b.invalidate,
	}

	b.view = conf.StorageView
	b.MapAppId.SaltFunc = b.Salt
	b.MapUserId.SaltFunc = b.Salt

	return &b, nil
}

type backend struct {
	*framework.Backend

	salt      *salt.Salt
	SaltMutex sync.RWMutex
	view      logical.Storage
	MapAppId  *framework.PolicyMap
	MapUserId *framework.PathMap
}

func (b *backend) Salt() (*salt.Salt, error) {
	b.SaltMutex.RLock()
	if b.salt != nil {
		defer b.SaltMutex.RUnlock()
		return b.salt, nil
	}
	b.SaltMutex.RUnlock()
	b.SaltMutex.Lock()
	defer b.SaltMutex.Unlock()
	if b.salt != nil {
		return b.salt, nil
	}
	salt, err := salt.NewSalt(b.view, &salt.Config{
		HashFunc: salt.SHA1Hash,
		Location: salt.DefaultLocation,
	})
	if err != nil {
		return nil, err
	}
	b.salt = salt
	return salt, nil
}

func (b *backend) invalidate(key string) {
	switch key {
	case salt.DefaultLocation:
		b.SaltMutex.Lock()
		defer b.SaltMutex.Unlock()
		b.salt = nil
	}
}

const backendHelp = `
The App ID credential provider is used to perform authentication from
within applications or machine by pairing together two hard-to-guess
unique pieces of information: a unique app ID, and a unique user ID.

The goal of this credential provider is to allow elastic users
(dynamic machines, containers, etc.) to authenticate with Vault without
having to store passwords outside of Vault. It is a single method of
solving the chicken-and-egg problem of setting up Vault access on a machine.
With this provider, nobody except the machine itself has access to both
pieces of information necessary to authenticate. For example:
configuration management will have the app IDs, but the machine itself
will detect its user ID based on some unique machine property such as a
MAC address (or a hash of it with some salt).

An example, real world process for using this provider:

  1. Create unique app IDs (UUIDs work well) and map them to policies.
     (Path: map/app-id/<app-id>)

  2. Store the app IDs within configuration management systems.

  3. An out-of-band process run by security operators map unique user IDs
     to these app IDs. Example: when an instance is launched, a cloud-init
     system tells security operators a unique ID for this machine. This
     process can be scripted, but the key is that it is out-of-band and
     out of reach of configuration management.
	 (Path: map/user-id/<user-id>)

  4. A new server is provisioned. Configuration management configures the
     app ID, the server itself detects its user ID. With both of these
     pieces of information, Vault can be accessed according to the policy
     set by the app ID.

More details on this process follow:

The app ID is a unique ID that maps to a set of policies. This ID is
generated by an operator and configured into the backend. The ID itself
is usually a UUID, but any hard-to-guess unique value can be used.

After creating app IDs, an operator authorizes a fixed set of user IDs
with each app ID. When a valid {app ID, user ID} tuple is given to the
"login" path, then the user is authenticated with the configured app
ID policies.

The user ID can be any value (just like the app ID), however it is
generally a value unique to a machine, such as a MAC address or instance ID,
or a value hashed from these unique values.

It is possible to authorize multiple app IDs with each
user ID by writing them as comma-separated values to the map/user-id/<user-id>
path.

It is also possible to renew the auth tokens with 'vault token-renew <token>' command.
Before the token is renewed, the validity of app ID, user ID and the associated
policies are checked again.
`
