package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExportRecords returns records from the project.
func (c *Client) ExportRecords(ctx context.Context, opts ...ExportOption) ([]Record, error) {
	params := map[string]string{
		"content": "record",
		"format":  "json",
	}

	for _, opt := range opts {
		opt(params)
	}

	body, err := c.Request(ctx, "", params)
	if err != nil {
		return nil, err
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, fmt.Errorf("unmarshaling records: %w", err)
	}

	result := make([]Record, len(records))
	for i, r := range records {
		record := Record{
			Fields: make(map[string]any),
		}
		for k, v := range r {
			switch k {
			case RecordIDField:
				record.ID = fmt.Sprintf("%v", v)
			case "redcap_event_name":
				record.EventName = fmt.Sprintf("%v", v)
			default:
				record.Fields[k] = v
			}
		}
		result[i] = record
	}

	return result, nil
}

// ExportRecordsRaw returns raw format (CSV/JSON) for records.
func (c *Client) ExportRecordsRaw(ctx context.Context, opts ...ExportOption) ([]byte, error) {
	params := map[string]string{
		"content": "record",
	}

	for _, opt := range opts {
		opt(params)
	}

	return c.Request(ctx, "", params)
}

// RecordIDField is the default field name for record IDs
const RecordIDField = "record_id"
