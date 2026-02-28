package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

type Event struct {
	Name            string `json:"event_name"`
	ArmNum          int    `json:"arm_num"`
	DayOffset       string `json:"day_offset"`
	OffsetMin       string `json:"offset_min"`
	OffsetMax       string `json:"offset_max"`
	UniqueEventName string `json:"unique_event_name"`
}

// ExportEvents returns the list of events for longitudinal projects.
func (c *Client) ExportEvents(ctx context.Context) ([]Event, error) {
	body, err := c.Request(ctx, "event", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("unmarshaling events: %w", err)
	}

	return events, nil
}

// ExportArms returns the list of arms for longitudinal projects.
func (c *Client) ExportArms(ctx context.Context) ([]Arm, error) {
	body, err := c.Request(ctx, "arm", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var arms []Arm
	if err := json.Unmarshal(body, &arms); err != nil {
		return nil, fmt.Errorf("unmarshaling arms: %w", err)
	}

	return arms, nil
}
