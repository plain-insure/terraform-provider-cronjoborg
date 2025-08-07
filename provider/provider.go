// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/plain-insure/terraform-provider-cron-job.org/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.cron-job.org/",
				Description: "Base URL for the cron-job API.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CRON_JOB_API_KEY", nil),
				Description: "API key for the cron-job API. Can also be set via CRON_JOB_API_KEY env variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// Note: The cron-job.org API only supports read operations for jobs.
			// Create, update, and delete operations are not supported via the API.
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cronjob_job":         dataSourceJob(),
			"cronjob_jobs":        dataSourceJobs(),
			"cronjob_job_history": dataSourceJobHistory(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiUrl, ok := d.Get("api_url").(string)
	if !ok {
		return nil, diag.Errorf("api_url must be a string")
	}
	apiKey, ok := d.Get("api_key").(string)
	if !ok {
		return nil, diag.Errorf("api_key must be a string")
	}
	if apiKey == "" {
		apiKey = os.Getenv("CRON_JOB_API_KEY")
	}

	if apiKey == "" {
		return nil, diag.Errorf("API key must be provided via provider configuration or CRON_JOB_API_KEY environment variable")
	}

	return client.NewClient(apiUrl, apiKey), nil
}
