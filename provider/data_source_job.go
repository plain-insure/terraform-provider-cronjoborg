// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/plain-insure/terraform-provider-cronjoborg/client"
)

func dataSourceJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeInt,
				Required:    true,
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
			"auth": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "HTTP authentication settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether HTTP basic authentication is enabled",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP basic auth username",
						},
						"password": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "HTTP basic auth password",
						},
					},
				},
			},
			"notification": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Notification settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_failure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to send notification on job failure",
						},
						"on_success": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to send notification when job succeeds after prior failure",
						},
						"on_disable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to send notification when job is disabled automatically",
						},
					},
				},
			},
			"extended_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Extended request data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Request headers",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request body data",
						},
					},
				},
			},
		},
	}
}

func dataSourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*client.Client)
	if !ok {
		return diag.Errorf("expected *client.Client, got %T", m)
	}

	jobIDVal, ok := d.Get("job_id").(int)
	if !ok {
		return diag.Errorf("job_id must be an integer")
	}
	jobIDStr := strconv.Itoa(jobIDVal)

	job, err := c.GetJob(jobIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(jobIDStr)

	// Set all fields from the job
	if err := d.Set("enabled", job.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("title", job.Title); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("save_responses", job.SaveResponses); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", job.URL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_status", job.LastStatus); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_duration", job.LastDuration); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_execution", job.LastExecution); err != nil {
		return diag.FromErr(err)
	}
	if job.NextExecution != nil {
		if err := d.Set("next_execution", *job.NextExecution); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("type", job.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("request_timeout", job.RequestTimeout); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("redirect_success", job.RedirectSuccess); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("folder_id", job.FolderID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("request_method", job.RequestMethod); err != nil {
		return diag.FromErr(err)
	}

	// For data sources, we need to get the detailed job to access auth, notification, and extendedData
	detailedJob, err := c.GetJobDetails(jobIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set schedule
	schedule := []interface{}{
		map[string]interface{}{
			"timezone":   job.Schedule.Timezone,
			"expires_at": job.Schedule.ExpiresAt,
			"hours":      job.Schedule.Hours,
			"mdays":      job.Schedule.MDays,
			"minutes":    job.Schedule.Minutes,
			"months":     job.Schedule.Months,
			"wdays":      job.Schedule.WDays,
		},
	}
	if err := d.Set("schedule", schedule); err != nil {
		return diag.FromErr(err)
	}

	// Set auth
	auth := []interface{}{
		map[string]interface{}{
			"enable":   detailedJob.Auth.Enable,
			"user":     detailedJob.Auth.User,
			"password": detailedJob.Auth.Password,
		},
	}
	if err := d.Set("auth", auth); err != nil {
		return diag.FromErr(err)
	}

	// Set notification
	notification := []interface{}{
		map[string]interface{}{
			"on_failure": detailedJob.Notification.OnFailure,
			"on_success": detailedJob.Notification.OnSuccess,
			"on_disable": detailedJob.Notification.OnDisable,
		},
	}
	if err := d.Set("notification", notification); err != nil {
		return diag.FromErr(err)
	}

	// Set extended_data
	extendedData := []interface{}{
		map[string]interface{}{
			"headers": detailedJob.ExtendedData.Headers,
			"body":    detailedJob.ExtendedData.Body,
		},
	}
	if err := d.Set("extended_data", extendedData); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
