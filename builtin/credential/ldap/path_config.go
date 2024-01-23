// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const userFilterWarning = "userfilter configured does not consider userattr and may result in colliding entity aliases on logins"

func pathConfig(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: `config`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixLDAP,
			Action:          "Configure",
		},

		Fields: ldaputil.ConfigFields(),

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "auth-configuration",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure-auth",
				},
			},
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}

	tokenutil.AddTokenFields(p.Fields)
	p.Fields["token_policies"].Description += ". This will apply to all tokens generated by this auth method, in addition to any configured for specific users/groups."

	p.Fields["password_policy"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "Password policy to use to rotate the root password",
	}

	p.Fields["rotation_schedule"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "CRON-style string that will define the schedule on which rotations should occur. Mutually exclusive with TTL",
	}

	p.Fields["rotation_window"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Description: "Specifies the amount of time in which the rotation is allowed to occur starting from a given rotation_schedule",
	}

	p.Fields["ttl"] = &framework.FieldSchema{
		Type:        framework.TypeInt,
		Description: "TTL for automatic credential rotation of the given username. Mutually exclusive with rotation_schedule",
	}

	return p
}

/*
 * Construct ConfigEntry struct using stored configuration.
 */
func (b *backend) Config(ctx context.Context, req *logical.Request) (*ldapConfigEntry, error) {
	storedConfig, err := req.Storage.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if storedConfig == nil {
		// Create a new ConfigEntry, filling in defaults where appropriate
		fd, err := b.getConfigFieldData()
		if err != nil {
			return nil, err
		}

		result, err := ldaputil.NewConfigEntry(nil, fd)
		if err != nil {
			return nil, err
		}

		// No user overrides, return default configuration
		result.CaseSensitiveNames = new(bool)
		*result.CaseSensitiveNames = false

		result.UsePre111GroupCNBehavior = new(bool)
		*result.UsePre111GroupCNBehavior = false

		return &ldapConfigEntry{ConfigEntry: result}, nil
	}

	// Deserialize stored configuration.
	// Fields not specified in storedConfig will retain their defaults.
	result := new(ldapConfigEntry)
	result.ConfigEntry = new(ldaputil.ConfigEntry)
	if err := storedConfig.DecodeJSON(result); err != nil {
		return nil, err
	}

	var persistNeeded bool
	if result.CaseSensitiveNames == nil {
		// Upgrade from before switching to case-insensitive
		result.CaseSensitiveNames = new(bool)
		*result.CaseSensitiveNames = true
		persistNeeded = true
	}

	if result.UsePre111GroupCNBehavior == nil {
		result.UsePre111GroupCNBehavior = new(bool)
		*result.UsePre111GroupCNBehavior = true
		persistNeeded = true
	}

	// leave these blank if unset, which would mean no rotation at all
	//if result.RotationSchedule == "" {
	//	result.RotationSchedule = "0 0 0 0 0"
	//}

	//if result.RotationWindow == 0 {
	//	// default rotation windoe
	//}

	if persistNeeded && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby)) {
		entry, err := logical.StorageEntryJSON("config", result)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	data := cfg.PasswordlessMap()
	cfg.PopulateTokenData(data)
	data["password_policy"] = cfg.PasswordPolicy

	resp := &logical.Response{
		Data: data,
	}

	if warnings := b.checkConfigUserFilter(cfg); len(warnings) > 0 {
		resp.Warnings = warnings
	}

	return resp, nil
}

