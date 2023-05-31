// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"time"
)

// Database to manipulate users within an external system (typically a database).
type Database interface {
	// Initialize the database plugin. This is the equivalent of a constructor for the
	// database object itself.
	Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error)

	// NewUser creates a new user within the database. This user is temporary in that it
	// will exist until the TTL expires.
	NewUser(ctx context.Context, req NewUserRequest) (NewUserResponse, error)

	// UpdateUser updates an existing user within the database.
	UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error)

	// DeleteUser from the database. This should not error if the user didn't
	// exist prior to this call.
	DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error)

	// Type returns the Name for the particular database backend implementation.
	// This type name is usually set as a constant within the database backend
	// implementation, e.g. "mysql" for the MySQL database backend. This is used
	// for things like metrics and logging. No behavior is switched on this.
	Type() (string, error)

	// Close attempts to close the underlying database connection that was
	// established by the backend.
	Close() error
}

// ///////////////////////////////////////////////////////////////////////////
// Database Request & Response Objects
// These request and response objects are *not* protobuf types because gRPC does not
// support all types that we need in a nice way. For instance, gRPC does not support
// map[string]interface{}. It does have an `Any` type, but converting it to a map
// requires extensive use of reflection and knowing what types to support ahead of
// time. Instead these types are made as user-friendly as possible so the conversion
// between protobuf types and request/response objects is handled by Vault developers
// rather than needing to be handled by external plugin developers.
// ///////////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////
// Initialize()
// ///////////////////////////////////////////////////////

// InitializeRequest contains all information needed to initialize a database plugin.
type InitializeRequest struct {
	// Config to initialize the database with. This can include things like connection details,
	// a "root" username & password, etc. This will not include all configuration items specified
	// when configuring the database. Some values will be stripped out by the database engine
	// prior to being passed to the plugin.
	Config map[string]interface{}

	// VerifyConnection during initialization. If true, a connection should be made to the
	// database to verify the connection can be made. If false, no connection should be made
	// on initialization.
	VerifyConnection bool
}

// InitializeResponse returns any information Vault needs to know after initializing
// a database plugin.
type InitializeResponse struct {
	// Config that should be saved in Vault. This may differ from the config in the request,
	// but should contain everything required to Initialize the database.
	// REQUIRED in order to save the configuration into Vault after initialization
	Config map[string]interface{}
}

// SupportedCredentialTypesKey is used to get and set the supported
// CredentialType values in database plugins and Vault.
const SupportedCredentialTypesKey = "supported_credential_types"

// SetSupportedCredentialTypes sets the CredentialType values that are
// supported by the database plugin. It can be used by database plugins
// to communicate what CredentialType values it supports managing.
func (ir InitializeResponse) SetSupportedCredentialTypes(credTypes []CredentialType) {
	sct := make([]interface{}, 0, len(credTypes))
	for _, t := range credTypes {
		sct = append(sct, t.String())
	}

	ir.Config[SupportedCredentialTypesKey] = sct
}

// ///////////////////////////////////////////////////////
// NewUser()
// ///////////////////////////////////////////////////////

// NewUserRequest request a new user is created
type NewUserRequest struct {
	// UsernameConfig is metadata that can be used to generate a username
	// within the database plugin
	UsernameConfig UsernameMetadata

	// Statements is an ordered list of commands to run within the database when
	// creating a new user. This frequently includes permissions to give the
	// user or similar actions.
	Statements Statements

	// RollbackStatements is an ordered list of commands to run within the database
	// if the new user creation process fails.
	RollbackStatements Statements

	// CredentialType is the type of credential to use when creating a user.
	// Respective fields for the credential type will contain the credential
	// value that was generated by Vault.
	CredentialType CredentialType

	// Password credential to use when creating the user.
	// Value is set when the credential type is CredentialTypePassword.
	Password string

	// PublicKey credential to use when creating the user.
	// The value is a PKIX marshaled, PEM encoded public key.
	// The value is set when the credential type is CredentialTypeRSAPrivateKey.
	PublicKey []byte

	// Subject is the distinguished name for the client certificate credential.
	// Value is set when the credential type is CredentialTypeClientCertificate.
	Subject string

	// Expiration of the user. Not all database plugins will support this.
	Expiration time.Time
}

// UsernameMetadata is metadata the database plugin can use to generate a username
type UsernameMetadata struct {
	DisplayName string
	RoleName    string
}

// NewUserResponse returns any information Vault needs to know after creating a new user.
type NewUserResponse struct {
	// Username of the user created within the database.
	// REQUIRED so Vault knows the name of the user that was created
	Username string
}

// CredentialType is a type of database credential.
type CredentialType int

const (
	CredentialTypePassword CredentialType = iota
	CredentialTypeRSAPrivateKey
	CredentialTypeClientCertificate
)

func (k CredentialType) String() string {
	switch k {
	case CredentialTypePassword:
		return "password"
	case CredentialTypeRSAPrivateKey:
		return "rsa_private_key"
	case CredentialTypeClientCertificate:
		return "client_certificate"
	default:
		return "unknown"
	}
}

// ///////////////////////////////////////////////////////
// UpdateUser()
// ///////////////////////////////////////////////////////

type UpdateUserRequest struct {
	// Username to make changes to.
	Username string

	// CredentialType is the type of credential to use when updating a user.
	// Respective fields for the credential type will contain the credential
	// value that was generated by Vault.
	CredentialType CredentialType

	// Password indicates the new password to change to.
	// The value is set when the credential type is CredentialTypePassword.
	// If nil, no change is requested.
	Password *ChangePassword

	// PublicKey indicates the new public key to change to.
	// The value is set when the credential type is CredentialTypeRSAPrivateKey.
	// If nil, no change is requested.
	PublicKey *ChangePublicKey

	// Expiration indicates the new expiration date to change to.
	// If nil, no change is requested.
	Expiration *ChangeExpiration
}

// ChangePublicKey of a given user
type ChangePublicKey struct {
	// NewPublicKey is the new public key credential for the user.
	// The value is a PKIX marshaled, PEM encoded public key.
	NewPublicKey []byte

	// Statements is an ordered list of commands to run within the database
	// when changing the user's public key credential.
	Statements Statements
}

// ChangePassword of a given user
type ChangePassword struct {
	// NewPassword for the user
	NewPassword string

	// Statements is an ordered list of commands to run within the database
	// when changing the user's password.
	Statements Statements
}

// ChangeExpiration of a given user
type ChangeExpiration struct {
	// NewExpiration of the user
	NewExpiration time.Time

	// Statements is an ordered list of commands to run within the database
	// when changing the user's expiration.
	Statements Statements
}

type UpdateUserResponse struct{}

// ///////////////////////////////////////////////////////
// DeleteUser()
// ///////////////////////////////////////////////////////

type DeleteUserRequest struct {
	// Username to delete from the database
	Username string

	// Statements is an ordered list of commands to run within the database
	// when deleting a user.
	Statements Statements
}

type DeleteUserResponse struct{}

// ///////////////////////////////////////////////////////
// Used across multiple functions
// ///////////////////////////////////////////////////////

// Statements wraps a collection of statements to run in a database when an
// operation is performed (create, update, etc.). This is a struct rather than
// a string slice so we can easily add more information to this in the future.
type Statements struct {
	// Commands is an ordered list of commands to execute in the database.
	// These commands may include templated fields such as {{username}} and {{password}}
	Commands []string
}
