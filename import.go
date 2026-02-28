package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// ImportRecords imports records into the project.
func (c *Client) ImportRecords(ctx context.Context, records []Record, opts ...ImportOption) (*ImportResult, error) {
	params := map[string]string{
		"content": "record",
		"format":  "json",
		"type":    "flat",
	}

	for _, opt := range opts {
		opt(params)
	}

	data, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("marshaling records: %w", err)
	}

	params["data"] = string(data)

	body, err := c.Request(ctx, "", params)
	if err != nil {
		return nil, err
	}

	var result ImportResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling import result: %w", err)
	}

	return &result, nil
}

// GenerateNextRecordName generates the next sequential record name.
func (c *Client) GenerateNextRecordName(ctx context.Context) (string, error) {
	body, err := c.Request(ctx, "generateNextRecordName", map[string]string{
		"format": "json",
	})
	if err != nil {
		return "", err
	}

	var result struct {
		NextRecordName string `json:"next_record_name"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshaling next record name: %w", err)
	}

	return result.NextRecordName, nil
}
