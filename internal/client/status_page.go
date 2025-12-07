// Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// StatusPageMonitor represents a monitor reference in a status page group.
type StatusPageMonitor struct {
	ID int `json:"id"`
}

// PublicGroup represents a group of monitors on a status page.
type PublicGroup struct {
	ID          int                 `json:"id,omitempty"`
	Name        string              `json:"name"`
	Weight      int                 `json:"weight"`
	MonitorList []StatusPageMonitor `json:"monitorList"`
}

// StatusPage represents an Uptime Kuma status page.
type StatusPage struct {
	ID                int           `json:"id,omitempty"`
	Slug              string        `json:"slug"`
	Title             string        `json:"title"`
	Description       string        `json:"description,omitempty"`
	Theme             string        `json:"theme"`
	Published         bool          `json:"published"`
	ShowTags          bool          `json:"showTags"`
	DomainNameList    []string      `json:"domainNameList"`
	FooterText        string        `json:"footerText,omitempty"`
	CustomCSS         string        `json:"customCSS,omitempty"`
	GoogleAnalyticsID string        `json:"googleAnalyticsId,omitempty"`
	Icon              string        `json:"icon,omitempty"`
	ShowPoweredBy     bool          `json:"showPoweredBy"`
	PublicGroupList   []PublicGroup `json:"publicGroupList,omitempty"`
}

// StatusPageList represents a list of status pages.
type StatusPageList struct {
	StatusPages []StatusPage `json:"statuspages"`
}

// AddStatusPageRequest represents the request to create a status page.
type AddStatusPageRequest struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Msg   string `json:"msg,omitempty"`
}

// AddStatusPageResponse represents the response from creating a status page.
type AddStatusPageResponse struct {
	Msg string `json:"msg"`
}

// SaveStatusPageRequest represents the request to update a status page.
type SaveStatusPageRequest struct {
	Title             string        `json:"title,omitempty"`
	Description       string        `json:"description,omitempty"`
	Theme             string        `json:"theme,omitempty"`
	Published         bool          `json:"published,omitempty"`
	ShowTags          bool          `json:"showTags,omitempty"`
	DomainNameList    []string      `json:"domainNameList,omitempty"`
	FooterText        string        `json:"footerText,omitempty"`
	CustomCSS         string        `json:"customCSS,omitempty"`
	GoogleAnalyticsID string        `json:"googleAnalyticsId,omitempty"`
	Icon              string        `json:"icon,omitempty"`
	ShowPoweredBy     bool          `json:"showPoweredBy,omitempty"`
	PublicGroupList   []PublicGroup `json:"publicGroupList,omitempty"`
}

// SaveStatusPageResponse represents the response from updating a status page.
type SaveStatusPageResponse struct {
	Detail interface{} `json:"detail"`
}

// DeleteStatusPageResponse represents the response from deleting a status page.
type DeleteStatusPageResponse struct {
	Detail string `json:"detail"`
}

// GetStatusPages retrieves all status pages.
func (c *Client) GetStatusPages(ctx context.Context) ([]StatusPage, error) {
	var result StatusPageList
	if err := c.Get(ctx, "/statuspages", &result); err != nil {
		return nil, fmt.Errorf("failed to get status pages: %w", err)
	}
	return result.StatusPages, nil
}

// GetStatusPage retrieves a specific status page by slug.
func (c *Client) GetStatusPage(ctx context.Context, slug string) (*StatusPage, error) {
	var result StatusPage
	path := fmt.Sprintf("/statuspages/%s", slug)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get status page %s: %w", slug, err)
	}
	return &result, nil
}

// CreateStatusPage creates a new status page.
func (c *Client) CreateStatusPage(ctx context.Context, request *AddStatusPageRequest) (*AddStatusPageResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status page: %w", err)
	}

	var result AddStatusPageResponse
	if err := c.Post(ctx, "/statuspages", bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("failed to create status page: %w", err)
	}
	return &result, nil
}

// UpdateStatusPage updates an existing status page.
func (c *Client) UpdateStatusPage(ctx context.Context, slug string, request *SaveStatusPageRequest) (*SaveStatusPageResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status page: %w", err)
	}

	var result SaveStatusPageResponse
	path := fmt.Sprintf("/statuspages/%s", slug)
	if err := c.Post(ctx, path, bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("failed to update status page %s: %w", slug, err)
	}
	return &result, nil
}

// DeleteStatusPage deletes a status page.
func (c *Client) DeleteStatusPage(ctx context.Context, slug string) (*DeleteStatusPageResponse, error) {
	var result DeleteStatusPageResponse
	path := fmt.Sprintf("/statuspages/%s", slug)
	if err := c.Delete(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to delete status page %s: %w", slug, err)
	}
	return &result, nil
}

// PostIncidentRequest represents a request to post an incident to a status page.
type PostIncidentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Style   string `json:"style,omitempty"` // primary, info, warning, danger, light, dark
}

// PostIncidentResponse represents the response from posting an incident.
type PostIncidentResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Style       string `json:"style"`
	CreatedDate string `json:"createdDate"`
	Pin         bool   `json:"pin"`
}

// UnpinIncidentResponse represents the response from unpinning an incident.
type UnpinIncidentResponse struct {
	Detail string `json:"detail"`
}

// PostIncident posts an incident to a status page.
func (c *Client) PostIncident(ctx context.Context, slug string, request *PostIncidentRequest) (*PostIncidentResponse, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal incident: %w", err)
	}

	var result PostIncidentResponse
	path := fmt.Sprintf("/statuspages/%s/incident", slug)
	if err := c.Post(ctx, path, bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("failed to post incident to status page %s: %w", slug, err)
	}
	return &result, nil
}

// UnpinIncident unpins an incident from a status page.
func (c *Client) UnpinIncident(ctx context.Context, slug string) (*UnpinIncidentResponse, error) {
	var result UnpinIncidentResponse
	path := fmt.Sprintf("/statuspages/%s/incident/unpin", slug)
	if err := c.Delete(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to unpin incident from status page %s: %w", slug, err)
	}
	return &result, nil
}
