// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
)

func (c *SQLConnectionProducer) getCloudSQLDriverName() (string, error) {
	var driverName string
	// using switch case for future extensibility
	switch c.Type {
	case dbTypePostgres:
		driverName = cloudSQLPostgres
	default:
		return "", fmt.Errorf("unsupported DB type for cloud IAM: %s", c.Type)
	}

	return driverName, nil
}

func (c *SQLConnectionProducer) registerDrivers(driverName string, credentials string) (func() error, error) {
	typ, err := c.getCloudSQLDriverName()
	if err != nil {
		return nil, err
	}

	opts, err := GetCloudSQLAuthOptions(credentials)
	if err != nil {
		return nil, err
	}

	// using switch case for future extensibility
	switch typ {
	case cloudSQLPostgres:
		return pgxv4.RegisterDriver(driverName, opts...)
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

// GetCloudSQLAuthOptions takes a credentials JSON and returns
// a set of GCP CloudSQL options - always WithIAMAUthN, and then the appropriate file/JSON option.
func GetCloudSQLAuthOptions(credentials string) ([]cloudsqlconn.Option, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithIAMAuthN()}

	if credentials != "" {
		opts = append(opts, cloudsqlconn.WithCredentialsJSON([]byte(credentials)))
	}

	return opts, nil
}
