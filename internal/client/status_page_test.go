// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestStatusPageOperations tests status page API operations.
func TestStatusPageOperations(t *testing.T) {
	writeJSON := func(w http.ResponseWriter, payload interface{}) {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		}
	}

	// Setup status pages for the mock server
	statusPages := []StatusPage{
		{
			ID:             1,
			Slug:           "main-status",
			Title:          "Main Status Page",
			Description:    "Main system status",
			Theme:          "light",
			Published:      true,
			ShowTags:       false,
			DomainNameList: []string{"status.example.com"},
			Icon:           "/icon.svg",
			ShowPoweredBy:  true,
			PublicGroupList: []PublicGroup{
				{
					ID:          1,
					Name:        "API Services",
					Weight:      1,
					MonitorList: []StatusPageMonitor{{ID: 1}, {ID: 2}, {ID: 3}},
				},
			},
		},
		{
			ID:             2,
			Slug:           "dev-status",
			Title:          "Development Status",
			Description:    "Development system status",
			Theme:          "dark",
			Published:      false,
			ShowTags:       true,
			DomainNameList: []string{"dev-status.example.com"},
			Icon:           "/dev-icon.svg",
			ShowPoweredBy:  false,
			PublicGroupList: []PublicGroup{
				{
					ID:          2,
					Name:        "Dev Services",
					Weight:      1,
					MonitorList: []StatusPageMonitor{{ID: 4}, {ID: 5}},
				},
			},
		},
	}

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle authentication
		if r.URL.Path == "/login/access-token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			writeJSON(w, TokenResponse{
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

		// Handle status page requests
		if r.URL.Path == "/statuspages" {
			switch r.Method {
			case http.MethodGet:
				// List status pages
				w.WriteHeader(http.StatusOK)
				writeJSON(w, StatusPageList{
					StatusPages: statusPages,
				})
				return
			case http.MethodPost:
				// Create status page
				var request AddStatusPageRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Check for duplicate slug
				for _, sp := range statusPages {
					if sp.Slug == request.Slug {
						w.WriteHeader(http.StatusConflict)
						return
					}
				}

				// Create new status page
				newStatusPage := StatusPage{
					ID:             len(statusPages) + 1,
					Slug:           request.Slug,
					Title:          request.Title,
					Description:    request.Msg,
					Theme:          "light",
					Published:      true,
					ShowTags:       false,
					DomainNameList: []string{},
					Icon:           "/icon.svg",
					ShowPoweredBy:  true,
				}
				statusPages = append(statusPages, newStatusPage)

				w.WriteHeader(http.StatusOK)
				writeJSON(w, AddStatusPageResponse{
					Msg: "Status page created",
				})
				return
			}
		} else if strings.HasPrefix(r.URL.Path, "/statuspages/") {
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 3 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			slug := parts[2]

			// Handle special endpoints
			if len(parts) > 3 {
				switch parts[3] {
				case "incident":
					if len(parts) > 4 && parts[4] == "unpin" {
						// Unpin incident
						if r.Method == http.MethodDelete {
							w.WriteHeader(http.StatusOK)
							writeJSON(w, UnpinIncidentResponse{
								Detail: "Incident unpinned",
							})
							return
						}
					} else {
						// Post incident
						if r.Method == http.MethodPost {
							var request PostIncidentRequest
							if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
								w.WriteHeader(http.StatusBadRequest)
								return
							}

							w.WriteHeader(http.StatusOK)
							writeJSON(w, PostIncidentResponse{
								ID:          1,
								Title:       request.Title,
								Content:     request.Content,
								Style:       request.Style,
								CreatedDate: time.Now().Format(time.RFC3339),
								Pin:         true,
							})
							return
						}
					}
				}
			}

			// Find status page by slug
			var spIndex = -1
			for i, sp := range statusPages {
				if sp.Slug == slug {
					spIndex = i
					break
				}
			}

			if spIndex == -1 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			switch r.Method {
			case http.MethodGet:
				// Get status page
				w.WriteHeader(http.StatusOK)
				writeJSON(w, statusPages[spIndex])
				return
			case http.MethodPost:
				// Update status page
				var request SaveStatusPageRequest
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if request.Title != "" {
					statusPages[spIndex].Title = request.Title
				}
				if request.Description != "" {
					statusPages[spIndex].Description = request.Description
				}
				if request.Theme != "" {
					statusPages[spIndex].Theme = request.Theme
				}
				statusPages[spIndex].Published = request.Published
				statusPages[spIndex].ShowTags = request.ShowTags
				if len(request.DomainNameList) > 0 {
					statusPages[spIndex].DomainNameList = request.DomainNameList
				}
				if request.FooterText != "" {
					statusPages[spIndex].FooterText = request.FooterText
				}
				if request.CustomCSS != "" {
					statusPages[spIndex].CustomCSS = request.CustomCSS
				}
				if request.GoogleAnalyticsID != "" {
					statusPages[spIndex].GoogleAnalyticsID = request.GoogleAnalyticsID
				}
				if request.Icon != "" {
					statusPages[spIndex].Icon = request.Icon
				}
				statusPages[spIndex].ShowPoweredBy = request.ShowPoweredBy
				if len(request.PublicGroupList) > 0 {
					statusPages[spIndex].PublicGroupList = request.PublicGroupList
				}

				w.WriteHeader(http.StatusOK)
				writeJSON(w, SaveStatusPageResponse{
					Detail: "Status page updated",
				})
				return
			case http.MethodDelete:
				// Delete status page
				statusPages = append(statusPages[:spIndex], statusPages[spIndex+1:]...)
				w.WriteHeader(http.StatusOK)
				writeJSON(w, DeleteStatusPageResponse{
					Detail: "Status page deleted",
				})
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

	// Test GetStatusPages
	retrievedPages, err := client.GetStatusPages(ctx)
	if err != nil {
		t.Fatalf("GetStatusPages failed: %v", err)
	}
	if len(retrievedPages) != len(statusPages) {
		t.Errorf("Expected %d status pages, got %d", len(statusPages), len(retrievedPages))
	}
	if !reflect.DeepEqual(retrievedPages, statusPages) {
		t.Errorf("Status pages don't match:\nExpected: %+v\nGot: %+v", statusPages, retrievedPages)
	}

	// Test GetStatusPage
	page, err := client.GetStatusPage(ctx, "main-status")
	if err != nil {
		t.Fatalf("GetStatusPage failed: %v", err)
	}
	if page.Slug != "main-status" || page.Title != "Main Status Page" {
		t.Errorf("GetStatusPage returned unexpected result: %+v", page)
	}

	// Test CreateStatusPage
	newPage := &AddStatusPageRequest{
		Slug:  "new-status",
		Title: "New Status Page",
		Msg:   "A new status page",
	}
	createResult, err := client.CreateStatusPage(ctx, newPage)
	if err != nil {
		t.Fatalf("CreateStatusPage failed: %v", err)
	}
	if createResult.Msg != "Status page created" {
		t.Errorf("CreateStatusPage returned unexpected result: %+v", createResult)
	}

	// Test UpdateStatusPage
	updatePage := &SaveStatusPageRequest{
		Title:          "Updated Status Page",
		Description:    "Updated description",
		Theme:          "dark",
		Published:      false,
		ShowTags:       true,
		DomainNameList: []string{"updated-status.example.com"},
		FooterText:     "Updated footer",
		CustomCSS:      ".header { color: blue; }",
		ShowPoweredBy:  false,
	}
	updateResult, err := client.UpdateStatusPage(ctx, "main-status", updatePage)
	if err != nil {
		t.Fatalf("UpdateStatusPage failed: %v", err)
	}
	if updateResult.Detail != "Status page updated" {
		t.Errorf("UpdateStatusPage returned unexpected result: %+v", updateResult)
	}

	// Test PostIncident
	incident := &PostIncidentRequest{
		Title:   "Test Incident",
		Content: "There is a problem with the service",
		Style:   "warning",
	}
	incidentResult, err := client.PostIncident(ctx, "main-status", incident)
	if err != nil {
		t.Fatalf("PostIncident failed: %v", err)
	}
	if incidentResult.Title != "Test Incident" || incidentResult.Style != "warning" {
		t.Errorf("PostIncident returned unexpected result: %+v", incidentResult)
	}

	// Test UnpinIncident
	unpinResult, err := client.UnpinIncident(ctx, "main-status")
	if err != nil {
		t.Fatalf("UnpinIncident failed: %v", err)
	}
	if unpinResult.Detail != "Incident unpinned" {
		t.Errorf("UnpinIncident returned unexpected result: %+v", unpinResult)
	}

	// Test DeleteStatusPage
	deleteResult, err := client.DeleteStatusPage(ctx, "dev-status")
	if err != nil {
		t.Fatalf("DeleteStatusPage failed: %v", err)
	}
	if deleteResult.Detail != "Status page deleted" {
		t.Errorf("DeleteStatusPage returned unexpected result: %+v", deleteResult)
	}
}
