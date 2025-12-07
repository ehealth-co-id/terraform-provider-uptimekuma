package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// MonitorType represents the type of monitor
type MonitorType string

// Monitor types
const (
	MonitorTypeHTTP      MonitorType = "http"
	MonitorTypePing      MonitorType = "ping"
	MonitorTypePort      MonitorType = "port"
	MonitorTypeDNS       MonitorType = "dns"
	MonitorTypeKeyword   MonitorType = "keyword"
	MonitorTypeGRPC      MonitorType = "grpc-keyword"
	MonitorTypeDocker    MonitorType = "docker"
	MonitorTypePush      MonitorType = "push"
	MonitorTypeSteam     MonitorType = "steam"
	MonitorTypeGamedig   MonitorType = "gamedig"
	MonitorTypeMQTT      MonitorType = "mqtt"
	MonitorTypeSQLServer MonitorType = "sqlserver"
	MonitorTypePostgres  MonitorType = "postgres"
	MonitorTypeMySQL     MonitorType = "mysql"
	MonitorTypeMongoDB   MonitorType = "mongodb"
	MonitorTypeRadius    MonitorType = "radius"
	MonitorTypeRedis     MonitorType = "redis"
)

// AuthMethod represents the authentication method for monitors
type AuthMethod string

// Auth methods
const (
	AuthMethodNone  AuthMethod = ""
	AuthMethodBasic AuthMethod = "basic"
	AuthMethodNTLM  AuthMethod = "ntlm"
	AuthMethodMTLS  AuthMethod = "mtls"
)

// Monitor represents an Uptime Kuma monitor
type Monitor struct {
	ID                  int           `json:"id,omitempty"`
	Type                MonitorType   `json:"type"`
	Name                string        `json:"name"`
	URL                 string        `json:"url,omitempty"`
	Method              string        `json:"method,omitempty"`
	Hostname            string        `json:"hostname,omitempty"`
	Port                int           `json:"port,omitempty"`
	Interval            int           `json:"interval"`
	RetryInterval       int           `json:"retryInterval"`
	ResendInterval      int           `json:"resendInterval"`
	MaxRetries          int           `json:"maxretries"`
	UpsideDown          bool          `json:"upsideDown"`
	NotificationIDList  []interface{} `json:"notificationIDList"`
	ExpiryNotification  bool          `json:"expiryNotification"`
	IgnoreTLS           bool          `json:"ignoreTls"`
	MaxRedirects        int           `json:"maxredirects,omitempty"`
	AcceptedStatusCodes []interface{} `json:"accepted_statuscodes,omitempty"`
	ProxyID             int           `json:"proxyId,omitempty"`
	Body                string        `json:"body,omitempty"`
	Headers             string        `json:"headers,omitempty"`
	AuthMethod          AuthMethod    `json:"authMethod,omitempty"`
	BasicAuthUser       string        `json:"basic_auth_user,omitempty"`
	BasicAuthPass       string        `json:"basic_auth_pass,omitempty"`
	AuthDomain          string        `json:"authDomain,omitempty"`
	AuthWorkstation     string        `json:"authWorkstation,omitempty"`
	Keyword             string        `json:"keyword,omitempty"`
	DNSResolveServer    string        `json:"dns_resolve_server,omitempty"`
	DNSResolveType      string        `json:"dns_resolve_type,omitempty"`
	DockerContainer     string        `json:"docker_container,omitempty"`
	DockerHost          int           `json:"docker_host,omitempty"`
}

// MonitorsResponse represents the API response for listing monitors
type MonitorsResponse struct {
	Monitors []Monitor `json:"monitors"`
}

// MonitorResponse represents the API response for a single monitor
type MonitorResponse struct {
	Monitor Monitor `json:"monitor"`
}

// MonitorCreateResponse represents the API response for creating a monitor
type MonitorCreateResponse struct {
	MonitorID int    `json:"monitorId"`
	Msg       string `json:"msg"`
}

