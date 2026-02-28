package redcap

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExportUsers returns the list of users in the project.
func (c *Client) ExportUsers(ctx context.Context) ([]User, error) {
	body, err := c.Request(ctx, "user", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("unmarshaling users: %w", err)
	}

	return users, nil
}

// DAG represents a Data Access Group.
type DAG struct {
	UniqueGroupName string `json:"unique_group_name"`
	GroupName       string `json:"group_name"`
}

// ExportDAGs returns the list of Data Access Groups.
func (c *Client) ExportDAGs(ctx context.Context) ([]DAG, error) {
	body, err := c.Request(ctx, "dag", map[string]string{
		"format": "json",
	})
	if err != nil {
		return nil, err
	}

	var dags []DAG
	if err := json.Unmarshal(body, &dags); err != nil {
		return nil, fmt.Errorf("unmarshaling DAGs: %w", err)
	}

	return dags, nil
}
