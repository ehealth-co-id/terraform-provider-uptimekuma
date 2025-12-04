// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestTagOperations tests tag API operations
func TestTagOperations(t *testing.T) {
	// Setup tags for the mock server
	tags := []Tag{
		{
			ID:    1,
			Name:  "production",
			Color: "#00FF00",
		},
		{
			ID:    2,
			Name:  "development",
			Color: "#0000FF",
		},
	}

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle authentication
		if r.URL.Path == "/login/access-token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken: "test-token-12345",
				TokenType:   "Bearer",
			})
			return
		}

		// Check auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token-12345" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Handle tag requests
		if r.URL.Path == "/tags" {
			switch r.Method {
			case http.MethodGet:
				// List tags
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tags)
				return
			case http.MethodPost:
				// Create tag
				var newTag Tag
				if err := json.NewDecoder(r.Body).Decode(&newTag); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				newTag.ID = len(tags) + 1
				tags = append(tags, newTag)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(newTag)
				return
			}
		} else if strings.HasPrefix(r.URL.Path, "/tags/") {
			// Extract tag ID from path
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 3 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			idStr := parts[2]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Find tag by ID
			var tagIndex = -1
			for i, t := range tags {
				if t.ID == id {
					tagIndex = i
					break
				}
			}

			if tagIndex == -1 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			switch r.Method {
			case http.MethodGet:
				// Get tag
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tags[tagIndex])
				return
			case http.MethodDelete:
				// Delete tag
				tags = append(tags[:tagIndex], tags[tagIndex+1:]...)
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// If we get here, return 404
		w.WriteHeader(http.StatusNotFound)
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

	// Test GetTags
	retrievedTags, err := client.GetTags(ctx)
	if err != nil {
		t.Fatalf("GetTags failed: %v", err)
	}
	if len(retrievedTags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(retrievedTags))
	}
	if !reflect.DeepEqual(retrievedTags, tags) {
		t.Errorf("Tags don't match:\nExpected: %+v\nGot: %+v", tags, retrievedTags)
	}

	// Test GetTag
	tag, err := client.GetTag(ctx, 1)
	if err != nil {
		t.Fatalf("GetTag failed: %v", err)
	}
	if tag.ID != 1 || tag.Name != "production" || tag.Color != "#00FF00" {
		t.Errorf("GetTag returned unexpected result: %+v", tag)
	}

	// Test CreateTag
	newTag := &Tag{
		Name:  "testing",
		Color: "#FF0000",
	}
	createdTag, err := client.CreateTag(ctx, newTag)
	if err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}
	if createdTag.ID != 3 || createdTag.Name != "testing" || createdTag.Color != "#FF0000" {
		t.Errorf("CreateTag returned unexpected result: %+v", createdTag)
	}

	// Test DeleteTag
	if err := client.DeleteTag(ctx, 2); err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}

	// Verify tag was deleted
	retrievedTags, err = client.GetTags(ctx)
	if err != nil {
		t.Fatalf("GetTags failed after deletion: %v", err)
	}
	
	// Should have 2 tags now (ID 1 and ID 3, ID 2 was deleted)
	if len(retrievedTags) != 2 {
		t.Errorf("Expected 2 tags after deletion, got %d", len(retrievedTags))
	}
	
	// Verify the deleted tag is not in the list
	for _, tag := range retrievedTags {
		if tag.ID == 2 {
			t.Errorf("Tag ID 2 still exists after deletion")
		}
	}
}