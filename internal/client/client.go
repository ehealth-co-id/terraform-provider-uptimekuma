package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Config holds the configuration for the Uptime Kuma client
type Config struct {
	BaseURL          string
	Username         string
	Password         string
	Timeout          time.Duration
	InsecureHTTPS    bool
	CustomHTTPClient *http.Client
	DefaultHeaders   map[string]string
}

// Client is the API client for Uptime Kuma
type Client struct {
	config     *Config
	authClient *AuthClient
	httpClient *http.Client
}

// New creates a new Uptime Kuma API client
func New(config *Config) (*Client, error) {
	// Validate config
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create HTTP client
	httpClient := config.CustomHTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	// Create auth client
	authClient := NewAuthClient(
		config.BaseURL,
		config.Username,
		config.Password,
		httpClient,
	)

	// Create API client with authenticated http client
	return &Client{
		config:     config,
		authClient: authClient,
		httpClient: authClient.AuthenticatedClient(),
	}, nil
}

// doRequest performs an HTTP request and decodes the response
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, result interface{}) error {
	// Create request
	url := fmt.Sprintf("%s%s", c.config.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set default headers from config
	for k, v := range c.config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode response if result is provided
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body io.Reader, result interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body io.Reader, result interface{}) error {
	return c.doRequest(ctx, http.MethodPut, path, body, result)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body io.Reader, result interface{}) error {
	return c.doRequest(ctx, http.MethodPatch, path, body, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, result)
}
