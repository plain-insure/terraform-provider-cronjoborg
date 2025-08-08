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
		ReadContext: dataSourceJobsRead,
		Schema: map[string]*schema.Schema{
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
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The title of the job",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to be called by the job",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the job is enabled",
						},
						"save_responses": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to save HTTP responses",
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
									"hours": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Hours when the job should run",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"mday": {
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
									"wday": {
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

	// Convert jobs to the expected format
	jobList := make([]interface{}, len(jobs))
	for i, job := range jobs {
		jobMap := map[string]interface{}{
			"job_id":         job.JobID,
			"title":          job.Title,
			"url":            job.URL,
			"enabled":        job.Enabled,
			"save_responses": job.SaveResponses,
			"schedule": []interface{}{
				map[string]interface{}{
					"timezone": job.Schedule.Timezone,
					"hours":    job.Schedule.Hours,
					"mday":     job.Schedule.MDay,
					"minutes":  job.Schedule.Minutes,
					"months":   job.Schedule.Months,
					"wday":     job.Schedule.WDay,
				},
			},
		}
		jobList[i] = jobMap
	}

	if err := d.Set("jobs", jobList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
