package mfa

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func methodIdentifiersListPaths(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "methods/" + framework.GenericNameRegex("method_name") + "/?$",

		Fields: map[string]*framework.FieldSchema{
			"method_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: mfaMethodNameHelp,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.mfaBackendMethodIdentifiersList,
		},

		HelpSynopsis:    mfaListMethodsHelp,
		HelpDescription: mfaListMethodsHelp,
	}
}

func methodIdentifiersPaths(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "methods/" + framework.GenericNameRegex("method_name") + "/(?P<identifier>.+)",

		Fields: map[string]*framework.FieldSchema{
			"method_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: mfaMethodNameHelp,
			},

			"identifier": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: mfaTypesHelp,
			},

			"totp_account_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: mfaTypesHelp,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.mfaBackendMethodIdentifiersRead,
			logical.CreateOperation: b.mfaBackendMethodIdentifiersCreate,
			logical.DeleteOperation: b.mfaBackendMethodIdentifiersDelete,
		},

		ExistenceCheck: b.mfaBackendMethodIdentifiersExistenceCheck,
	}
}

type mfaIdentifierEntry struct {
	CreationTime    time.Time `json:"creation_time" mapstructure:"creation_time" structs:"creation_time"`
	TOTPURL         string    `json:"totp_url" mapstructure:"totp_url" structs:"totp_url"`
	TOTPSecret      string    `json:"totp_secret" mapstructure:"totp_secret" structs:"totp_secret"`
	TOTPAccountName string    `json:"totp_account_name" mapstructure:"totp_account_name" structs:"totp_account_name"`
	Identifier      string    `json:"identifier" mapstructure:"identifier" structs:"identifier"`
}

func (b *backend) mfaBackendMethodIdentifiers(methodName, identifier string) (*mfaMethodEntry, *mfaIdentifierEntry, error) {
	b.RLock()
	defer b.RUnlock()

	return b.mfaBackendMethodIdentifiersInternal(methodName, identifier)
}

func (b *backend) mfaBackendMethodIdentifiersInternal(methodName, identifier string) (*mfaMethodEntry, *mfaIdentifierEntry, error) {
	method, err := b.mfaBackendMethodInternal(methodName)
	if err != nil {
		return nil, nil, err
	}
	if method == nil {
		return nil, nil, fmt.Errorf("method %s does not exist", methodName)
	}

	entry, err := b.storage.Get(fmt.Sprintf("method/%s/identifiers/%s", methodName, strings.ToLower(identifier)))
	if err != nil {
		return nil, nil, err
	}
	if entry == nil {
		return method, nil, nil
	}

	var result mfaIdentifierEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, nil, err
	}

	return method, &result, nil
}

func (b *backend) mfaBackendMethodIdentifiersList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	b.RLock()
	defer b.RUnlock()

	method, err := b.mfaBackendMethodInternal(methodName)
	if err != nil {
		return nil, err
	}
	if method == nil {
		return logical.ErrorResponse(fmt.Sprintf("method %s doesn't exist", methodName)), nil
	}

	entries, err := b.storage.List(fmt.Sprintf("method/%s/identifiers/", methodName))
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(entries))
	for i, entry := range entries {
		ret[i] = strings.TrimPrefix(entry, fmt.Sprintf("method/%s/identifiers/", methodName))
	}

	return logical.ListResponse(ret), nil
}

func (b *backend) mfaBackendMethodIdentifiersDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	identifier := data.Get("identifier").(string)
	if identifier == "" {
		return logical.ErrorResponse("identifier cannot be empty"), nil
	}

	b.Lock()
	defer b.Unlock()

	err := b.storage.Delete(fmt.Sprintf("method/%s/identifiers/%s", methodName, strings.ToLower(identifier)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) mfaBackendMethodIdentifiersRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	identifier := data.Get("identifier").(string)
	if identifier == "" {
		return logical.ErrorResponse("identifier cannot be empty"), nil
	}

	_, entry, err := b.mfaBackendMethodIdentifiersInternal(methodName, identifier)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"creation_time_utc":    entry.CreationTime.Unix(),
			"creation_time_string": entry.CreationTime.String(),
			"totp_account_name":    entry.TOTPAccountName,
			"identifier":           entry.Identifier,
		},
	}

	return resp, nil
}

func (b *backend) mfaBackendMethodIdentifiersExistenceCheck(
	req *logical.Request, data *framework.FieldData) (bool, error) {
	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return false, fmt.Errorf("method name cannot be empty")
	}

	identifier := data.Get("identifier").(string)
	if identifier == "" {
		return false, fmt.Errorf("identifier cannot be empty")
	}

	_, entry, err := b.mfaBackendMethodIdentifiers(methodName, identifier)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (b *backend) mfaBackendMethodIdentifiersCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.Operation == logical.UpdateOperation {
		return logical.ErrorResponse("identifiers cannot be updated; delete and recreate the identifier to reregister"), nil
	}

	methodName := data.Get("method_name").(string)
	if methodName == "" {
		return logical.ErrorResponse("method name cannot be empty"), nil
	}

	identifier := data.Get("identifier").(string)
	if identifier == "" {
		return logical.ErrorResponse("identifier cannot be empty"), nil
	}

	totpAccountName := data.Get("totp_account_name").(string)

	b.Lock()
	defer b.Unlock()

	method, err := b.mfaBackendMethodInternal(methodName)
	if err != nil {
		return nil, err
	}
	if method == nil {
		return logical.ErrorResponse(fmt.Sprintf("method %s does not exist", methodName)), nil
	}

	entry := &mfaIdentifierEntry{
		CreationTime:    time.Now().UTC(),
		TOTPAccountName: totpAccountName,
		Identifier:      identifier,
	}

	resp := &logical.Response{}

	switch method.Type {
	case "totp":
		err := b.createTOTPKey(method, entry, resp)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return resp, nil
		}

	default:
		return logical.ErrorResponse(fmt.Sprintf("identifier registration not supported for method type %s", method.Type)), nil
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON(fmt.Sprintf("method/%s/identifiers/%s", methodName, strings.ToLower(identifier)), entry)
	if err != nil {
		return nil, err
	}
	if err := b.storage.Put(jsonEntry); err != nil {
		return nil, err
	}

	return resp, nil
}
