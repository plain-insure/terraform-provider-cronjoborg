// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
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
		if r.URL.Path != "/jobs" {
			t.Errorf("Expected path to be /jobs, got %s", r.URL.Path)
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
				}
			]
		}`))
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
					"jobId": 1,
					"date": "2023-01-01T10:00:00Z",
					"status": "OK",
					"httpStatus": 200,
					"duration": 150
				},
				{
					"jobId": 1,
					"date": "2023-01-01T09:00:00Z",
					"status": "FAILED",
					"httpStatus": 404,
					"duration": 100
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	history, err := client.GetJobHistory("1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(history))
	}

	if history[0].Status != "OK" {
		t.Errorf("Expected first entry status to be 'OK', got '%s'", history[0].Status)
	}

	if history[1].Status != "FAILED" {
		t.Errorf("Expected second entry status to be 'FAILED', got '%s'", history[1].Status)
	}
}
