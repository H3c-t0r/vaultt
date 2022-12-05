package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	vaultVersionPath string = "core/versions/"
)

// storeVersionEntry will store the version, timestamp, and build date to storage
// only if no entry for that version already exists in storage. Version
// timestamps were initially stored in local time. UTC should be used. Existing
// entries can be overwritten via the force flag. A bool will be returned
// denoting whether the entry was updated
func (c *Core) storeVersionEntry(ctx context.Context, vaultVersion *VaultVersion, force bool) (bool, error) {
	key := vaultVersionPath + vaultVersion.Version

	if vaultVersion.TimestampInstalled.Location() != time.UTC {
		vaultVersion.TimestampInstalled = vaultVersion.TimestampInstalled.UTC()
	}

	marshalledVaultVersion, err := json.Marshal(vaultVersion)
	if err != nil {
		return false, err
	}

	newEntry := &logical.StorageEntry{
		Key:   key,
		Value: marshalledVaultVersion,
	}

	if force {
		// avoid storage lookup and write immediately
		err = c.barrier.Put(ctx, newEntry)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	existingEntry, err := c.barrier.Get(ctx, key)
	if err != nil {
		return false, err
	}

	if existingEntry != nil {
		return false, nil
	}

	err = c.barrier.Put(ctx, newEntry)

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindOldestVersionTimestamp searches for the vault version with the oldest
// upgrade timestamp from storage. The earliest version this can be is 1.9.0.
func (c *Core) FindOldestVersionTimestamp() (string, time.Time, error) {
	if c.versionHistory == nil {
		return "", time.Time{}, fmt.Errorf("version history is not initialized")
	}

	oldestUpgradeTime := time.Now().UTC()
	var oldestVersion string

	for versionStr, versionEntry := range c.versionHistory {
		if versionEntry.TimestampInstalled.Before(oldestUpgradeTime) {
			oldestVersion = versionStr
			oldestUpgradeTime = versionEntry.TimestampInstalled
		}
	}
	return oldestVersion, oldestUpgradeTime, nil
}

func (c *Core) FindNewestVersionTimestamp() (string, time.Time, error) {
	if c.versionHistory == nil {
		return "", time.Time{}, fmt.Errorf("version history is not initialized")
	}

	var newestUpgradeTime time.Time
	var newestVersion string

	for versionStr, versionEntry := range c.versionHistory {
		if versionEntry.TimestampInstalled.After(newestUpgradeTime) {
			newestVersion = versionStr
			newestUpgradeTime = versionEntry.TimestampInstalled
		}
	}

	return newestVersion, newestUpgradeTime, nil
}

// loadVersionHistory loads all the vault versions entries from storage.
// Version timestamps were originally stored in local time. A timestamp
// that is not in UTC will be rewritten to storage as UTC.
func (c *Core) loadVersionHistory(ctx context.Context) error {
	vaultVersions, err := c.barrier.List(ctx, vaultVersionPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve vault versions from storage: %w", err)
	}

	for _, versionPath := range vaultVersions {
		version, err := c.barrier.Get(ctx, vaultVersionPath+versionPath)
		if err != nil {
			return fmt.Errorf("unable to read vault version at path %s: err %w", versionPath, err)
		}
		if version == nil {
			return fmt.Errorf("nil version stored at path %s", versionPath)
		}
		var vaultVersion VaultVersion
		err = json.Unmarshal(version.Value, &vaultVersion)
		if err != nil {
			return fmt.Errorf("unable to unmarshal vault version for path %s: err %w", versionPath, err)
		}
		if vaultVersion.Version == "" || vaultVersion.TimestampInstalled.IsZero() {
			return fmt.Errorf("found empty serialized vault version at path %s", versionPath)
		}

		// self-heal entries that were not stored in UTC
		if vaultVersion.TimestampInstalled.Location() != time.UTC {
			vaultVersion.TimestampInstalled = vaultVersion.TimestampInstalled.UTC()

			isUpdated, err := c.storeVersionEntry(ctx, &vaultVersion, true)
			if err != nil {
				c.logger.Warn("failed to rewrite vault version timestamp as UTC", "error", err)
			}

			if isUpdated {
				c.logger.Info("self-healed pre-existing vault version in UTC",
					"vault version", vaultVersion.Version, "UTC time", vaultVersion.TimestampInstalled)
			}
		}
	}
	return nil
}

// isMajorOrMinorUpgrade compares two versions of Vault to see if currentVersion is is a
// major/minor upgrade from the prevVersion. This is useful in determining
// shutdown behavior for deprecated builtins.
func isMajorOrMinorUpgrade(currentVersion, prevVersion string) bool {
	// Get versions into comparable form
	curr, err := semver.NewSemver(currentVersion)
	if err != nil {
		return false
	}
	prev, err := semver.NewSemver(prevVersion)
	if err != nil {
		// If we can't find a previous version, this is effectively an upgrade
		return true
	}

	// Check for major version upgrade
	if curr.Segments()[0] > prev.Segments()[0] {
		return true
	}

	// Check for minor version upgrade
	if curr.Segments()[1] > prev.Segments()[1] {
		return true
	}

	return false
}

func IsJWT(token string) bool {
	return len(token) > 3 && strings.Count(token, ".") == 2 &&
		(token[3] != '.' && token[1] != '.')
}

func IsSSCToken(token string) bool {
	return len(token) > MaxNsIdLength+TokenLength+TokenPrefixLength &&
		strings.HasPrefix(token, consts.ServiceTokenPrefix)
}

func IsServiceToken(token string) bool {
	return strings.HasPrefix(token, consts.ServiceTokenPrefix) ||
		strings.HasPrefix(token, consts.LegacyServiceTokenPrefix)
}

func IsBatchToken(token string) bool {
	return strings.HasPrefix(token, consts.LegacyBatchTokenPrefix) ||
		strings.HasPrefix(token, consts.BatchTokenPrefix)
}
