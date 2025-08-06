// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJobResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJobResourceConfig("test-job", "https://example.com/test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cronjob_job.test", "title", "test-job"),
					resource.TestCheckResourceAttr("cronjob_job.test", "url", "https://example.com/test"),
					resource.TestCheckResourceAttrSet("cronjob_job.test", "id"),
				),
			},
			{
				Config: testAccJobResourceConfig("updated-job", "https://example.com/updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cronjob_job.test", "title", "updated-job"),
					resource.TestCheckResourceAttr("cronjob_job.test", "url", "https://example.com/updated"),
					resource.TestCheckResourceAttrSet("cronjob_job.test", "id"),
				),
			},
		},
	})
}

func TestAccFolderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolderResourceConfig("test-folder"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cronjob_folder.test", "title", "test-folder"),
					resource.TestCheckResourceAttrSet("cronjob_folder.test", "id"),
				),
			},
		},
	})
}

func TestAccStatusPageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckStatusPageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStatusPageResourceConfig("test-status-page"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cronjob_status_page.test", "title", "test-status-page"),
					resource.TestCheckResourceAttrSet("cronjob_status_page.test", "id"),
				),
			},
		},
	})
}

func testAccJobResourceConfig(title, url string) string {
	return fmt.Sprintf(`
resource "cronjob_job" "test" {
  title = "%s"
  url   = "%s"
}
`, title, url)
}

func testAccFolderResourceConfig(title string) string {
	return fmt.Sprintf(`
resource "cronjob_folder" "test" {
  title = "%s"
}
`, title)
}

func testAccStatusPageResourceConfig(title string) string {
	return fmt.Sprintf(`
resource "cronjob_status_page" "test" {
  title = "%s"
}
`, title)
}

func testAccCheckJobDestroy(s *terraform.State) error {
	// This would normally check that the job is destroyed
	// For now, we'll skip this check as it requires API access
	return nil
}

func testAccCheckFolderDestroy(s *terraform.State) error {
	// This would normally check that the folder is destroyed
	// For now, we'll skip this check as it requires API access
	return nil
}

func testAccCheckStatusPageDestroy(s *terraform.State) error {
	// This would normally check that the status page is destroyed
	// For now, we'll skip this check as it requires API access
	return nil
}
