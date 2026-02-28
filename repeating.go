package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// RepeatingForm represents a repeating form or event.
type RepeatingForm struct {
	FormName          string `json:"form_name"`
	CustomRecordLabel string `json:"custom_record_label"`
}

// ExportRepeatingFormsEvents returns repeating form/event information.
func (c *Client) ExportRepeatingFormsEvents(ctx context.Context) ([]RepeatingForm, error) {
	body, err := c.Request(ctx, "repeatingFormsEvents", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var forms []RepeatingForm
	if err := json.Unmarshal(body, &forms); err != nil {
		return nil, fmt.Errorf("unmarshaling repeating forms: %w", err)
	}

	return forms, nil
}

// FormEventMapping represents form-event mappings.
type FormEventMapping struct {
	FormName        string `json:"form_name"`
	UniqueEventName string `json:"unique_event_name"`
}

// ExportFormEventMapping returns form-event mappings.
func (c *Client) ExportFormEventMapping(ctx context.Context) ([]FormEventMapping, error) {
	body, err := c.Request(ctx, "formEventMapping", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var mappings []FormEventMapping
	if err := json.Unmarshal(body, &mappings); err != nil {
		return nil, fmt.Errorf("unmarshaling form-event mappings: %w", err)
	}

	return mappings, nil
}
