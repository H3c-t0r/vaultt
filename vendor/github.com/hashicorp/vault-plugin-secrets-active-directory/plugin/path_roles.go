package plugin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	rolePath       = "roles"
	rolePrefix     = "roles/"
	roleStorageKey = "roles"

	roleCacheCleanup    = time.Second / 2
	roleCacheExpiration = time.Second
)

func (b *backend) invalidateRole(ctx context.Context, key string) {
	if strings.HasPrefix(key, rolePrefix) {
		roleName := key[len(rolePrefix):]
		b.roleCache.Delete(roleName)
	}
}

func (b *backend) pathListRoles() *framework.Path {
	return &framework.Path{
		Pattern: rolePrefix + "?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.roleListOperation,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func (b *backend) pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: rolePrefix + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
			"service_account_name": {
				Type:        framework.TypeString,
				Description: "The username/logon name for the service account with which this role will be associated.",
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "In seconds, the default password time-to-live.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.roleUpdateOperation,
			logical.ReadOperation:   b.roleReadOperation,
			logical.DeleteOperation: b.roleDeleteOperation,
		},
		HelpSynopsis:    roleHelpSynopsis,
		HelpDescription: roleHelpDescription,
	}
}

func (b *backend) readRole(ctx context.Context, storage logical.Storage, roleName string) (*backendRole, error) {
	// If it's cached, return it from there.
	roleIfc, found := b.roleCache.Get(roleName)
	if found {
		return roleIfc.(*backendRole), nil
	}

	// It's not, read it from storage.
	entry, err := storage.Get(ctx, roleStorageKey+"/"+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	role := &backendRole{}
	if err := entry.DecodeJSON(role); err != nil {
		return nil, err
	}

	// Always check when ActiveDirectory shows the password as last set on the fly.
	engineConf, err := b.readConfig(ctx, storage)
	if err != nil {
		return nil, err
	}
	if engineConf == nil {
		return nil, errors.New("the config is currently unset")
	}

	passwordLastSet, err := b.client.GetPasswordLastSet(engineConf.ADConf, role.ServiceAccountName)
	if err != nil {
		return nil, err
	}
	role.PasswordLastSet = passwordLastSet

	// Cache it.
	b.roleCache.SetDefault(roleName, role)
	return role, nil
}

func (b *backend) writeRole(ctx context.Context, storage logical.Storage, roleName string, role *backendRole) error {
	entry, err := logical.StorageEntryJSON(roleStorageKey+"/"+roleName, role)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return err
	}
	b.roleCache.SetDefault(roleName, role)
	return nil
}

func (b *backend) roleUpdateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	// Get everything we need to construct the role.
	roleName := fieldData.Get("name").(string)

	engineConf, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if engineConf == nil {
		return nil, errors.New("the config is currently unset")
	}

	// Actually construct it.
	serviceAccountName, err := getServiceAccountName(fieldData)
	if err != nil {
		return nil, err
	}

	// verify service account exists
	_, err = b.client.Get(engineConf.ADConf, serviceAccountName)
	if err != nil {
		return nil, err
	}

	ttl, err := getValidatedTTL(engineConf.PasswordConf, fieldData)
	if err != nil {
		return nil, err
	}
	role := &backendRole{
		ServiceAccountName: serviceAccountName,
		TTL:                ttl,
	}

	// Was there already a role before that we're now overwriting? If so, let's carry forward the LastVaultRotation.
	oldRole, err := b.readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	} else {
		if oldRole != nil {
			role.LastVaultRotation = oldRole.LastVaultRotation
		}
	}

	// writeRole it to storage and the roleCache.
	if err := b.writeRole(ctx, req.Storage, roleName, role); err != nil {
		return nil, err
	}

	// Return a 204.
	return nil, nil
}

func (b *backend) roleReadOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	roleName := fieldData.Get("name").(string)

	role, err := b.readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

func (b *backend) roleListOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	keys, err := req.Storage.List(ctx, roleStorageKey+"/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

func (b *backend) roleDeleteOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	roleName := fieldData.Get("name").(string)

	if err := req.Storage.Delete(ctx, roleStorageKey+"/"+roleName); err != nil {
		return nil, err
	}

	b.roleCache.Delete(roleName)

	if err := b.deleteCred(ctx, req.Storage, roleName); err != nil {
		return nil, err
	}
	return nil, nil
}

func getServiceAccountName(fieldData *framework.FieldData) (string, error) {
	serviceAccountName := fieldData.Get("service_account_name").(string)
	if serviceAccountName == "" {
		return "", errors.New("\"service_account_name\" is required")
	}
	return serviceAccountName, nil
}

func getValidatedTTL(passwordConf *passwordConf, fieldData *framework.FieldData) (int, error) {
	ttl := fieldData.Get("ttl").(int)
	if ttl == 0 {
		ttl = passwordConf.TTL
	}
	if ttl > passwordConf.MaxTTL {
		return 0, fmt.Errorf("requested ttl of %d seconds is over the max ttl of %d seconds", ttl, passwordConf.MaxTTL)
	}
	if ttl < 0 {
		return 0, fmt.Errorf("ttl can't be negative")
	}
	return ttl, nil
}

const (
	roleHelpSynopsis    = ``
	roleHelpDescription = ``

	pathListRolesHelpSyn  = ``
	pathListRolesHelpDesc = ``
)
