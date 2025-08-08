// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/plain-insure/terraform-provider-cronjoborg/client"
)

func dataSourceJobs() *schema.Resource {
	return &schema.Resource{
		Description: "Fetch information about all cron jobs in your account.",
		ReadContext: dataSourceJobsRead,
		Schema: map[string]*schema.Schema{
			"some_failed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if some jobs could not be retrieved due to internal errors",
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all cron jobs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the job",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the job is enabled",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The title of the job",
						},
						"save_responses": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to save HTTP responses",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to be called by the job",
						},
						"last_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last execution status",
						},
						"last_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last execution duration in milliseconds",
						},
						"last_execution": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp of last execution (in seconds)",
						},
						"next_execution": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp of predicted next execution (in seconds)",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Job type (0=Default job, 1=Monitoring job)",
						},
						"request_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Job timeout in seconds",
						},
						"redirect_success": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to treat 3xx HTTP redirect status codes as success",
						},
						"folder_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The identifier of the folder this job resides in",
						},
						"request_method": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP request method",
						},
						"schedule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The schedule configuration for the job",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The timezone for the schedule",
									},
									"expires_at": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Date/time after which the job expires",
									},
									"hours": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Hours when the job should run",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"mdays": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Days of the month when the job should run",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"minutes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Minutes when the job should run",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"months": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Months when the job should run",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"wdays": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Days of the week when the job should run",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceJobsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*client.Client)
	if !ok {
		return diag.Errorf("expected *client.Client, got %T", m)
	}

	jobs, err := c.GetJobs()
	if err != nil {
		return diag.FromErr(err)
	}

	// Set a composite ID based on the number of jobs
	d.SetId(fmt.Sprintf("jobs-%d", len(jobs)))

	// For now, we'll set some_failed to false since we don't have that info from the current GetJobs response
	// In a future enhancement, this could be updated to parse the actual API response
	if err := d.Set("some_failed", false); err != nil {
		return diag.FromErr(err)
	}

	// Convert jobs to the expected format
	jobList := make([]interface{}, len(jobs))
	for i, job := range jobs {
		jobMap := map[string]interface{}{
			"job_id":           job.JobID,
			"enabled":          job.Enabled,
			"title":            job.Title,
			"save_responses":   job.SaveResponses,
			"url":              job.URL,
			"last_status":      job.LastStatus,
			"last_duration":    job.LastDuration,
			"last_execution":   job.LastExecution,
			"type":             job.Type,
			"request_timeout":  job.RequestTimeout,
			"redirect_success": job.RedirectSuccess,
			"folder_id":        job.FolderID,
			"request_method":   job.RequestMethod,
			"schedule": []interface{}{
				map[string]interface{}{
					"timezone":   job.Schedule.Timezone,
					"expires_at": job.Schedule.ExpiresAt,
					"hours":      job.Schedule.Hours,
					"mdays":      job.Schedule.MDays,
					"minutes":    job.Schedule.Minutes,
					"months":     job.Schedule.Months,
					"wdays":      job.Schedule.WDays,
				},
			},
		}

		// Set next_execution only if it's not nil
		if job.NextExecution != nil {
			jobMap["next_execution"] = *job.NextExecution
		}

		jobList[i] = jobMap
	}

	if err := d.Set("jobs", jobList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
