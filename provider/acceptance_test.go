// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJobDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJobDataSourceConfig(1), // Assuming job ID 1 exists
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cronjoborg_job.test", "title"),
					resource.TestCheckResourceAttrSet("data.cronjoborg_job.test", "url"),
					resource.TestCheckResourceAttrSet("data.cronjoborg_job.test", "enabled"),
					resource.TestCheckResourceAttrSet("data.cronjoborg_job.test", "save_responses"),
				),
			},
		},
	})
}

func TestAccJobsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJobsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cronjoborg_jobs.test", "jobs.#"),
				),
			},
		},
	})
}

func TestAccJobHistoryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJobHistoryDataSourceConfig(1), // Assuming job ID 1 exists
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cronjoborg_job_history.test", "history.#"),
				),
			},
		},
	})
}

func testAccJobDataSourceConfig(jobID int) string {
	return fmt.Sprintf(`
data "cronjoborg_job" "test" {
  job_id = %d
}
`, jobID)
}

func testAccJobsDataSourceConfig() string {
	return `
data "cronjoborg_jobs" "test" {
}
`
}

func testAccJobHistoryDataSourceConfig(jobID int) string {
	return fmt.Sprintf(`
data "cronjoborg_job_history" "test" {
  job_id = %d
}
`, jobID)
}
