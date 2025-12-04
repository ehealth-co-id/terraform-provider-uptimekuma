// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// Tag represents an Uptime Kuma tag.
type Tag struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// GetTags retrieves all tags.
func (c *Client) GetTags(ctx context.Context) ([]Tag, error) {
	var result []Tag
	if err := c.Get(ctx, "/tags", &result); err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return result, nil
}

// GetTag retrieves a specific tag by ID.
func (c *Client) GetTag(ctx context.Context, id int) (*Tag, error) {
	var result Tag
	path := fmt.Sprintf("/tags/%d", id)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get tag %d: %w", id, err)
	}
	return &result, nil
}

// CreateTag creates a new tag.
func (c *Client) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	data, err := json.Marshal(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tag: %w", err)
	}

	var result Tag
	if err := c.Post(ctx, "/tags", bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return &result, nil
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(ctx context.Context, id int) error {
	path := fmt.Sprintf("/tags/%d", id)
	if err := c.Delete(ctx, path, nil); err != nil {
		return fmt.Errorf("failed to delete tag %d: %w", id, err)
	}
	return nil
}
