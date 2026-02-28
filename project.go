package redcap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Project struct {
	ProjectID                      int       `json:"project_id"`
	ProjectTitle                   string    `json:"project_title"`
	CreationTime                   time.Time `json:"creation_time"`
	ProductionTime                 time.Time `json:"production_time"`
	Purpose                        int       `json:"purpose"`
	PurposeOther                   string    `json:"purpose_other"`
	ProjectNotes                   string    `json:"project_notes"`
	CustomRecordLabel              string    `json:"custom_record_label"`
	SecondaryUniqueField           string    `json:"secondary_unique_field"`
	IsLongitudinal                 bool      `json:"is_longitudinal"`
	HasSurveys                     bool      `json:"has_surveys"`
	HasRepetingInstrumentsOrEvents bool      `json:"has_repeating_instruments_or_events"`
	ExternalModules                []string  `json:"external_modules"`

	//Cached Metadata
	mu       sync.RWMutex
	metadata []Field
	forms    map[string]*Form
	events   []Event
	arms     []Arm
	users    []User
}

// Ping verifies API connectivity.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Request(ctx, "project", map[string]string{
		"format": "json",
	})
	return err
}

// ExportProject returns project information.
func (c *Client) ExportProject(ctx context.Context) (map[string]interface{}, error) {
	body, err := c.Request(ctx, "project", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling project: %w", err)
	}

	if len(result) == 0 {
		return nil, errors.New("no project data returned")
	}

	return result[0], nil
}
