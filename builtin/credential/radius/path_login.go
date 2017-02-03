package radius

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login/(?P<username>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username to be used for login.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	policies, resp, err := b.Login(req, username, password)
	// Handle an internal error
	if err != nil {
		return nil, err
	}
	if resp != nil {
		// Handle a logical error
		if resp.IsError() {
			return resp, nil
		}
	} else {
		resp = &logical.Response{}
	}

	sort.Strings(policies)

	resp.Auth = &logical.Auth{
		Policies: policies,
		Metadata: map[string]string{
			"username": username,
			"policies": strings.Join(policies, ","),
		},
		InternalData: map[string]interface{}{
			"password": password,
		},
		DisplayName: username,
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
	}
	return resp, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var err error

	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	var resp *logical.Response
	var loginPolicies []string
	var user *UserEntry
	if cfg.ReauthOnRenew {
		loginPolicies, resp, err = b.Login(req, username, password)
		if len(loginPolicies) == 0 {
			return resp, err
		}
	} else {
		user, err = b.user(req.Storage, username)
		if err != nil {
			return nil, err
		}
		if user == nil {
			// User no longer exists, do not renew
			return nil, fmt.Errorf("Renew not allowed for unknown users")
		}
		loginPolicies = user.Policies
	}

	if !policyutil.EquivalentPolicies(loginPolicies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	return framework.LeaseExtend(0, 0, b.System())(req, d)
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password. Please be sure to
read the note on escaping from the path-help for the 'config' endpoint.
`
