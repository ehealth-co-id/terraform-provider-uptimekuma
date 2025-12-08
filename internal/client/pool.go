// Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"sync"
)

// Pool manages a shared connection to Uptime Kuma for testing scenarios.
// This prevents "login: Too frequently" errors during acceptance tests by
// reusing a single Socket.IO connection across multiple provider instances.
type Pool struct {
	mu     sync.RWMutex
	client *Client
	config *Config
	refs   int // Reference counter for tracking active users
}

var (
	globalPool     *Pool
	globalPoolOnce sync.Once
	globalPoolMu   sync.Mutex
)

// GetGlobalPool returns the global connection pool instance.
// This should only be used in testing scenarios.
func GetGlobalPool() *Pool {
	globalPoolOnce.Do(func() {
		globalPool = &Pool{}
	})
	return globalPool
}

// GetOrCreate returns an existing client from the pool or creates a new one.
// If a client already exists with different configuration, an error is returned
// to prevent credential confusion.
func (p *Pool) GetOrCreate(config *Config) (*Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If we already have a client, verify config matches
	if p.client != nil {
		if !p.configMatches(config) {
			return nil, fmt.Errorf(
				"connection pool config mismatch: existing connection uses different credentials (URL: %s vs %s)",
				p.config.BaseURL, config.BaseURL,
			)
		}

		// Reuse existing connection
		p.refs++
		return p.client, nil
	}

	// Create new connection
	client, err := newClientDirect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pooled connection: %w", err)
	}

	// Store in pool
	p.client = client
	p.config = config
	p.refs = 1

	return client, nil
}

// configMatches checks if the provided config matches the pool's config.
func (p *Pool) configMatches(config *Config) bool {
	if p.config == nil {
		return false
	}
	return p.config.BaseURL == config.BaseURL &&
		p.config.Username == config.Username &&
		p.config.Password == config.Password
}

// Release decrements the reference counter for the pooled connection.
// This should be called when a client is no longer needed, but it does not
// actually close the connection (connection remains pooled for reuse).
func (p *Pool) Release() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.refs > 0 {
		p.refs--
	}
}

// Close forcefully closes the pooled connection and resets the pool.
// This should only be called during test cleanup (e.g., in TestMain).
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		err := p.client.Disconnect()
		p.client = nil
		p.config = nil
		p.refs = 0
		return err
	}

	return nil
}

// CloseGlobalPool closes the global connection pool.
// This is a convenience function for test cleanup.
func CloseGlobalPool() error {
	globalPoolMu.Lock()
	defer globalPoolMu.Unlock()

	if globalPool != nil {
		return globalPool.Close()
	}
	return nil
}
