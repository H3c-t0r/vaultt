// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-uuid"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/mitchellh/mapstructure"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
)

const (
	AuthTypeGCPIAM = "gcp_iam"
)

const (
	dbTypePostgres   = "pgx"
	dbTypeMSSQL      = "mssql"
	cloudSQLPostgres = "cloudsql-postgres"
	cloudSQLMSSQL    = "cloudsql-sqlserver"
)

var _ ConnectionProducer = &SQLConnectionProducer{}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]cloudSQLCleanup)
)

type cloudSQLCleanup func() error

// SQLConnectionProducer implements ConnectionProducer and provides a generic producer for most sql databases
type SQLConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url" mapstructure:"connection_url" structs:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections" mapstructure:"max_open_connections" structs:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections" mapstructure:"max_idle_connections" structs:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime" structs:"max_connection_lifetime"`
	Username                 string      `json:"username" mapstructure:"username" structs:"username"`
	Password                 string      `json:"password" mapstructure:"password" structs:"password"`
	AuthType                 string      `json:"auth_type" mapstructure:"auth_type" structs:"auth_type"`
	Credentials              string      `json:"credentials" mapstructure:"credentials" structs:"credentials"`
	DisableEscaping          bool        `json:"disable_escaping" mapstructure:"disable_escaping" structs:"disable_escaping"`

	// cloud options here - cloudDriverName is globally unique, but only needs to be retained for the lifetime
	// of driver registration, not from run to run.
	cloudDriverName string

	Type                  string
	RawConfig             map[string]interface{}
	maxConnectionLifetime time.Duration
	Initialized           bool
	db                    *sql.DB
	sync.Mutex
}

func (c *SQLConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

func (c *SQLConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.RawConfig = conf

	err := mapstructure.WeakDecode(conf, &c)
	if err != nil {
		return nil, err
	}

	if len(c.ConnectionURL) == 0 {
		return nil, fmt.Errorf("connection_url cannot be empty")
	}

	// Do not allow the username or password template pattern to be used as
	// part of the user-supplied username or password
	if strings.Contains(c.Username, "{{username}}") ||
		strings.Contains(c.Username, "{{password}}") ||
		strings.Contains(c.Password, "{{username}}") ||
		strings.Contains(c.Password, "{{password}}") {

		return nil, fmt.Errorf("username and/or password cannot contain the template variables")
	}

	// Don't escape special characters for MySQL password
	// Also don't escape special characters for the username and password if
	// the disable_escaping parameter is set to true
	username := c.Username
	password := c.Password
	if !c.DisableEscaping {
		username = url.PathEscape(c.Username)
	}
	if (c.Type != "mysql") && !c.DisableEscaping {
		password = url.PathEscape(c.Password)
	}

	// QueryHelper doesn't do any SQL escaping, but if it starts to do so
	// then maybe we won't be able to use it to do URL substitution any more.
	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": username,
		"password": password,
	})

	if c.MaxOpenConnections == 0 {
		c.MaxOpenConnections = 4
	}

	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = c.MaxOpenConnections
	}
	if c.MaxIdleConnections > c.MaxOpenConnections {
		c.MaxIdleConnections = c.MaxOpenConnections
	}
	if c.MaxConnectionLifetimeRaw == nil {
		c.MaxConnectionLifetimeRaw = "0s"
	}

	c.maxConnectionLifetime, err = parseutil.ParseDurationSecond(c.MaxConnectionLifetimeRaw)
	if err != nil {
		return nil, errwrap.Wrapf("invalid max_connection_lifetime: {{err}}", err)
	}

	if c.AuthType == AuthTypeGCPIAM {
		c.cloudDriverName, err = uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("unable to generate UUID for IAM configuration: %w", err)
		}

		// There are a few important points to keep in mind with this line of code, for more information
		// see the connection_producer.go
		_, err := c.registerDrivers(c.cloudDriverName, c.Credentials)
		if err != nil {
			return nil, err
		}
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}

		if err := c.db.PingContext(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}
	}

	return c.RawConfig, nil
}

func (c *SQLConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, ErrNotInitialized
	}

	// If we already have a DB, test it and return
	if c.db != nil {
		if err := c.db.PingContext(ctx); err == nil {
			return c.db, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		c.db.Close()
	}

	// default non-IAM behavior
	driverName := c.Type

	if c.AuthType == AuthTypeGCPIAM {
		driverName = c.cloudDriverName
	} else if c.Type == "mssql" {
		// For mssql backend, switch to sqlserver instead
		driverName = "sqlserver"
	}

	// Otherwise, attempt to make connection
	conn := c.ConnectionURL

	// PostgreSQL specific settings
	if strings.HasPrefix(conn, "postgres://") || strings.HasPrefix(conn, "postgresql://") {
		// Ensure timezone is set to UTC for all the connections
		if strings.Contains(conn, "?") {
			conn += "&timezone=UTC"
		} else {
			conn += "?timezone=UTC"
		}

		// Ensure a reasonable application_name is set
		if !strings.Contains(conn, "application_name") {
			conn += "&application_name=vault"
		}
	}

	var err error
	c.db, err = sql.Open(driverName, conn)
	if err != nil {
		return nil, err
	}

	// Set some connection pool settings. We don't need much of this,
	// since the request rate shouldn't be high.
	c.db.SetMaxOpenConns(c.MaxOpenConnections)
	c.db.SetMaxIdleConns(c.MaxIdleConnections)
	c.db.SetConnMaxLifetime(c.maxConnectionLifetime)

	return c.db, nil
}

func (c *SQLConnectionProducer) SecretValues() map[string]interface{} {
	return map[string]interface{}{
		c.Password: "[password]",
	}
}

// Close attempts to close the connection
func (c *SQLConnectionProducer) Close() error {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		// if auth_type is IAM, ensure cleanup
		// of cloudSQL resources
		if c.AuthType == AuthTypeGCPIAM {
			driversMu.Lock()
			defer driversMu.Unlock()
		} else {
			c.db.Close()
		}
	}

	c.db = nil

	return nil
}

// SetCredentials uses provided information to set/create a user in the
// database. Unlike CreateUser, this method requires a username be provided and
// uses the name given, instead of generating a name. This is used for creating
// and setting the password of static accounts, as well as rolling back
// passwords in the database in the event an updated database fails to save in
// Vault's storage.
func (c *SQLConnectionProducer) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticUser dbplugin.StaticUserConfig) (username, password string, err error) {
	return "", "", dbutil.Unimplemented()
}
