package pki

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	keyConfig    = "/config/keys"
	issuerConfig = "/config/issuers"
	keyPrefix    = "/config/key/"
	issuerPrefix = "/config/issuer/"
)

type keyId string

func (p keyId) String() string {
	return string(p)
}

type issuerId string

func (p issuerId) String() string {
	return string(p)
}

type key struct {
	ID             keyId                   `json:"id" structs:"id" mapstructure:"id"`
	Name           string                  `json:"name" structs:"name" mapstructure:"name"`
	PrivateKeyType certutil.PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	PrivateKey     string                  `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

type issuer struct {
	ID           issuerId `json:"id" structs:"id" mapstructure:"id"`
	Name         string   `json:"name" structs:"name" mapstructure:"name"`
	KeyID        keyId    `json:"key_id" structs:"key_id" mapstructure:"key_id"`
	Certificate  string   `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	CAChain      []string `json:"ca_chain" structs:"ca_chain" mapstructure:"ca_chain"`
	SerialNumber string   `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
}

func listKeys(ctx context.Context, s logical.Storage) ([]keyId, error) {
	strList, err := s.List(ctx, keyPrefix)
	if err != nil {
		return nil, err
	}

	keyIds := make([]keyId, 0, len(strList))
	for _, entry := range strList {
		keyIds = append(keyIds, keyId(entry))
	}

	return keyIds, nil
}

func fetchKeyById(ctx context.Context, s logical.Storage, keyId keyId) (*key, error) {
	keyEntry, err := s.Get(ctx, keyPrefix+keyId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: %v", err)}
	}
	if keyEntry == nil {
		// FIXME: Dedicated/specific error for this?
		return nil, errutil.UserError{Err: fmt.Sprintf("pki key id %s does not exist", keyId.String())}
	}

	var key key
	if err := keyEntry.DecodeJSON(&key); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki key with id %s: %v", keyId.String(), err)}
	}

	return &key, nil
}

func writeKey(ctx context.Context, s logical.Storage, key key) error {
	keyId := key.ID

	json, err := logical.StorageEntryJSON(keyPrefix+keyId.String(), key)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func listIssuers(ctx context.Context, s logical.Storage) ([]issuerId, error) {
	strList, err := s.List(ctx, issuerPrefix)
	if err != nil {
		return nil, err
	}

	issuerIds := make([]issuerId, 0, len(strList))
	for _, entry := range strList {
		issuerIds = append(issuerIds, issuerId(entry))
	}

	return issuerIds, nil
}

func resolveKeyReference(ctx context.Context, s logical.Storage, reference string) (keyId, error) {
	if reference == "default" {
		// Handle fetching the default key.
		config, err := getKeysConfig(ctx, s)
		if err != nil {
			return keyId("config-error"), err
		}

		return keyId(config["default"]), nil
	}

	keys, err := listKeys(ctx, s)
	if err != nil {
		return keyId("list-error"), err
	}

	// Cheaper to list keys and check if an id is a match...
	for _, key_id := range keys {
		if key_id == keyId(reference) {
			return key_id, nil
		}
	}

	// ... than to pull all keys from storage.
	for _, key_id := range keys {
		key, err := fetchKeyById(ctx, s, key_id)
		if err != nil {
			return keyId("key-read"), err
		}

		if key.Name == reference {
			return key.ID, nil
		}
	}

	// Otherwise, we must not have found the key.
	return keyId("not-found"), errutil.UserError{Err: fmt.Sprintf("unable to find PKI key for reference: %v", reference)}
}

func fetchIssuerById(ctx context.Context, s logical.Storage, issuerId issuerId) (*issuer, error) {
	issuerEntry, err := s.Get(ctx, issuerPrefix+issuerId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
	}
	if issuerEntry == nil {
		// FIXME: Dedicated/specific error for this?
		return nil, errutil.UserError{Err: fmt.Sprintf("pki issuer id %s does not exist", issuerId.String())}
	}

	var issuer issuer
	if err := issuerEntry.DecodeJSON(&issuer); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki issuer with id %s: %v", issuerId.String(), err)}
	}

	return &issuer, nil
}

func writeIssuer(ctx context.Context, s logical.Storage, issuer issuer) error {
	issuerId := issuer.ID

	json, err := logical.StorageEntryJSON(issuerPrefix+issuerId.String(), issuer)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func setKeysConfig(ctx context.Context, s logical.Storage, config map[string]keyId) error {
	json, err := logical.StorageEntryJSON(keyConfig, config)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getKeysConfig(ctx context.Context, s logical.Storage) (map[string]keyId, error) {
	keyConfigEntry, err := s.Get(ctx, keyConfig)
	if err != nil {
		return nil, err
	}

	var keyConfig map[string]keyId
	if err := keyConfigEntry.DecodeJSON(&keyConfig); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode key configuration: %v", err)}
	}

	return keyConfig, nil
}

func setIssuersConfig(ctx context.Context, s logical.Storage, config map[string]issuerId) error {
	json, err := logical.StorageEntryJSON(issuerConfig, config)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getIssuersConfig(ctx context.Context, s logical.Storage) (map[string]issuerId, error) {
	issuerConfigEntry, err := s.Get(ctx, issuerConfig)
	if err != nil {
		return nil, err
	}

	var issuerConfig map[string]issuerId
	if err := issuerConfigEntry.DecodeJSON(&issuerConfig); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode issuer configuration: %v", err)}
	}

	return issuerConfig, nil
}

func resolveIssuerReference(ctx context.Context, s logical.Storage, reference string) (issuerId, error) {
	if reference == "default" {
		// Handle fetching the default issuer.
		config, err := getIssuersConfig(ctx, s)
		if err != nil {
			return issuerId("config-error"), err
		}

		return issuerId(config["default"]), nil
	}

	issuers, err := listIssuers(ctx, s)
	if err != nil {
		return issuerId("list-error"), err
	}

	// Cheaper to list issuers and check if an id is a match...
	for _, issuer_id := range issuers {
		if issuer_id == issuerId(reference) {
			return issuer_id, nil
		}
	}

	// ... than to pull all issuers from storage.
	for _, issuer_id := range issuers {
		issuer, err := fetchIssuerById(ctx, s, issuer_id)
		if err != nil {
			return issuerId("issuer-read"), err
		}

		if issuer.Name == reference {
			return issuer.ID, nil
		}
	}

	// Otherwise, we must not have found the issuer.
	return issuerId("not-found"), errutil.UserError{Err: fmt.Sprintf("unable to find PKI issuer for reference: %v", reference)}
}
