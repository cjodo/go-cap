package redcap

import (
	"context"
)

// ExportFile exports a file field from a record.
func (c *Client) ExportFile(ctx context.Context, recordID, field, event string) ([]byte, error) {
	params := map[string]string{
		"content": "file",
		"record":  recordID,
		"field":   field,
	}
	if event != "" {
		params["event"] = event
	}

	return c.Request(ctx, "", params)
}

// ImportFile imports a file into a record field.
func (c *Client) ImportFile(ctx context.Context, recordID, field, event string, data []byte, opts ...ImportOption) error {
	params := map[string]string{
		"content": "file",
		"record":  recordID,
		"field":   field,
	}
	if event != "" {
		params["event"] = event
	}

	for _, opt := range opts {
		opt(params)
	}

	_, err := c.Request(ctx, "", params)
	return err
}

// DeleteFile deletes a file from a record field.
func (c *Client) DeleteFile(ctx context.Context, recordID, field, event string) error {
	params := map[string]string{
		"content": "file",
		"record":  recordID,
		"field":   field,
	}
	if event != "" {
		params["event"] = event
	}

	_, err := c.Request(ctx, "", params)
	return err
}
