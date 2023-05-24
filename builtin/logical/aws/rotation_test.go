package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

// TestRotation verifies that the rotation code and priority queue correctly selects and rotates credentials
// for static secrets.
func TestRotation(t *testing.T) {
	bgCTX := context.Background()

	type credToInsert struct {
		config staticRoleEntry // role configuration from a normal createRole request
		age    time.Duration   // how old the cred should be - if this is longer than the config.RotationPeriod,
		// the cred is 'pre-expired'

		changed bool // whether we expect the cred to change - this is technically redundant to a comparison between
		// rotationPeriod and age.
	}

	cases := []struct {
		name  string
		creds []credToInsert
	}{
		{
			name: "refresh one",
			creds: []credToInsert{
				{
					config: staticRoleEntry{
						Name:           "test",
						Username:       "jane-doe",
						RotationPeriod: 2 * time.Second,
					},
					age:     5 * time.Second,
					changed: true,
				},
			},
		},
		{
			name: "refresh none",
			creds: []credToInsert{
				{
					config: staticRoleEntry{
						Name:           "test",
						Username:       "jane-doe",
						RotationPeriod: 1 * time.Minute,
					},
					age:     5 * time.Second,
					changed: false,
				},
			},
		},
		{
			name: "refresh one of two",
			creds: []credToInsert{
				{
					config: staticRoleEntry{
						Name:           "test",
						Username:       "jane-doe",
						RotationPeriod: 1 * time.Minute,
					},
					age:     5 * time.Second,
					changed: false,
				},
				{
					config: staticRoleEntry{
						Name:           "toast",
						Username:       "john-doe",
						RotationPeriod: 1 * time.Second,
					},
					age:     5 * time.Second,
					changed: true,
				},
			},
		},
		{
			name:  "no creds to rotate",
			creds: []credToInsert{},
		},
	}

	ak := "long-access-key-id"
	oldSecret := "abcdefghijklmnopqrstuvwxyz"
	newSecret := "zyxwvutsrqponmlkjihgfedcba"

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := logical.TestBackendConfig()
			config.StorageView = &logical.InmemStorage{}

			b := Backend(config)
			// this means the creds will be the same for every user, but that's okay
			// since what we care about is whether they changed on a single-user basis.
			miam, err := awsutil.NewMockIAM(
				// blank list for existing user
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{},
					},
				}),
				// initial key to store
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String(ak),
						SecretAccessKey: aws.String(oldSecret),
					},
				}),
			)(nil)
			if err != nil {
				t.Fatalf("couldn't initialze mock IAM handler: %s", err)
			}

			// insert all our creds
			for i, cred := range c.creds {
				b.iamClient = miam

				err = b.createCredential(bgCTX, config.StorageView, cred.config, true)
				if err != nil {
					t.Fatalf("couldn't insert credential %d: %s", i, err)
				}

				item := &queue.Item{
					Key:      cred.config.Name,
					Value:    cred.config,
					Priority: time.Now().Add(-1 * cred.age).Add(cred.config.RotationPeriod).Unix(),
				}
				err = b.credRotationQueue.Push(item)
				if err != nil {
					t.Fatalf("couldn't push item onto queue: %s", err)
				}
			}

			// update aws responses, same argument for why it's okay every cred will be the same
			miam, err = awsutil.NewMockIAM(
				// old key
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{
							AccessKeyId: aws.String(ak),
						},
					},
				}),
				// new key
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String(ak),
						SecretAccessKey: aws.String(newSecret),
					},
				}),
			)(nil)
			if err != nil {
				t.Fatalf("couldn't initialze mock IAM handler: %s", err)
			}
			b.iamClient = miam

			req := &logical.Request{
				Storage: config.StorageView,
			}
			err = b.rotateExpiredStaticCreds(bgCTX, req)
			if err != nil {
				t.Fatalf("got an error rotating credentials: %s", err)
			}

			// check our credentials
			for i, cred := range c.creds {
				entry, err := config.StorageView.Get(bgCTX, formatCredsStoragePath(cred.config.Name))
				if err != nil {
					t.Fatalf("got an error retrieving credentials %d", i)
				}
				var out awsCredentials
				err = entry.DecodeJSON(&out)
				if err != nil {
					t.Fatalf("could not unmarshal storage view entry for cred %d to an aws credential: %s", i, err)
				}

				if cred.changed && out.SecretAccessKey != newSecret {
					t.Fatalf("expected the key for cred %d to have changed, but it hasn't", i)
				} else if !cred.changed && out.SecretAccessKey != oldSecret {
					t.Fatalf("expected the key for cred %d to have stayed the same, but it changed", i)
				}
			}
		})
	}
}

type fakeIAM struct {
	iamiface.IAMAPI
	delReqs []*iam.DeleteAccessKeyInput
}

func (f *fakeIAM) DeleteAccessKey(r *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	f.delReqs = append(f.delReqs, r)
	return f.IAMAPI.DeleteAccessKey(r)
}

// TestCreateCredential verifies that credential creation firstly only deletes credentials if it needs to (i.e., two
// or more credentials on IAM), and secondly correctly deletes the oldest one.
func TestCreateCredential(t *testing.T) {
	cases := []struct {
		name       string
		username   string
		deletedKey string
		opts       []awsutil.MockIAMOption
	}{
		{
			name:     "zero keys",
			username: "jane-doe",
			opts: []awsutil.MockIAMOption{
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{},
				}),
				// delete should _not_ be called
				awsutil.WithDeleteAccessKeyError(errors.New("should not have been called")),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("key"),
						SecretAccessKey: aws.String("itsasecret"),
					},
				}),
			},
		},
		{
			name:     "one key",
			username: "jane-doe",
			opts: []awsutil.MockIAMOption{
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{AccessKeyId: aws.String("foo"), CreateDate: aws.Time(time.Now())},
					},
				}),
				// delete should _not_ be called
				awsutil.WithDeleteAccessKeyError(errors.New("should not have been called")),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("key"),
						SecretAccessKey: aws.String("itsasecret"),
					},
				}),
			},
		},
		{
			name:       "two keys",
			username:   "jane-doe",
			deletedKey: "foo",
			opts: []awsutil.MockIAMOption{
				awsutil.WithListAccessKeysOutput(&iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{AccessKeyId: aws.String("foo"), CreateDate: aws.Time(time.Time{})},
						{AccessKeyId: aws.String("bar"), CreateDate: aws.Time(time.Now())},
					},
				}),
				awsutil.WithCreateAccessKeyOutput(&iam.CreateAccessKeyOutput{
					AccessKey: &iam.AccessKey{
						AccessKeyId:     aws.String("key"),
						SecretAccessKey: aws.String("itsasecret"),
					},
				}),
			},
		},
	}

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			miam, err := awsutil.NewMockIAM(
				c.opts...,
			)(nil)
			if err != nil {
				t.Fatal(err)
			}
			fiam := &fakeIAM{
				IAMAPI: miam,
			}

			b := Backend(config)
			b.iamClient = fiam

			err = b.createCredential(context.Background(), config.StorageView, staticRoleEntry{Username: c.username}, true)
			if err != nil {
				t.Fatalf("got an error we didn't expect: %q", err)
			}

			if c.deletedKey != "" {
				if len(fiam.delReqs) != 1 {
					t.Fatalf("called the wrong number of deletes (called %d deletes)", len(fiam.delReqs))
				}
				actualKey := *fiam.delReqs[0].AccessKeyId
				if c.deletedKey != actualKey {
					t.Fatalf("we deleted the wrong key: %q instead of %q", actualKey, c.deletedKey)
				}
			}
		})
	}
}
