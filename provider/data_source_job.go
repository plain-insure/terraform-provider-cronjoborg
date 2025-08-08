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

	if err := d.Set("title", job.Title); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", job.URL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", job.Enabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("save_responses", job.SaveResponses); err != nil {
		return diag.FromErr(err)
	}

	schedule := []interface{}{
		map[string]interface{}{
			"timezone": job.Schedule.Timezone,
			"hours":    job.Schedule.Hours,
			"mday":     job.Schedule.MDay,
			"minutes":  job.Schedule.Minutes,
			"months":   job.Schedule.Months,
			"wday":     job.Schedule.WDay,
		},
	}

	if err := d.Set("schedule", schedule); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
