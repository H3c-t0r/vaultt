package api

import "context"

func (c *Sys) StepDown() error {
	return c.StepDownWithContext(context.Background())
}

func (c *Sys) StepDownWithContext(ctx context.Context) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest("PUT", "/v1/sys/step-down")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err
}
