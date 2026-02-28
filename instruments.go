package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExportInstruments returns the list of instruments/forms in the project.
func (c *Client) ExportInstruments(ctx context.Context) ([]Instrument, error) {
	body, err := c.Request(ctx, "instrument", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var instruments []Instrument
	if err := json.Unmarshal(body, &instruments); err != nil {
		return nil, fmt.Errorf("unmarshaling instruments: %w", err)
	}

	return instruments, nil
}