// checkConfigUserFilter performs a best-effort check the config's userfilter.
// It will checked whether the templated or literal userattr value is present,
// and if not return a warning.
func (b *backend) checkConfigUserFilter(cfg *ldapConfigEntry) []string {
	if cfg == nil || cfg.UserFilter == "" {
		return nil
	}

	var warnings []string

	switch {
	case strings.Contains(cfg.UserFilter, "{{.UserAttr}}"):
		// Case where the templated userattr value is provided
	case strings.Contains(cfg.UserFilter, cfg.UserAttr):
		// Case where the literal userattr value is provided
	default:
		b.Logger().Debug(userFilterWarning, "userfilter", cfg.UserFilter, "userattr", cfg.UserAttr)
		warnings = append(warnings, userFilterWarning)
	}

	return warnings
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	// Build a ConfigEntry struct out of the supplied FieldData
	cfg.ConfigEntry, err = ldaputil.NewConfigEntry(cfg.ConfigEntry, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// On write, if not specified, use false. We do this here so upgrade logic
	// works since it calls the same newConfigEntry function
	if cfg.CaseSensitiveNames == nil {
		cfg.CaseSensitiveNames = new(bool)
		*cfg.CaseSensitiveNames = false
	}

	if cfg.UsePre111GroupCNBehavior == nil {
		cfg.UsePre111GroupCNBehavior = new(bool)
		*cfg.UsePre111GroupCNBehavior = false
	}

	if err := cfg.ParseTokenFields(req, d); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	if passwordPolicy, ok := d.GetOk("password_policy"); ok {
		cfg.PasswordPolicy = passwordPolicy.(string)
	}

	ttl, ttlOk := d.GetOk("ttl")
	rotationSchedule, rotationScheduleOk := d.GetOk("rotation_schedule")
	rotationWindow, rotationWindowOk := d.GetOk("rotation_window")

	var rc *logical.RootCredential
	if rotationScheduleOk && ttlOk {
		return logical.ErrorResponse("mutually exclusive fields rotation_schedule and ttl were both specified; only one of them can be provided"), nil
	} else if rotationWindowOk && ttlOk {
		return logical.ErrorResponse("rotation_window does not apply to ttl"), nil
	} else if rotationScheduleOk && !rotationWindowOk || rotationWindowOk && !rotationScheduleOk {
		return logical.ErrorResponse("must include both rotation_schedule and rotation_window"), nil
	}

	if rotationScheduleOk && rotationWindowOk {
		cfg.RotationSchedule = rotationSchedule.(string)
		cfg.RotationWindow = rotationWindow.(int)

		rc, err = logical.GetRootCredential(cfg.RotationSchedule, "ldap/config",
			"ldap-root-creds", cfg.RotationWindow, 0)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
		// unset ttl if rotation_schedule is set since these are mutually exclusive
		cfg.TTL = 0

		b.Logger().Info("rotation", "window", cfg.RotationWindow, "schedule", cfg.RotationSchedule, "ttl", cfg.TTL)
	}

	if ttlOk {
		cfg.TTL = ttl.(int)

		rc, err = logical.GetRootCredential("", "ldap/config",
			"ldap-root-creds", 0, cfg.TTL)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		cfg.RotationSchedule = ""
		cfg.RotationWindow = 0

		b.Logger().Info("rotation", "window", cfg.RotationWindow, "schedule", cfg.RotationSchedule, "ttl", cfg.TTL)
	}

	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if warnings := b.checkConfigUserFilter(cfg); len(warnings) > 0 {
		return &logical.Response{
			Warnings: warnings,
		}, nil
	}

	if rc != nil {
		return &logical.Response{
			RootCredential: rc,
		}, nil
	} else {
		return nil, nil
	}
}

/*
 * Returns FieldData describing our ConfigEntry struct schema
 */
func (b *backend) getConfigFieldData() (*framework.FieldData, error) {
	configPath := b.Route("config")

	if configPath == nil {
		return nil, logical.ErrUnsupportedPath
	}

	raw := make(map[string]interface{}, len(configPath.Fields))

	fd := framework.FieldData{
		Raw:    raw,
		Schema: configPath.Fields,
	}

	return &fd, nil
}

type ldapConfigEntry struct {
	tokenutil.TokenParams
	*ldaputil.ConfigEntry

	PasswordPolicy   string `json:"password_policy"`
	RotationSchedule string `json:"rotation_schedule"`
	RotationWindow   int    `json:"rotation_window"`
	TTL              int    `json:"ttl"`
}

const pathConfigHelpSyn = `
Configure the LDAP server to connect to, along with its options.
`

const pathConfigHelpDesc = `
This endpoint allows you to configure the LDAP server to connect to and its
configuration options.

The LDAP URL can use either the "ldap://" or "ldaps://" schema. In the former
case, an unencrypted connection will be made with a default port of 389, unless
the "starttls" parameter is set to true, in which case TLS will be used. In the
latter case, a SSL connection will be established with a default port of 636.

## A NOTE ON ESCAPING

It is up to the administrator to provide properly escaped DNs. This includes
the user DN, bind DN for search, and so on.

The only DN escaping performed by this backend is on usernames given at login
time when they are inserted into the final bind DN, and uses escaping rules
defined in RFC 4514.

Additionally, Active Directory has escaping rules that differ slightly from the
RFC; in particular it requires escaping of '#' regardless of position in the DN
(the RFC only requires it to be escaped when it is the first character), and
'=', which the RFC indicates can be escaped with a backslash, but does not
contain in its set of required escapes. If you are using Active Directory and
these appear in your usernames, please ensure that they are escaped, in
addition to being properly escaped in your configured DNs.

For reference, see https://www.ietf.org/rfc/rfc4514.txt and
http://social.technet.microsoft.com/wiki/contents/articles/5312.active-directory-characters-to-escape.aspx
`
