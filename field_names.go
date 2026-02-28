package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExportFieldNames returns the list of export field names.
func (c *Client) ExportFieldNames(ctx context.Context) ([]string, error) {
	body, err := c.Request(ctx, "exportFieldNames", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var fields []struct {
		OriginalName string `json:"original_field_name"`
		ExportName   string `json:"export_field_name"`
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, fmt.Errorf("unmarshaling field names: %w", err)
	}

	result := make([]string, len(fields))
	for i, f := range fields {
		result[i] = f.ExportName
	}

	return result, nil
}
