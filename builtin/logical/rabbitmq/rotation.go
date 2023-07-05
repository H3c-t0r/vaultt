// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

const (
	// Default interval to check the queue for items needing rotation
	defaultQueueTickSeconds = 5

	// Config key to set an alternate interval
	queueTickIntervalKey = "rotation_queue_tick_interval"

	// WAL storage key used for static account rotations
	staticWALKey = "staticRotationKey"
)

// rotateExpiredStaticCreds will pop expired credentials (credentials whose priority
// represents a time before the present), rotate the associated credential, and push
// them back onto the queue with the new priority.
func (b *backend) rotateExpiredStaticCreds(ctx context.Context, req *logical.Request) error {
	var errs *multierror.Error

	for {
		keepGoing, err := b.rotateCredential(ctx, req.Storage)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		if !keepGoing {
			if errs.ErrorOrNil() != nil {
				return fmt.Errorf("error(s) occurred while rotating expired static credentials: %w", errs)
			} else {
				return nil
			}
		}
	}
}

// TODO check if not pop by key needed
func (b *backend) rotateCredential(ctx context.Context, storage logical.Storage) (rotated bool, err error) {
	item, err := b.credRotationQueue.Pop()
	if err != nil {
		if err == queue.ErrEmpty {
			return false, nil
		}
		return false, fmt.Errorf("failed to pop from queue for role %q: %w", item.Key, err)
	}
	if item.Priority > time.Now().Unix() {
		err = b.credRotationQueue.Push(item)
		if err != nil {
			return false, fmt.Errorf("failed to add item into the rotation queue for username %q: %w", item.Key, err)
		}
		return false, nil
	}

	cfg := item.Value.(staticRoleEntry)

	err = b.createStaticCredential(ctx, storage, cfg, item.Key)
	if err != nil {
		return false, err
	}

	// set new priority and re-queue
	item.Priority = time.Now().Add(cfg.RotationPeriod).Unix()
	err = b.credRotationQueue.Push(item)
	if err != nil {
		return false, fmt.Errorf("failed to add item into the rotation queue for username %q: %w", cfg.Username, err)
	}

	return true, nil
}

// TODO add option to lock storage while performing updates
func (b *backend) createStaticCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, entryName string) error {
	config, err := readConfig(ctx, storage)
	if err != nil {
		return fmt.Errorf("unable to read configuration: %w", err)
	}
	newPassword, err := b.generatePassword(ctx, config.PasswordPolicy)
	if err != nil {
		return err
	}
	cfg.Password = newPassword
	client, err := b.Client(ctx, storage)
	if err != nil {
		return err
	}
	if _, err = client.DeleteUser(cfg.Username); err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}
	_, err = client.PutUser(cfg.Username, rabbithole.UserSettings{
		Password: cfg.Password,
		Tags:     []string{cfg.RoleEntry.Tags},
	})
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}
	// update storage with new password and new rotation
	cfg.LastVaultRotation = time.Now()
	entry, err := logical.StorageEntryJSON(rabbitMQStaticRolePath+entryName, cfg)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return err
	}
	// TODO: refactor host and topic setting function from path_role_create and set permissions here
	return nil
}

func (b *backend) deleteStaticCredential(ctx context.Context, storage logical.Storage, cfg staticRoleEntry, shouldLockStorage bool) error {
	// TODO needed to pop from queue?
	// TODO remove from logical storage?
	if cfg.RevokeUserOnDelete {
		client, err := b.Client(ctx, storage)
		if err != nil {
			return err
		}
		if _, err = client.DeleteUser(cfg.Username); err != nil {
			return fmt.Errorf("could not delete user: %w", err)
		}
	}
	return nil
}
