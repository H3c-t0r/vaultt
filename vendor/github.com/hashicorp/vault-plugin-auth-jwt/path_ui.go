// A throwaway file for super simple testing via a UI
package jwtauth

import (
	"context"
	"io/ioutil"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUI(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: `ui$`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathUI,
		},
	}
}

func (b *jwtAuthBackend) pathUI(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	data, err := ioutil.ReadFile("test_ui.html")
	if err != nil {
		panic(err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPRawBody:     string(data),
			logical.HTTPContentType: "text/html",
		},
	}

	return resp, nil
}
