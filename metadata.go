package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExportMetadata returns the data dictionary (metadata) for the project.
func (c *Client) ExportMetadata(ctx context.Context) ([]Field, error) {
	body, err := c.Request(ctx, "metadata", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var fields []Field
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, fmt.Errorf("unmarshaling metadata: %w", err)
	}

	return fields, nil
}
