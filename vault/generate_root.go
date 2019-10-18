package vault

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/shamir"
	shamirseal "github.com/hashicorp/vault/vault/seal/shamir"
)

const coreDROperationTokenPath = "core/dr-operation-token"

var (
	// GenerateStandardRootTokenStrategy is the strategy used to generate a
	// typical root token
	GenerateStandardRootTokenStrategy GenerateRootStrategy = generateStandardRootToken{}

	// GenerateDROperationTokenStrategy is the strategy used to generate a
	// DR operational token
	GenerateDROperationTokenStrategy GenerateRootStrategy = generateStandardRootToken{}
)

// GenerateRootStrategy allows us to swap out the strategy we want to use to
// create a token upon completion of the generate root process.
type GenerateRootStrategy interface {
	generate(context.Context, *Core) (string, func(), error)
}

// generateStandardRootToken implements the GenerateRootStrategy and is in
// charge of creating standard root tokens.
type generateStandardRootToken struct{}

func (g generateStandardRootToken) generate(ctx context.Context, c *Core) (string, func(), error) {
	te, err := c.tokenStore.rootToken(ctx)
	if err != nil {
		c.logger.Error("root token generation failed", "error", err)
		return "", nil, err
	}
	if te == nil {
		c.logger.Error("got nil token entry back from root generation")
		return "", nil, fmt.Errorf("got nil token entry back from root generation")
	}

	cleanupFunc := func() {
		c.tokenStore.revokeOrphan(ctx, te.ID)
	}

	return te.ID, cleanupFunc, nil
}

// GenerateRootConfig holds the configuration for a root generation
// command.
type GenerateRootConfig struct {
	Nonce          string
	PGPKey         string
	PGPFingerprint string
	OTP            string
	Strategy       GenerateRootStrategy
}

// GenerateRootResult holds the result of a root generation update
// command
type GenerateRootResult struct {
	Progress       int
	Required       int
	EncodedToken   string
	PGPFingerprint string
}

// GenerateRootProgress is used to return the root generation progress (num shares)
func (c *Core) GenerateRootProgress() (int, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return 0, consts.ErrSealed
	}
	if c.standby && !c.recoveryMode {
		return 0, consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	return len(c.generateRootProgress), nil
}

// GenerateRootConfiguration is used to read the root generation configuration
// It stubbornly refuses to return the OTP if one is there.
func (c *Core) GenerateRootConfiguration() (*GenerateRootConfig, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return nil, consts.ErrSealed
	}
	if c.standby && !c.recoveryMode {
		return nil, consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Copy the config if any
	var conf *GenerateRootConfig
	if c.generateRootConfig != nil {
		conf = new(GenerateRootConfig)
		*conf = *c.generateRootConfig
		conf.OTP = ""
		conf.Strategy = nil
	}
	return conf, nil
}

