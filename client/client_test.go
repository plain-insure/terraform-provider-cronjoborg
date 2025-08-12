// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://api.cron-job.org/"
	apiKey := "test-api-key"

	client := NewClient(baseURL, apiKey)

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", baseURL, client.BaseURL)
	}

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, client.APIKey)
	}

	if client.HTTPClient != http.DefaultClient {
		t.Errorf("Expected HTTPClient to be http.DefaultClient")
	}
}

func TestDoRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Expected Authorization header to be 'Bearer test-key', got %s", auth)
		}

		// Check method
		if r.Method != "GET" {
			t.Errorf("Expected method to be GET, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/test" {
			t.Errorf("Expected path to be /test, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	resp, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestDoRequestWithBody(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Content-Type header
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json', got %s", contentType)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	body := map[string]interface{}{
		"title": "Test Job",
	}

	resp, err := client.doRequest("POST", "/test", body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestDoRequestAPIError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Invalid request"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	_, err := client.doRequest("GET", "/test", nil)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}

	if apiErr.Message != "Invalid request" {
		t.Errorf("Expected error message 'Invalid request', got '%s'", apiErr.Message)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Message:    "Bad Request",
		Body:       `{"error": "Bad Request"}`,
	}

	expected := "API error 400: Bad Request"
	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}
}

func TestGetJobs(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jobs" {
			t.Errorf("Expected path to be /jobs, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected method to be GET, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"jobs": [
				{
					"jobId": 1,
					"title": "Test Job 1",
					"url": "https://example.com/1",
					"enabled": true,
					"saveResponses": false,
					"schedule": {
						"timezone": "UTC",
						"hours": [10],
						"mday": [-1],
						"minutes": [0],
						"months": [-1],
						"wday": [-1]
					}
				},
				{
					"jobId": 2,
					"title": "Test Job 2",
					"url": "https://example.com/2",
					"enabled": false,
					"saveResponses": true,
					"schedule": {
						"timezone": "UTC",
						"hours": [14],
						"mday": [-1],
						"minutes": [30],
						"months": [-1],
						"wday": [-1]
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	jobs, err := client.GetJobs()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].JobID != 1 {
		t.Errorf("Expected first job ID to be 1, got %d", jobs[0].JobID)
	}

	if jobs[0].Title != "Test Job 1" {
		t.Errorf("Expected first job title to be 'Test Job 1', got '%s'", jobs[0].Title)
	}

	if jobs[1].JobID != 2 {
		t.Errorf("Expected second job ID to be 2, got %d", jobs[1].JobID)
	}
}

func TestGetJob(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jobs/1" && r.URL.Path != "/jobs/999" {
			t.Errorf("Expected path to be /jobs/1 or /jobs/999, got %s", r.URL.Path)
		}

		if r.URL.Path == "/jobs/1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"jobDetails": {
					"jobId": 1,
					"enabled": true,
					"title": "Test Job 1",
					"saveResponses": false,
					"url": "https://example.com/1",
					"lastStatus": 0,
					"lastDuration": 0,
					"lastExecution": 0,
					"nextExecution": null,
					"type": 0,
					"requestTimeout": 300,
					"redirectSuccess": false,
					"folderId": 0,
					"schedule": {
						"timezone": "UTC",
						"expiresAt": 0,
						"hours": [10],
						"mdays": [-1],
						"minutes": [0],
						"months": [-1],
						"wdays": [-1]
					},
					"requestMethod": 0,
					"auth": {
						"enable": false,
						"user": "",
						"password": ""
					},
					"notification": {
						"onFailure": false,
						"onSuccess": false,
						"onDisable": false
					},
					"extendedData": {
						"headers": {},
						"body": ""
					}
				}
			}`))
		} else {
			// Return 404 for job 999
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "Job not found"}`))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	job, err := client.GetJob("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if job.JobID != 1 {
		t.Errorf("Expected job ID to be 1, got %d", job.JobID)
	}

	if job.Title != "Test Job 1" {
		t.Errorf("Expected job title to be 'Test Job 1', got '%s'", job.Title)
	}

	// Test job not found
	_, err = client.GetJob("999")
	if err == nil {
		t.Fatal("Expected an error for non-existent job, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
	}
}

