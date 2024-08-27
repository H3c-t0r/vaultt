// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package connutil

import (
	"context"
	"database/sql"
	"errors"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func (c *SQLConnectionProducer) StaticConnection(_ context.Context, _, _ string) (*sql.DB, error) {
	return nil, errors.New("self-managed static roles not implemented in CE")
}