// GenerateRootInit is used to initialize the root generation settings
func (c *Core) GenerateRootInit(otp, pgpKey string, strategy GenerateRootStrategy) error {
	var fingerprint string
	switch {
	case len(otp) > 0:
		if len(otp) != TokenLength+2 {
			return fmt.Errorf("OTP string is wrong length")
		}

	case len(pgpKey) > 0:
		fingerprints, err := pgpkeys.GetFingerprints([]string{pgpKey}, nil)
		if err != nil {
			return errwrap.Wrapf("error parsing PGP key: {{err}}", err)
		}
		if len(fingerprints) != 1 || fingerprints[0] == "" {
			return fmt.Errorf("could not acquire PGP key entity")
		}
		fingerprint = fingerprints[0]

	default:
		return fmt.Errorf("otp or pgp_key parameter must be provided")
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return consts.ErrSealed
	}
	barrierSealed, err := c.barrier.Sealed()
	if err != nil {
		return errors.New("unable to check barrier seal status")
	}
	if !barrierSealed && c.recoveryMode {
		return errors.New("attempt to generate recovery operation token when already unsealed")
	}
	if c.standby && !c.recoveryMode {
		return consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Prevent multiple concurrent root generations
	if c.generateRootConfig != nil {
		return fmt.Errorf("root generation already in progress")
	}

	// Copy the configuration
	generationNonce, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	c.generateRootConfig = &GenerateRootConfig{
		Nonce:          generationNonce,
		OTP:            otp,
		PGPKey:         pgpKey,
		PGPFingerprint: fingerprint,
		Strategy:       strategy,
	}

	if c.logger.IsInfo() {
		switch strategy.(type) {
		case generateStandardRootToken:
			c.logger.Info("root generation initialized", "nonce", c.generateRootConfig.Nonce)
		case *generateRecoveryToken:
			c.logger.Info("recovery operation token generation initialized", "nonce", c.generateRootConfig.Nonce)
		default:
			c.logger.Info("dr operation token generation initialized", "nonce", c.generateRootConfig.Nonce)
		}
	}

	return nil
}

// GenerateRootUpdate is used to provide a new key part
func (c *Core) GenerateRootUpdate(ctx context.Context, key []byte, nonce string, strategy GenerateRootStrategy) (*GenerateRootResult, error) {
	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return nil, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	// Get the seal configuration
	var config *SealConfig
	var err error
	if c.seal.RecoveryKeySupported() {
		config, err = c.seal.RecoveryConfig(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = c.seal.BarrierConfig(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Ensure the barrier is initialized
	if config == nil {
		return nil, ErrNotInit
	}

	// Ensure we are already unsealed
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return nil, consts.ErrSealed
	}

	barrierSealed, err := c.barrier.Sealed()
	if err != nil {
		return nil, errors.New("unable to check barrier seal status")
	}
	if !barrierSealed && c.recoveryMode {
		return nil, errors.New("attempt to generate recovery operation token when already unsealed")
	}

	if c.standby && !c.recoveryMode {
		return nil, consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Ensure a generateRoot is in progress
	if c.generateRootConfig == nil {
		return nil, fmt.Errorf("no root generation in progress")
	}

	if nonce != c.generateRootConfig.Nonce {
		return nil, fmt.Errorf("incorrect nonce supplied; nonce for this root generation operation is %q", c.generateRootConfig.Nonce)
	}

	if strategy != c.generateRootConfig.Strategy {
		return nil, fmt.Errorf("incorrect strategy supplied; a generate root operation of another type is already in progress")
	}

	// Check if we already have this piece
	for _, existing := range c.generateRootProgress {
		if bytes.Equal(existing, key) {
			return nil, fmt.Errorf("given key has already been provided during this generation operation")
		}
	}

	// Store this key
	c.generateRootProgress = append(c.generateRootProgress, key)
	progress := len(c.generateRootProgress)

	// Check if we don't have enough keys to unlock
	if len(c.generateRootProgress) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot generate root, not enough keys", "keys", progress, "threshold", config.SecretThreshold)
		}
		return &GenerateRootResult{
			Progress:       progress,
			Required:       config.SecretThreshold,
			PGPFingerprint: c.generateRootConfig.PGPFingerprint,
		}, nil
	}

	// Combine the key parts
	var combinedKey []byte
	if config.SecretThreshold == 1 {
		combinedKey = c.generateRootProgress[0]
		c.generateRootProgress = nil
	} else {
		combinedKey, err = shamir.Combine(c.generateRootProgress)
		c.generateRootProgress = nil
		if err != nil {
			return nil, errwrap.Wrapf("failed to compute master key: {{err}}", err)
		}
	}

	switch {
	case c.seal.RecoveryKeySupported():
		// Ensure that the combined recovery key is valid
		if err := c.seal.VerifyRecoveryKey(ctx, combinedKey); err != nil {
			c.logger.Error("root generation aborted, recovery key verification failed", "error", err)
			return nil, err
		}
		// If we are in recovery mode, then retrieve
		// the stored keys and unseal the barrier
		if c.recoveryMode {
			storedKeys, err := c.seal.GetStoredKeys(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("unable to retrieve stored keys in recovery mode: {{err}}", err)
			}

			// Use the retrieved master key to unseal the barrier
			if err := c.barrier.Unseal(ctx, storedKeys[0]); err != nil {
				c.logger.Error("root generation aborted, recovery operation token verification failed", "error", err)
				return nil, err
			}
		}
	default:
		masterKey := combinedKey
		if c.seal.StoredKeysSupported() == StoredKeysSupportedShamirMaster {
			testseal := NewDefaultSeal(shamirseal.NewSeal(c.logger.Named("testseal")))
			testseal.SetCore(c)
			cfg, err := c.seal.BarrierConfig(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("failed to setup test barrier config: {{err}}", err)
			}
			testseal.SetCachedBarrierConfig(cfg)
			err = testseal.GetAccess().(*shamirseal.ShamirSeal).SetKey(combinedKey)
			if err != nil {
				return nil, errwrap.Wrapf("failed to setup unseal key: {{err}}", err)
			}
			stored, err := testseal.GetStoredKeys(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("failed to read master key: {{err}}", err)
			}
			masterKey = stored[0]
		}
		switch {
		case c.recoveryMode:
			// If we are in recovery mode, being able to unseal
			// the barrier is how we establish authentication
			if err := c.barrier.Unseal(ctx, masterKey); err != nil {
				c.logger.Error("root generation aborted, recovery operation token verification failed", "error", err)
				return nil, err
			}
		default:
			if err := c.barrier.VerifyMaster(masterKey); err != nil {
				c.logger.Error("root generation aborted, master key verification failed", "error", err)
				return nil, err
			}
		}
	}

	// Authentication in recovery mode is successful
	if c.recoveryMode {
		// Run any post unseal functions that are set
		for _, v := range c.postRecoveryUnsealFuncs {
			if err := v(); err != nil {
				return nil, errwrap.Wrapf("failed to run post unseal func: {{err}}", err)
			}
		}
	}

	// Run the generate strategy
	token, cleanupFunc, err := strategy.generate(ctx, c)
	if err != nil {
		return nil, err
	}

	var tokenBytes []byte

	// Get the encoded value first so that if there is an error we don't create
	// the root token.
	switch {
	case len(c.generateRootConfig.OTP) > 0:
		// This function performs decoding checks so rather than decode the OTP,
		// just encode the value we're passing in.
		tokenBytes, err = xor.XORBytes([]byte(c.generateRootConfig.OTP), []byte(token))
		if err != nil {
			cleanupFunc()
			c.logger.Error("xor of root token failed", "error", err)
			return nil, err
		}
		token = base64.RawStdEncoding.EncodeToString(tokenBytes)

	case len(c.generateRootConfig.PGPKey) > 0:
		_, tokenBytesArr, err := pgpkeys.EncryptShares([][]byte{[]byte(token)}, []string{c.generateRootConfig.PGPKey})
		if err != nil {
			cleanupFunc()
			c.logger.Error("error encrypting new root token", "error", err)
			return nil, err
		}
		token = base64.StdEncoding.EncodeToString(tokenBytesArr[0])

	default:
		cleanupFunc()
		return nil, fmt.Errorf("unreachable condition")
	}

	results := &GenerateRootResult{
		Progress:       progress,
		Required:       config.SecretThreshold,
		EncodedToken:   token,
		PGPFingerprint: c.generateRootConfig.PGPFingerprint,
	}

	switch strategy.(type) {
	case generateStandardRootToken:
		c.logger.Info("root generation finished", "nonce", c.generateRootConfig.Nonce)
	case *generateRecoveryToken:
		c.logger.Info("recovery operation token generation finished", "nonce", c.generateRootConfig.Nonce)
	default:
		c.logger.Info("dr operation token generation finished", "nonce", c.generateRootConfig.Nonce)
	}

	c.generateRootProgress = nil
	c.generateRootConfig = nil
	return results, nil
}

// GenerateRootCancel is used to cancel an in-progress root generation
func (c *Core) GenerateRootCancel() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.Sealed() && !c.recoveryMode {
		return consts.ErrSealed
	}
	if c.standby && !c.recoveryMode {
		return consts.ErrStandby
	}

	c.generateRootLock.Lock()
	defer c.generateRootLock.Unlock()

	// Clear any progress or config
	c.generateRootConfig = nil
	c.generateRootProgress = nil
	return nil
}