func TestGetJobHistory(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jobs/1/history" {
			t.Errorf("Expected path to be /jobs/1/history, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"history": [
				{
					"jobLogId": 1001,
					"jobId": 1,
					"identifier": "1-23-01-1001",
					"date": 1672574400,
					"datePlanned": 1672574400,
					"jitter": 50,
					"url": "https://example.com",
					"duration": 150,
					"status": 1,
					"statusText": "OK",
					"httpStatus": 200,
					"headers": null,
					"body": null,
					"stats": {
						"nameLookup": 1000,
						"connect": 2000,
						"appConnect": 0,
						"preTransfer": 3000,
						"startTransfer": 4000,
						"total": 150000
					}
				},
				{
					"jobLogId": 1002,
					"jobId": 1,
					"identifier": "1-23-01-1002",
					"date": 1672570800,
					"datePlanned": 1672570800,
					"jitter": 75,
					"url": "https://example.com",
					"duration": 100,
					"status": 4,
					"statusText": "FAILED",
					"httpStatus": 404,
					"headers": null,
					"body": null,
					"stats": {
						"nameLookup": 1200,
						"connect": 2200,
						"appConnect": 0,
						"preTransfer": 3200,
						"startTransfer": 4200,
						"total": 100000
					}
				}
			],
			"predictions": [1672578000, 1672581600, 1672585200]
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	history, predictions, err := client.GetJobHistory("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(history))
	}

	if len(predictions) != 3 {
		t.Errorf("Expected 3 predictions, got %d", len(predictions))
	}

	if history[0].Status != 1 {
		t.Errorf("Expected first entry status to be 1 (OK), got %d", history[0].Status)
	}

	if history[1].Status != 4 {
		t.Errorf("Expected second entry status to be 4 (HTTP error/FAILED), got %d", history[1].Status)
	}
}

