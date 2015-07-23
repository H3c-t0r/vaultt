package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config",
				"groups/*",
				"users/*",
			},

			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathLogin(&b),
			pathConfig(&b),
			pathGroups(&b),
			pathUsers(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

func EscapeLDAPValue(input string) string {
	// RFC4514 forbids un-escaped:
	// - leading space or hash
	// - trailing space
	// - special characters '"', '+', ',', ';', '<', '>', '\\'
	// - null
	for i := 0; i < len(input); i++ {
		escaped := false
		if input[i] == '\\' {
			i++
			escaped = true
		}
		switch input[i] {
		case '"', '+', ',', ';', '<', '>', '\\':
			if !escaped {
				input = input[0:i] + "\\" + input[i:]
				i++
			}
			continue
		}
		if escaped {
			input = input[0:i] + "\\" + input[i:]
			i++
		}
	}
	if input[0] == ' ' || input[0] == '#' {
		input = "\\" + input
	}
	if input[len(input)-1] == ' ' {
		input = input[0:len(input)-1] + "\\ "
	}
	return input
}

func (b *backend) Login(req *logical.Request, username string, password string) ([]string, *logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("ldap backend not configured"), nil
	}

	c, err := cfg.DialLDAP()
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}

	binddn := ""
	if cfg.UPNDomain != "" {
		binddn = fmt.Sprintf("%s@%s", EscapeLDAPValue(username), cfg.UPNDomain)
	} else {
		binddn = fmt.Sprintf("%s=%s,%s", cfg.UserAttr, EscapeLDAPValue(username), cfg.UserDN)
	}

	var binduser string
	var bindpassword string
	configuredBinding := false

	if cfg.BindDN == "" { //Bind with requesting credentials
		binduser = binddn
		bindpassword = password
	} else { //Bind with configured credentials
		binduser = cfg.BindDN
		bindpassword = cfg.BindDNPassword
		configuredBinding = true
	}

	//Attempt bind
	if err = c.Bind(binduser, bindpassword); err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("LDAP bind (%s) failed: %v", binduser, err)), nil
	}

	//If using cofigured binding credentials, we need to make sure
	//the requesting credentials are correct
	if configuredBinding {
		filter := fmt.Sprintf("(%s=%s)", cfg.UserAttr, username)
		//Find the user requesting to login
		sresult, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.UserDN,
			Scope:  2, // subtree
			Filter: filter,
		})
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP user search (%s) failed: %v", filter, err)), nil
		}

		//Requesting user wasn't found
		if len(sresult.Entries) == 0 {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP user search (%s) failed to find user: %s", filter, username)), nil
		}

		discoveredUserDN := sresult.Entries[0].DN

		//Check requesting credentials
		if err = c.Bind(discoveredUserDN, password); err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP bind (%s) failed: %v", discoveredUserDN, err)), nil
		}

		//User the user dn to perform group searches
		binddn = discoveredUserDN
	}

	userdn := ""
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		sresult, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.UserDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(userPrincipalName=%s)", binddn),
		})
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP search failed: %v", err)), nil
		}
		for _, e := range sresult.Entries {
			userdn = e.DN
		}
	} else {
		userdn = binddn
	}

	// Enumerate all groups the user is member of. The search filter should
	// work with both openldap and MS AD standard schemas.
	sresult, err := c.Search(&ldap.SearchRequest{
		BaseDN: cfg.GroupDN,
		Scope:  2, // subtree
		Filter: fmt.Sprintf("(|(memberUid=%s)(member=%s)(uniqueMember=%s))", username, userdn, userdn),
	})
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("LDAP search failed: %v", err)), nil
	}

	var allgroups []string
	var policies []string

	user, err := b.User(req.Storage, username)
	if err == nil && user != nil {
		allgroups = append(allgroups, user.Groups...)
	}

	for _, e := range sresult.Entries {
		dn, err := ldap.ParseDN(e.DN)
		if err != nil || len(dn.RDNs) == 0 || len(dn.RDNs[0].Attributes) == 0 {
			continue
		}
		gname := dn.RDNs[0].Attributes[0].Value
		allgroups = append(allgroups, gname)
	}

	for _, gname := range allgroups {
		group, err := b.Group(req.Storage, gname)
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}

	if len(policies) == 0 {
		return nil, logical.ErrorResponse(fmt.Sprintf("user (%s) is not member of any authorized group", binddn)), nil
	}

	return policies, nil, nil
}

const backendHelp = `
The "ldap" credential provider allows authentication querying
a LDAP server, checking username and password, and associating groups
to set of policies.

Configuration of the server is done through the "config" and "groups"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
