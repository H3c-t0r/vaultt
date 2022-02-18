package api

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) CORSStatus() (*CORSResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.CORSStatusContext(ctx)
}

func (c *Sys) CORSStatusContext(ctx context.Context) (*CORSResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/config/cors")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result CORSResponse
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (c *Sys) ConfigureCORS(req *CORSRequest) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.ConfigureCORSContext(ctx, req)
}

func (c *Sys) ConfigureCORSContext(ctx context.Context, req *CORSRequest) error {
	r := c.c.NewRequest("PUT", "/v1/sys/config/cors")
	if err := r.SetJSONBody(req); err != nil {
		return err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) DisableCORS() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.DisableCORSContext(ctx)
}

func (c *Sys) DisableCORSContext(ctx context.Context) error {
	r := c.c.NewRequest("DELETE", "/v1/sys/config/cors")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

type CORSRequest struct {
	AllowedOrigins []string `json:"allowed_origins" mapstructure:"allowed_origins"`
	AllowedHeaders []string `json:"allowed_headers" mapstructure:"allowed_headers"`
	Enabled        bool     `json:"enabled" mapstructure:"enabled"`
}

type CORSResponse struct {
	AllowedOrigins []string `json:"allowed_origins" mapstructure:"allowed_origins"`
	AllowedHeaders []string `json:"allowed_headers" mapstructure:"allowed_headers"`
	Enabled        bool     `json:"enabled" mapstructure:"enabled"`
}