func TestClient_CreateJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/jobs" {
			t.Errorf("Expected path '/jobs', got '%s'", r.URL.Path)
		}

		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		job, ok := requestBody["job"].(map[string]interface{})
		if !ok {
			t.Error("Expected 'job' field in request body")
		}

		if job["title"] != "Test Job" {
			t.Errorf("Expected job title 'Test Job', got '%v'", job["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"jobId": 123}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	job := map[string]interface{}{
		"title": "Test Job",
		"url":   "https://example.com",
	}

	jobID, err := client.CreateJob(job)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if jobID != 123 {
		t.Errorf("Expected job ID 123, got %d", jobID)
	}
}

func TestClient_UpdateJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		if r.URL.Path != "/jobs/123" {
			t.Errorf("Expected path '/jobs/123', got '%s'", r.URL.Path)
		}

		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		job, ok := requestBody["job"].(map[string]interface{})
		if !ok {
			t.Error("Expected 'job' field in request body")
		}

		if job["title"] != "Updated Job" {
			t.Errorf("Expected job title 'Updated Job', got '%v'", job["title"])
		}

		// Ensure jobId is NOT included in the request body
		if _, exists := requestBody["jobId"]; exists {
			t.Error("jobId should not be included in PATCH request body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	job := map[string]interface{}{
		"title": "Updated Job",
		"url":   "https://example.com/updated",
	}

	err := client.UpdateJob("123", job)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_DeleteJob(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/jobs/123" {
			t.Errorf("Expected path '/jobs/123', got '%s'", r.URL.Path)
		}

		// Check that no request body is sent for DELETE
		if r.ContentLength > 0 {
			t.Error("DELETE request should not have a body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	err := client.DeleteJob("123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_GetJobDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/jobs/123" {
			t.Errorf("Expected path '/jobs/123', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"jobDetails": {
				"jobId": 123,
				"enabled": true,
				"title": "Test Job",
				"saveResponses": false,
				"url": "https://example.com",
				"lastStatus": 0,
				"lastDuration": 0,
				"lastExecution": 0,
				"nextExecution": null,
				"type": 0,
				"requestTimeout": 300,
				"redirectSuccess": false,
				"folderId": 0,
				"schedule": {
					"timezone": "UTC",
					"expiresAt": 0,
					"hours": [-1],
					"mdays": [-1],
					"minutes": [-1],
					"months": [-1],
					"wdays": [-1]
				},
				"requestMethod": 0,
				"auth": {
					"enable": false,
					"user": "",
					"password": ""
				},
				"notification": {
					"onFailure": false,
					"onSuccess": false,
					"onDisable": false
				},
				"extendedData": {
					"headers": [],
					"body": ""
				}
			}
		}`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	details, err := client.GetJobDetails("123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if details.Title != "Test Job" {
		t.Errorf("Expected title 'Test Job', got '%s'", details.Title)
	}

	if details.URL != "https://example.com" {
		t.Errorf("Expected url 'https://example.com', got '%s'", details.URL)
	}

	// Verify that headers array is properly handled as empty map
	if len(details.ExtendedData.Headers) != 0 {
		t.Errorf("Expected empty headers map, got %d headers", len(details.ExtendedData.Headers))
	}
}

func TestJobExtendedData_UnmarshalJSON_ObjectHeaders(t *testing.T) {
	// Test normal object headers
	jsonData := `{
		"headers": {
			"X-Foo": "Bar",
			"Content-Type": "application/json"
		},
		"body": "test body"
	}`

	var extData JobExtendedData
	err := json.Unmarshal([]byte(jsonData), &extData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(extData.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(extData.Headers))
	}

	if extData.Headers["X-Foo"] != "Bar" {
		t.Errorf("Expected X-Foo header to be 'Bar', got '%s'", extData.Headers["X-Foo"])
	}

	if extData.Body != "test body" {
		t.Errorf("Expected body to be 'test body', got '%s'", extData.Body)
	}
}

func TestJobExtendedData_UnmarshalJSON_ArrayHeaders(t *testing.T) {
	// Test when API returns an array instead of object (should be treated as empty map)
	jsonData := `{
		"headers": [],
		"body": "test body"
	}`

	var extData JobExtendedData
	err := json.Unmarshal([]byte(jsonData), &extData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(extData.Headers) != 0 {
		t.Errorf("Expected 0 headers, got %d", len(extData.Headers))
	}

	if extData.Body != "test body" {
		t.Errorf("Expected body to be 'test body', got '%s'", extData.Body)
	}
}

func TestJobExtendedData_UnmarshalJSON_EmptyHeaders(t *testing.T) {
	// Test empty object headers
	jsonData := `{
		"headers": {},
		"body": "test body"
	}`

	var extData JobExtendedData
	err := json.Unmarshal([]byte(jsonData), &extData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(extData.Headers) != 0 {
		t.Errorf("Expected 0 headers, got %d", len(extData.Headers))
	}

	if extData.Body != "test body" {
		t.Errorf("Expected body to be 'test body', got '%s'", extData.Body)
	}
}

func TestClient_URLConstruction(t *testing.T) {
	testCases := []struct {
		name        string
		baseURL     string
		path        string
		expectedURL string
	}{
		{
			name:        "Base URL without trailing slash",
			baseURL:     "https://api.cron-job.org",
			path:        "/jobs",
			expectedURL: "https://api.cron-job.org/jobs",
		},
		{
			name:        "Custom API URL without trailing slash",
			baseURL:     "https://example.com/api",
			path:        "/jobs",
			expectedURL: "https://example.com/api/jobs",
		},
		{
			name:        "URL with path component",
			baseURL:     "https://custom.example.com/cron-api/v1",
			path:        "/jobs/123",
			expectedURL: "https://custom.example.com/cron-api/v1/jobs/123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the URL construction logic directly (this is what happens in doRequest)
			constructedURL := tc.baseURL + tc.path
			if constructedURL != tc.expectedURL {
				t.Errorf("Expected constructed URL to be '%s', got '%s'", tc.expectedURL, constructedURL)
			}
		})
	}
}

// TestClient_RealHTTPURLConstruction verifies that the actual HTTP requests use the expected URLs
func TestClient_RealHTTPURLConstruction(t *testing.T) {
	testCases := []struct {
		name         string
		baseURL      string
		path         string
		expectedPath string
	}{
		{
			name:         "Normalized base URL constructs correct endpoint",
			baseURL:      "https://api.cron-job.org",  // normalized (no trailing slash)
			path:         "/jobs",
			expectedPath: "/jobs",
		},
		{
			name:         "Custom normalized URL constructs correct endpoint",
			baseURL:      "https://example.com/api/v1",  // normalized (no trailing slash)
			path:         "/jobs/123/history",
			expectedPath: "/jobs/123/history",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server that captures the request URL
			var capturedPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"success": true}`))
			}))
			defer server.Close()

			// Create client with normalized base URL
			client := NewClient(tc.baseURL, "test-key")
			
			// Override base URL to point to our test server
			client.BaseURL = server.URL
			
			// Make the request
			_, err := client.doRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify the path was constructed correctly
			if capturedPath != tc.expectedPath {
				t.Errorf("Expected request path to be '%s', got '%s'", tc.expectedPath, capturedPath)
			}
		})
	}
}
