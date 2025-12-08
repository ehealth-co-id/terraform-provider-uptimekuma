// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	kuma "github.com/breml/go-uptime-kuma-client"
)

// Config holds the configuration for the Uptime Kuma client.
type Config struct {
	BaseURL              string
	Username             string
	Password             string
	EnableConnectionPool bool // Enable connection pooling (test-only)
}

// Client is the API client for Uptime Kuma.
type Client struct {
	Kuma *kuma.Client
	// Mutex is handled internally by the library
}

// New creates a new Uptime Kuma API client.
// If connection pooling is enabled (via config or environment variable),
// it returns a shared connection from the pool. Otherwise, it creates
// a new direct connection with retry logic.
func New(config *Config) (*Client, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	// Check if connection pooling is enabled (test-only feature)
	// Priority: config flag > environment variable
	poolEnabled := config.EnableConnectionPool
	if !poolEnabled {
		// Check environment variable as fallback
		poolEnabled = (os.Getenv("UPTIMEKUMA_ENABLE_CONNECTION_POOL") == "true")
	}

	if poolEnabled {
		// Use connection pool (test scenarios)
		pool := GetGlobalPool()
		return pool.GetOrCreate(config)
	}

	// Create new direct connection (production scenarios)
	return newClientDirect(config)
}

// newClientDirect creates a new direct connection with retry logic.
// This is the original New() implementation, now extracted for reuse.
func newClientDirect(config *Config) (*Client, error) {
	ctx := context.Background() // TODO: Should we pass context in?

	// Retry configuration
	maxRetries := 5
	baseDelay := 5 * time.Second

	var k *kuma.Client
	var err error

	for i := 0; i <= maxRetries; i++ {
		k, err = kuma.New(ctx, config.BaseURL, config.Username, config.Password)
		if err == nil {
			return &Client{
				Kuma: k,
			}, nil
		}

		if i == maxRetries {
			break
		}

		// Calculate backoff with jitter
		// backoff = base * 2^i
		backoff := float64(baseDelay) * math.Pow(2, float64(i))

		// Jitter: +/- 20%
		// r = 0.8 to 1.2
		r := rand.Float64()*0.4 + 0.8
		sleepDuration := time.Duration(backoff * r)

		// Cap at 30 seconds
		if sleepDuration > 30*time.Second {
			sleepDuration = 30 * time.Second
		}

		fmt.Printf("Connection failed (attempt %d/%d): %v. Retrying in %v...\n", i+1, maxRetries+1, err, sleepDuration)
		time.Sleep(sleepDuration)
	}

	return nil, fmt.Errorf("failed to connect to Uptime Kuma after %d attempts: %w", maxRetries+1, err)
}

// Disconnect closes the connection.
func (c *Client) Disconnect() error {
	if c.Kuma != nil {
		return c.Kuma.Disconnect()
	}
	return nil
}
