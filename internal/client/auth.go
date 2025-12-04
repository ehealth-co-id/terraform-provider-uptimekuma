// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TokenResponse represents the OAuth token response from the API
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// AuthClient handles authentication with the Uptime Kuma API
type AuthClient struct {
	baseURL     string
	username    string
	password    string
	httpClient  *http.Client
	token       string
	tokenExpiry time.Time
	mutex       sync.RWMutex
}

// NewAuthClient creates a new auth client
func NewAuthClient(baseURL, username, password string, httpClient *http.Client) *AuthClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &AuthClient{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}
}

// GetToken returns a valid authentication token, refreshing if necessary
func (a *AuthClient) GetToken(ctx context.Context) (string, error) {
	a.mutex.RLock()
	token := a.token
	expiry := a.tokenExpiry
	a.mutex.RUnlock()

	// Check if we need a new token
	if token == "" || time.Now().After(expiry) {
		return a.refreshToken(ctx)
	}

	return token, nil
}

// refreshToken authenticates and gets a fresh token
func (a *AuthClient) refreshToken(ctx context.Context) (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Double check that we still need a token after acquiring the lock
	if a.token != "" && time.Now().Before(a.tokenExpiry) {
		return a.token, nil
	}

	// Prepare the authentication request
	data := url.Values{}
	data.Set("username", a.username)
	data.Set("password", a.password)

	authURL := fmt.Sprintf("%s/login/access-token", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
	}

	// Parse the token response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("received empty access token")
	}

	// Store the token and set expiry (assuming 1 hour validity, adjust as needed)
	a.token = tokenResp.AccessToken
	a.tokenExpiry = time.Now().Add(1 * time.Hour)

	return a.token, nil
}

// AddAuthHeader adds the authorization header to an HTTP request
func (a *AuthClient) AddAuthHeader(ctx context.Context, req *http.Request) error {
	token, err := a.GetToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

// AuthenticatedClient returns an http.Client that automatically handles authentication
func (a *AuthClient) AuthenticatedClient() *http.Client {
	return &http.Client{
		Transport: &authTransport{
			base:       a.httpClient.Transport,
			authClient: a,
		},
		Timeout: a.httpClient.Timeout,
	}
}

// authTransport is a custom http.RoundTripper that adds authentication headers
type authTransport struct {
	base       http.RoundTripper
	authClient *AuthClient
}

// RoundTrip implements the http.RoundTripper interface
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	req2 := req.Clone(req.Context())
	
	// Add authentication header
	if err := t.authClient.AddAuthHeader(req.Context(), req2); err != nil {
		return nil, err
	}
	
	// Use the base transport or default if none provided
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	
	return base.RoundTrip(req2)
}