// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cliconfig

import (
	"github.com/hashicorp/vault/api/tokenhelper"
)

// DefaultTokenHelper returns the token helper that is configured for Vault.
// This helper should only be used for non-server CLI commands.
func DefaultTokenHelper() (tokenhelper.TokenHelper, error) {
	config, err := LoadConfig("")
	if err != nil {
		return nil, err
	}

	path := config.TokenHelper
	if path == "" {
		return tokenhelper.NewInternalTokenHelper()
	}

	path, err = tokenhelper.ExternalTokenHelperPath(path)
	if err != nil {
		return nil, err
	}
	return &tokenhelper.ExternalTokenHelper{BinaryPath: path}, nil
}
