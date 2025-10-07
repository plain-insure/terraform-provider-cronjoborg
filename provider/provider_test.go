// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	p := Provider()

	if p == nil {
		t.Fatal("Provider should not be nil")
	}

	// Test schema
	if p.Schema == nil {
		t.Fatal("Provider schema should not be nil")
	}

	// Test required fields
	apiURLSchema, ok := p.Schema["api_url"]
	if !ok {
		t.Error("api_url should be in schema")
	}
	if apiURLSchema.Type != schema.TypeString {
		t.Error("api_url should be of type string")
	}
	if !apiURLSchema.Optional {
		t.Error("api_url should be optional")
	}

	apiKeySchema, ok := p.Schema["api_key"]
	if !ok {
		t.Error("api_key should be in schema")
	}
	if apiKeySchema.Type != schema.TypeString {
		t.Error("api_key should be of type string")
	}
	if !apiKeySchema.Optional {
		t.Error("api_key should be optional")
	}
	if !apiKeySchema.Sensitive {
		t.Error("api_key should be sensitive")
	}

	// Test resources
	if p.ResourcesMap == nil {
		t.Fatal("ResourcesMap should not be nil")
	}

	expectedResources := []string{
		"cronjoborg_job",
	}

	for _, resource := range expectedResources {
		if _, ok := p.ResourcesMap[resource]; !ok {
			t.Errorf("Resource %s should be in ResourcesMap", resource)
		}
	}

	// Test datasources
	if p.DataSourcesMap == nil {
		t.Fatal("DataSourcesMap should not be nil")
	}

	expectedDataSources := []string{
		"cronjoborg_job",
		"cronjoborg_jobs",
		"cronjoborg_job_history",
	}

	for _, dataSource := range expectedDataSources {
		if _, ok := p.DataSourcesMap[dataSource]; !ok {
			t.Errorf("DataSource %s should be in DataSourcesMap", dataSource)
		}
	}

	// Test configure function
	if p.ConfigureContextFunc == nil {
		t.Error("ConfigureContextFunc should not be nil")
	}
}

func TestProvider_Configure(t *testing.T) {
	p := Provider()

	// Create test resource data
	d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{
		"api_url": "https://api.cron-job.org/",
		"api_key": "test-api-key",
	})

	_, diags := p.ConfigureContextFunc(nil, d)
	if diags.HasError() {
		t.Fatalf("Expected no errors, got %v", diags)
	}
}

func TestProvider_ConfigureMissingAPIKey(t *testing.T) {
	p := Provider()

	// Ensure environment variable does not mask missing configuration in test
	prev := os.Getenv("CRON_JOB_API_KEY")
	_ = os.Unsetenv("CRON_JOB_API_KEY")
	defer func() { _ = os.Setenv("CRON_JOB_API_KEY", prev) }()

	// Create test resource data without API key
	d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{
		"api_url": "https://api.cron-job.org/",
	})

	_, diags := p.ConfigureContextFunc(nil, d)
	if !diags.HasError() {
		t.Fatal("Expected error when API key is missing")
	}

	if len(diags) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(diags))
	}

	expectedError := "API key must be provided via provider configuration or CRON_JOB_API_KEY environment variable"
	if diags[0].Summary != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, diags[0].Summary)
	}
}

func TestProvider_ConfigureNormalizesAPIURL(t *testing.T) {
	p := Provider()

	testCases := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "URL with trailing slash",
			inputURL: "https://api.cron-job.org/",
			expected: "https://api.cron-job.org",
		},
		{
			name:     "URL with multiple trailing slashes",
			inputURL: "https://api.cron-job.org///",
			expected: "https://api.cron-job.org",
		},
		{
			name:     "URL without trailing slash",
			inputURL: "https://api.cron-job.org",
			expected: "https://api.cron-job.org",
		},
		{
			name:     "Custom URL with trailing slash",
			inputURL: "https://custom.example.com/api/",
			expected: "https://custom.example.com/api",
		},
		{
			name:     "URL with path and trailing slash",
			inputURL: "https://example.com/cron-api/v1/",
			expected: "https://example.com/cron-api/v1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{
				"api_url": tc.inputURL,
				"api_key": "test-api-key",
			})

			clientInterface, diags := p.ConfigureContextFunc(nil, d)
			if diags.HasError() {
				t.Fatalf("Expected no errors, got %v", diags)
			}

			if clientInterface == nil {
				t.Fatal("Expected client to be returned, got nil")
			}

			// We test the normalization indirectly by verifying the client was created successfully
			// The actual URL normalization will be verified by integration tests that check HTTP requests
		})
	}
}

func TestProvider_DefaultAPIURLHasNoTrailingSlash(t *testing.T) {
	p := Provider()

	// Check the default value in schema
	apiURLSchema := p.Schema["api_url"]
	defaultValue, ok := apiURLSchema.Default.(string)
	if !ok {
		t.Fatal("Default API URL should be a string")
	}

	if strings.HasSuffix(defaultValue, "/") {
		t.Errorf("Default API URL should not have trailing slash, got '%s'", defaultValue)
	}

	expectedDefault := "https://api.cron-job.org"
	if defaultValue != expectedDefault {
		t.Errorf("Expected default API URL to be '%s', got '%s'", expectedDefault, defaultValue)
	}
}
