// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestClientAuthentication tests the authentication flow.
func TestClientAuthentication(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an authentication request
		if r.URL.Path == "/login/access-token" && r.Method == http.MethodPost {
			// Verify the auth request is correct
			if err := r.ParseForm(); err != nil {
				t.Fatalf("Failed to parse form: %v", err)
			}
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username != "testuser" || password != "testpass" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Return token response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := TokenResponse{
				AccessToken: "test-token-12345",
				TokenType:   "Bearer",
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
			return
		}

		// Check auth header for other requests
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token-12345" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Mock successful response for all other requests
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	config := &Config{
		BaseURL:  server.URL,
		Username: "testuser",
		Password: "testpass",
		Timeout:  5 * time.Second,
	}

	client, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test a request that should trigger authentication
	var result map[string]interface{}
	err = client.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got %v", result["status"])
	}
}

// TestInvalidAuthentication tests authentication failure.
func TestInvalidAuthentication(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is an authentication request
		if r.URL.Path == "/login/access-token" && r.Method == http.MethodPost {
			// Return unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}))
	defer server.Close()

	// Create client with invalid credentials
	config := &Config{
		BaseURL:  server.URL,
		Username: "invalid",
		Password: "invalid",
		Timeout:  5 * time.Second,
	}

	client, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test a request that should fail authentication
	var result map[string]interface{}
	err = client.Get(context.Background(), "/test", &result)
	if err == nil {
		t.Fatalf("Expected authentication error, got nil")
	}
}

// TestClientOperations tests all HTTP methods.
func TestClientOperations(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First handle auth request
		if r.URL.Path == "/login/access-token" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := TokenResponse{
				AccessToken: "test-token-12345",
				TokenType:   "Bearer",
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
			return
		}

		// Check auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token-12345" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Respond based on the HTTP method
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var data map[string]interface{}

		switch r.Method {
		case http.MethodGet:
			data = map[string]interface{}{"method": "GET", "path": r.URL.Path}
		case http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			data = map[string]interface{}{"method": "POST", "path": r.URL.Path, "body": string(body)}
		case http.MethodPut:
			body, _ := io.ReadAll(r.Body)
			data = map[string]interface{}{"method": "PUT", "path": r.URL.Path, "body": string(body)}
		case http.MethodPatch:
			body, _ := io.ReadAll(r.Body)
			data = map[string]interface{}{"method": "PATCH", "path": r.URL.Path, "body": string(body)}
		case http.MethodDelete:
			data = map[string]interface{}{"method": "DELETE", "path": r.URL.Path}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	config := &Config{
		BaseURL:  server.URL,
		Username: "testuser",
		Password: "testpass",
		Timeout:  5 * time.Second,
	}

	client, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	var result map[string]interface{}

	// Test GET
	if err := client.Get(ctx, "/test", &result); err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if result["method"] != "GET" || result["path"] != "/test" {
		t.Errorf("Unexpected GET response: %v", result)
	}

	// Test POST
	body := strings.NewReader(`{"key":"value"}`)
	if err := client.Post(ctx, "/test", body, &result); err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if result["method"] != "POST" || result["path"] != "/test" || result["body"] != `{"key":"value"}` {
		t.Errorf("Unexpected POST response: %v", result)
	}

	// Test PUT
	body = strings.NewReader(`{"key":"updated"}`)
	if err := client.Put(ctx, "/test", body, &result); err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if result["method"] != "PUT" || result["path"] != "/test" || result["body"] != `{"key":"updated"}` {
		t.Errorf("Unexpected PUT response: %v", result)
	}

	// Test PATCH
	body = strings.NewReader(`{"key":"patched"}`)
	if err := client.Patch(ctx, "/test", body, &result); err != nil {
		t.Fatalf("PATCH request failed: %v", err)
	}
	if result["method"] != "PATCH" || result["path"] != "/test" || result["body"] != `{"key":"patched"}` {
		t.Errorf("Unexpected PATCH response: %v", result)
	}

	// Test DELETE
	if err := client.Delete(ctx, "/test", &result); err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	if result["method"] != "DELETE" || result["path"] != "/test" {
		t.Errorf("Unexpected DELETE response: %v", result)
	}
}