// GetMonitors retrieves all monitors
func (c *Client) GetMonitors(ctx context.Context) ([]Monitor, error) {
	var result MonitorsResponse
	if err := c.Get(ctx, "/monitors", &result); err != nil {
		return nil, fmt.Errorf("failed to get monitors: %w", err)
	}
	return result.Monitors, nil
}

// GetMonitor retrieves a specific monitor by ID
func (c *Client) GetMonitor(ctx context.Context, id int) (*Monitor, error) {
	var result MonitorResponse
	path := fmt.Sprintf("/monitors/%d", id)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get monitor %d: %w", id, err)
	}
	return &result.Monitor, nil
}

// CreateMonitor creates a new monitor
func (c *Client) CreateMonitor(ctx context.Context, monitor *Monitor) (*Monitor, error) {
	data, err := json.Marshal(monitor)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal monitor: %w", err)
	}

	var result MonitorCreateResponse
	if err := c.Post(ctx, "/monitors", bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("failed to create monitor: %w", err)
	}
	
	// Return the monitor with the ID from the response
	monitor.ID = result.MonitorID
	return monitor, nil
}

// UpdateMonitor updates an existing monitor
func (c *Client) UpdateMonitor(ctx context.Context, id int, monitor *Monitor) (*Monitor, error) {
	data, err := json.Marshal(monitor)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal monitor: %w", err)
	}

	var result MonitorResponse
	path := fmt.Sprintf("/monitors/%d", id)
	if err := c.Patch(ctx, path, bytes.NewReader(data), &result); err != nil {
		return nil, fmt.Errorf("failed to update monitor %d: %w", id, err)
	}
	return &result.Monitor, nil
}

// DeleteMonitor deletes a monitor
func (c *Client) DeleteMonitor(ctx context.Context, id int) error {
	path := fmt.Sprintf("/monitors/%d", id)
	if err := c.Delete(ctx, path, nil); err != nil {
		return fmt.Errorf("failed to delete monitor %d: %w", id, err)
	}
	return nil
}

// PauseMonitor pauses a monitor
func (c *Client) PauseMonitor(ctx context.Context, id int) error {
	path := fmt.Sprintf("/monitors/%d/pause", id)
	if err := c.Post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to pause monitor %d: %w", id, err)
	}
	return nil
}

// ResumeMonitor resumes a paused monitor
func (c *Client) ResumeMonitor(ctx context.Context, id int) error {
	path := fmt.Sprintf("/monitors/%d/resume", id)
	if err := c.Post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to resume monitor %d: %w", id, err)
	}
	return nil
}

// GetMonitorBeats retrieves the heartbeats for a monitor
func (c *Client) GetMonitorBeats(ctx context.Context, id int, hours float64) (interface{}, error) {
	path := fmt.Sprintf("/monitors/%d/beats?hours=%s", id, strconv.FormatFloat(hours, 'f', -1, 64))
	var result interface{}
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("failed to get beats for monitor %d: %w", id, err)
	}
	return result, nil
}

// AddMonitorTag adds a tag to a monitor
func (c *Client) AddMonitorTag(ctx context.Context, monitorID int, tagID int, value string) error {
	tag := struct {
		TagID int    `json:"tag_id"`
		Value string `json:"value,omitempty"`
	}{
		TagID: tagID,
		Value: value,
	}

	data, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("failed to marshal tag: %w", err)
	}

	path := fmt.Sprintf("/monitors/%d/tag", monitorID)
	if err := c.Post(ctx, path, bytes.NewReader(data), nil); err != nil {
		return fmt.Errorf("failed to add tag %d to monitor %d: %w", tagID, monitorID, err)
	}
	return nil
}

// DeleteMonitorTag removes a tag from a monitor
func (c *Client) DeleteMonitorTag(ctx context.Context, monitorID int, tagID int) error {
	tag := struct {
		TagID int `json:"tag_id"`
	}{
		TagID: tagID,
	}

	data, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("failed to marshal tag: %w", err)
	}

	path := fmt.Sprintf("/monitors/%d/tag", monitorID)
	if err := c.Delete(ctx, path, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("failed to delete tag %d from monitor %d: %w", tagID, monitorID, err)
	}
	return nil
}