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