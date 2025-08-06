// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testAccProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"cronjob": func() (*schema.Provider, error) { //nolint:unparam // Required by Terraform test framework
		return Provider(), nil
	},
}

func testAccPreCheck(t *testing.T) {
	// You can add common checks here if needed
	// For example, check if required environment variables are set:
	if v := os.Getenv("CRON_JOB_API_KEY"); v == "" {
		t.Skip("CRON_JOB_API_KEY must be set for acceptance tests")
	}
}
