package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

// A kvv1 client is used to perform reads and writes against a KV v1 secrets engine in Vault.
//
// The mount path is the location where the target KV secrets engine resides
// in Vault.
//
// While v1 is not necessarily deprecated, Vault development servers tend to
// use v2 as the version of the KV secrets engine, as this is what's mounted
// by default when a server is started in -dev mode. See the kvv2 struct.
//
// Learn more about the KV secrets engine here:
// https://www.vaultproject.io/docs/secrets/kv
type kvv1 struct {
	c         *Client
	mountPath string
}

// A kvv2 client is used to perform reads and writes against a KV v2 secrets engine in Vault.
//
// The mount path is the location where the target KV secrets engine resides
// in Vault.
//
// Vault development servers tend to have "secret" as the mount path,
// as these are the default settings when a server is started in -dev mode.
//
// Learn more about the KV secrets engine here:
// https://www.vaultproject.io/docs/secrets/kv
type kvv2 struct {
	c         *Client
	mountPath string
}

type KVSecret struct {
	Data     map[string]interface{}
	Metadata *KVVersionMetadata
	Raw      *Secret
}

type KVMetadata struct {
	CASRequired        bool              `mapstructure:"cas_required"`
	CreatedTime        time.Time         `mapstructure:"created_time"`
	CurrentVersion     int               `mapstructure:"current_version"`
	CustomMetadata     map[string]string `mapstructure:"custom_metadata"`
	DeleteVersionAfter time.Duration     `mapstructure:"delete_version_after"`
	MaxVersions        int               `mapstructure:"max_versions"`
	OldestVersion      int               `mapstructure:"oldest_version"`
	UpdatedTime        time.Time         `mapstructure:"updated_time"`
	// Keys are stringified ints, e.g. "3"
	Versions map[string]KVVersionMetadata `mapstructure:"versions"`
}

type KVVersionMetadata struct {
	Version      int       `mapstructure:"version"`
	CreatedTime  time.Time `mapstructure:"created_time"`
	DeletionTime time.Time `mapstructure:"deletion_time"`
	Destroyed    bool      `mapstructure:"destroyed"`
	// There is currently no version-specific custom metadata.
	// This field is just a copy of what's in the CustomMetadata field
	// for the full KVMetadata of the secret.
	CustomMetadata map[string]string `mapstructure:"custom_metadata"`
}

func (c *Client) KVv1(mountPath string) *kvv1 {
	return &kvv1{c: c, mountPath: mountPath}
}

func (c *Client) KVv2(mountPath string) *kvv2 {
	return &kvv2{c: c, mountPath: mountPath}
}

//// KV v1 methods ////

// Get returns a secret from the KV v1 secrets engine.
//
// The Metadata field in the returned *KVSecret will always be nil.
// The Raw field can be inspected for information about the lease,
// and passed to a LifetimeWatcher object for periodic renewal.
func (kv *kvv1) Get(ctx context.Context, secretPath string) (*KVSecret, error) {
	pathToRead := fmt.Sprintf("%s/%s", kv.mountPath, secretPath)

	secret, err := kv.c.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return nil, fmt.Errorf("error encountered while reading secret at %s: %w", pathToRead, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %s", pathToRead)
	}

	return &KVSecret{
		Data:     secret.Data,
		Metadata: nil,
		Raw:      secret,
	}, nil
}

func (kv *kvv1) Put(ctx context.Context, secretPath string, data map[string]interface{}) error {
	pathToWriteTo := fmt.Sprintf("%s/%s", kv.mountPath, secretPath)

	_, err := kv.c.Logical().WriteWithContext(ctx, pathToWriteTo, data)
	if err != nil {
		return fmt.Errorf("error writing secret to %s: %w", pathToWriteTo, err)
	}

	return nil
}

func (kv *kvv1) Delete(ctx context.Context, secretPath string) error {
	pathToDelete := fmt.Sprintf("%s/%s", kv.mountPath, secretPath)

	_, err := kv.c.Logical().DeleteWithContext(ctx, pathToDelete)
	if err != nil {
		return fmt.Errorf("error deleting secret at %s: %w", pathToDelete, err)
	}

	return nil
}

//// KV v2 methods ////

// Get returns the latest version of a secret from the KV v2 secrets engine.
//
// If the latest version has been deleted, an error will not be thrown, but
// the Data field on the returned secret will be nil, and the Metadata field
// will contain the deletion time.
func (kv *kvv2) Get(ctx context.Context, secretPath string) (*KVSecret, error) {
	pathToRead := fmt.Sprintf("%s/data/%s", kv.mountPath, secretPath)

	secret, err := kv.c.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return nil, fmt.Errorf("error encountered while reading secret at %s: %w", pathToRead, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %s", pathToRead)
	}

	return extractDataAndVersionMetadata(secret)
}

// GetVersion returns the data and metadata for a specific version of the
// given secret. If that version has been deleted, the Data field on the
// returned secret will be nil, and the Metadata field will contain the deletion time.
func (kv *kvv2) GetVersion(ctx context.Context, secretPath string, version int) (*KVSecret, error) {
	pathToRead := fmt.Sprintf("%s/data/%s", kv.mountPath, secretPath)

	queryParams := map[string][]string{"version": {strconv.Itoa(version)}}
	secret, err := kv.c.Logical().ReadWithDataWithContext(ctx, pathToRead, queryParams)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret with version %d found at %s", version, pathToRead)
	}

	return extractDataAndVersionMetadata(secret)
}

