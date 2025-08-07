// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceJob_Schema(t *testing.T) {
	ds := dataSourceJob()
	
	if ds == nil {
		t.Fatal("dataSourceJob() returned nil")
	}
	
	if ds.Schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Check required fields
	jobIDSchema, ok := ds.Schema["job_id"]
	if !ok {
		t.Error("job_id should be in schema")
	}
	if jobIDSchema.Type != schema.TypeInt {
		t.Error("job_id should be of type int")
	}
	if !jobIDSchema.Required {
		t.Error("job_id should be required")
	}

	// Check computed fields
	expectedComputedFields := []string{"title", "url", "enabled", "save_responses", "schedule"}
	for _, field := range expectedComputedFields {
		fieldSchema, ok := ds.Schema[field]
		if !ok {
			t.Errorf("%s should be in schema", field)
		}
		if !fieldSchema.Computed {
			t.Errorf("%s should be computed", field)
		}
	}
}

func TestDataSourceJobs_Schema(t *testing.T) {
	ds := dataSourceJobs()
	
	if ds == nil {
		t.Fatal("dataSourceJobs() returned nil")
	}
	
	if ds.Schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Check jobs field
	jobsSchema, ok := ds.Schema["jobs"]
	if !ok {
		t.Error("jobs should be in schema")
	}
	if jobsSchema.Type != schema.TypeList {
		t.Error("jobs should be of type list")
	}
	if !jobsSchema.Computed {
		t.Error("jobs should be computed")
	}
}

func TestDataSourceJobHistory_Schema(t *testing.T) {
	ds := dataSourceJobHistory()
	
	if ds == nil {
		t.Fatal("dataSourceJobHistory() returned nil")
	}
	
	if ds.Schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Check required fields
	jobIDSchema, ok := ds.Schema["job_id"]
	if !ok {
		t.Error("job_id should be in schema")
	}
	if jobIDSchema.Type != schema.TypeInt {
		t.Error("job_id should be of type int")
	}
	if !jobIDSchema.Required {
		t.Error("job_id should be required")
	}

	// Check history field
	historySchema, ok := ds.Schema["history"]
	if !ok {
		t.Error("history should be in schema")
	}
	if historySchema.Type != schema.TypeList {
		t.Error("history should be of type list")
	}
	if !historySchema.Computed {
		t.Error("history should be computed")
	}
}