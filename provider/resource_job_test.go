// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceJob_Schema(t *testing.T) {
	resource := resourceJob()

	if resource == nil {
		t.Fatal("Resource should not be nil")
	}

	// Test that CRUD functions are defined
	if resource.Create == nil {
		t.Error("Create function should be defined")
	}
	if resource.Read == nil {
		t.Error("Read function should be defined")
	}
	if resource.Update == nil {
		t.Error("Update function should be defined")
	}
	if resource.Delete == nil {
		t.Error("Delete function should be defined")
	}

	// Test schema fields
	if resource.Schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Test title field
	titleSchema, ok := resource.Schema["title"]
	if !ok {
		t.Error("title should be in schema")
	}
	if titleSchema.Type != schema.TypeString {
		t.Error("title should be of type string")
	}
	if !titleSchema.Required {
		t.Error("title should be required")
	}

	// Test url field
	urlSchema, ok := resource.Schema["url"]
	if !ok {
		t.Error("url should be in schema")
	}
	if urlSchema.Type != schema.TypeString {
		t.Error("url should be of type string")
	}
	if !urlSchema.Required {
		t.Error("url should be required")
	}
}

func TestBuildSchedule_NoScheduleBlock(t *testing.T) {
	// Test that when no schedule block is provided, defaults are used
	resource := resourceJob()
	resourceData := resource.TestResourceData()
	resourceData.Set("title", "Test Job")
	resourceData.Set("url", "https://example.com")

	// No schedule block set

	schedule := buildScheduleFromResourceData(resourceData)

	// Check default values
	if schedule["timezone"] != "UTC" {
		t.Errorf("Expected timezone to be UTC, got %v", schedule["timezone"])
	}
	if schedule["expiresAt"] != 0 {
		t.Errorf("Expected expiresAt to be 0, got %v", schedule["expiresAt"])
	}

	checkScheduleField := func(fieldName string, expected []int) {
		actual, ok := schedule[fieldName].([]int)
		if !ok {
			t.Errorf("Expected %s to be []int, got %T", fieldName, schedule[fieldName])
			return
		}
		if len(actual) != len(expected) {
			t.Errorf("Expected %s to have length %d, got %d", fieldName, len(expected), len(actual))
			return
		}
		for i, v := range expected {
			if actual[i] != v {
				t.Errorf("Expected %s[%d] to be %d, got %d", fieldName, i, v, actual[i])
			}
		}
	}

	checkScheduleField("hours", []int{-1})
	checkScheduleField("mdays", []int{-1})
	checkScheduleField("minutes", []int{-1})
	checkScheduleField("months", []int{-1})
	checkScheduleField("wdays", []int{-1})
}

func TestBuildSchedule_PartialScheduleBlock(t *testing.T) {
	// Test that when schedule block is provided with some fields, others default to [-1]
	resource := resourceJob()
	resourceData := resource.TestResourceData()
	resourceData.Set("title", "Test Job")
	resourceData.Set("url", "https://example.com")

	// Set schedule with only hours specified
	schedule := []interface{}{
		map[string]interface{}{
			"timezone":   "America/New_York",
			"expires_at": 20241231235959,
			"hours":      []interface{}{9, 17}, // Only hours specified
			// mdays, minutes, months, wdays not specified
		},
	}
	resourceData.Set("schedule", schedule)

	result := buildScheduleFromResourceData(resourceData)

	// Check that explicitly set values are preserved
	if result["timezone"] != "America/New_York" {
		t.Errorf("Expected timezone to be America/New_York, got %v", result["timezone"])
	}
	if result["expiresAt"] != 20241231235959 {
		t.Errorf("Expected expiresAt to be 20241231235959, got %v", result["expiresAt"])
	}

	// Check that hours are preserved
	hours, ok := result["hours"].([]int)
	if !ok {
		t.Errorf("Expected hours to be []int, got %T", result["hours"])
	} else if len(hours) != 2 || hours[0] != 9 || hours[1] != 17 {
		t.Errorf("Expected hours to be [9, 17], got %v", hours)
	}

	// Check that other fields default to [-1]
	checkScheduleField := func(fieldName string, expected []int) {
		actual, ok := result[fieldName].([]int)
		if !ok {
			t.Errorf("Expected %s to be []int, got %T", fieldName, result[fieldName])
			return
		}
		if len(actual) != len(expected) {
			t.Errorf("Expected %s to have length %d, got %d", fieldName, len(expected), len(actual))
			return
		}
		for i, v := range expected {
			if actual[i] != v {
				t.Errorf("Expected %s[%d] to be %d, got %d", fieldName, i, v, actual[i])
			}
		}
	}

	checkScheduleField("mdays", []int{-1})
	checkScheduleField("minutes", []int{-1})
	checkScheduleField("months", []int{-1})
	checkScheduleField("wdays", []int{-1})
}

func TestBuildSchedule_EmptyScheduleFields(t *testing.T) {
	// Test that when schedule fields are explicitly set to empty arrays, they default to [-1]
	resource := resourceJob()
	resourceData := resource.TestResourceData()
	resourceData.Set("title", "Test Job")
	resourceData.Set("url", "https://example.com")

	// Set schedule with empty arrays
	schedule := []interface{}{
		map[string]interface{}{
			"timezone":   "UTC",
			"expires_at": 0,
			"hours":      []interface{}{}, // Empty array
			"mdays":      []interface{}{}, // Empty array
			"minutes":    []interface{}{}, // Empty array
			"months":     []interface{}{}, // Empty array
			"wdays":      []interface{}{}, // Empty array
		},
	}
	resourceData.Set("schedule", schedule)

	result := buildScheduleFromResourceData(resourceData)

	// Check that all fields default to [-1]
	checkScheduleField := func(fieldName string, expected []int) {
		actual, ok := result[fieldName].([]int)
		if !ok {
			t.Errorf("Expected %s to be []int, got %T", fieldName, result[fieldName])
			return
		}
		if len(actual) != len(expected) {
			t.Errorf("Expected %s to have length %d, got %d", fieldName, len(expected), len(actual))
			return
		}
		for i, v := range expected {
			if actual[i] != v {
				t.Errorf("Expected %s[%d] to be %d, got %d", fieldName, i, v, actual[i])
			}
		}
	}

	checkScheduleField("hours", []int{-1})
	checkScheduleField("mdays", []int{-1})
	checkScheduleField("minutes", []int{-1})
	checkScheduleField("months", []int{-1})
	checkScheduleField("wdays", []int{-1})
}