func (kv *kvv2) GetMetadata(ctx context.Context, secretPath string) (*KVMetadata, error) {
	pathToRead := fmt.Sprintf("%s/metadata/%s", kv.mountPath, secretPath)

	secret, err := kv.c.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret metadata found at %s", pathToRead)
	}

	return extractFullMetadata(secret)
}

func (kv *kvv2) Put(ctx context.Context, secretPath string, data map[string]interface{}) (*KVSecret, error) {
	pathToWriteTo := fmt.Sprintf("%s/data/%s", kv.mountPath, secretPath)

	secret, err := kv.c.Logical().WriteWithContext(ctx, pathToWriteTo, data)
	if err != nil {
		return nil, fmt.Errorf("error writing secret to %s: %w", pathToWriteTo, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret was written to %s", pathToWriteTo)
	}

	metadata, err := extractVersionMetadata(secret)
	if err != nil {
		return nil, fmt.Errorf("secret was written successfully, but unable to view version metadata from response: %w", err)
	}

	return &KVSecret{
		Data:     secret.Data,
		Metadata: metadata,
		Raw:      secret,
	}, nil
}

func (kv *kvv2) Delete(ctx context.Context, secretPath string) error {
	pathToDelete := fmt.Sprintf("%s/data/%s", kv.mountPath, secretPath)

	_, err := kv.c.Logical().DeleteWithContext(ctx, pathToDelete)
	if err != nil {
		return fmt.Errorf("error deleting secret at %s: %w", pathToDelete, err)
	}

	return nil
}

func (kv *kvv2) DeleteVersions(ctx context.Context, secretPath string, versions []int) error {
	// verb and path are different when trying to delete past versions
	pathToDelete := fmt.Sprintf("%s/delete/%s", kv.mountPath, secretPath)

	if len(versions) > 0 {
		var versionsToDelete []string
		for _, version := range versions {
			versionsToDelete = append(versionsToDelete, strconv.Itoa(version))
		}
		versionsMap := map[string]interface{}{
			"versions": versionsToDelete,
		}
		_, err := kv.c.Logical().Write(pathToDelete, versionsMap)
		if err != nil {
			return fmt.Errorf("error deleting secret at %s: %w", pathToDelete, err)
		}
	}

	return nil
}

func extractDataAndVersionMetadata(secret *Secret) (*KVSecret, error) {
	// A nil map is a valid value for data: secret.Data will be nil when this
	// version of the secret has been deleted, but the metadata is still
	// available.
	var data map[string]interface{}
	if secret.Data != nil {
		dataInterface, ok := secret.Data["data"]
		if !ok {
			return nil, fmt.Errorf("missing expected 'data' element")
		}

		if dataInterface != nil {
			data, ok = dataInterface.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected type for 'data' element: %T (%#v)", data, data)
			}
		}
	}

	metadata, err := extractVersionMetadata(secret)
	if err != nil {
		return nil, fmt.Errorf("unable to get version metadata: %w", err)
	}

	return &KVSecret{
		Data:     data,
		Metadata: metadata,
		Raw:      secret,
	}, nil
}

func extractVersionMetadata(secret *Secret) (*KVVersionMetadata, error) {
	var metadata *KVVersionMetadata

	if secret.Data != nil {
		// Writes return the metadata directly, Reads return it nested inside the "metadata" key
		var metadataMap map[string]interface{}
		metadataInterface, ok := secret.Data["metadata"]
		if ok {
			metadataMap, ok = metadataInterface.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected type for 'metadata' element: %T (%#v)", metadataInterface, metadataInterface)
			}
		} else {
			metadataMap = secret.Data
		}

		// deletion_time usually comes in as an empty string which can't be
		// processed as time.RFC3339, so we reset it to a convertible value
		if metadataMap["deletion_time"] == "" {
			metadataMap["deletion_time"] = time.Time{}
		}

		d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
			Result:     &metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("error setting up decoder for API response: %w", err)
		}

		err = d.Decode(metadataMap)
		if err != nil {
			return nil, fmt.Errorf("error decoding metadata from API response into VersionMetadata: %w", err)
		}
	}

	return metadata, nil
}

func extractFullMetadata(secret *Secret) (*KVMetadata, error) {
	//TODO: actually find a way to copy custom metadata from KVMetadata into the KVVersionMetadata struct on KV ReadMetadata
	var metadata *KVMetadata

	if secret.Data != nil {
		// deletion_time usually comes in as an empty string which can't be
		// processed as time.RFC3339, so we reset it to a convertible value
		if versions, ok := secret.Data["versions"]; ok {
			versionsMap := versions.(map[string]interface{})
			if len(versionsMap) > 0 {
				for version, metadata := range versionsMap {
					metadataMap := metadata.(map[string]interface{})
					if metadataMap["deletion_time"] == "" {
						metadataMap["deletion_time"] = time.Time{}
					}
					versionsMap[version] = metadataMap // save the updated copy of the metadata map
				}
			}
			secret.Data["versions"] = versionsMap // save the updated copy of the versions map
		}

		d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeHookFunc(time.RFC3339),
				mapstructure.StringToTimeDurationHookFunc(),
			),
			Result: &metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("error setting up decoder for API response: %w", err)
		}

		err = d.Decode(secret.Data)
		if err != nil {
			return nil, fmt.Errorf("error decoding metadata from API response into KVMetadata: %w", err)
		}
	}

	return metadata, nil
}
