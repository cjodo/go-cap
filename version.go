package redcap

import (
	"context"
	"strings"
)

func (c *Client) ExportVersion(ctx context.Context) (string, error) {
	resp, err := c.Request(ctx, "version", nil)
	if err != nil {
		return "", err
	}
	// REDCap version endpoint returns plain text, not JSON
	return strings.TrimSpace(strings.Trim(string(resp), `"`)), nil
}
