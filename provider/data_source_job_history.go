// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/plain-insure/terraform-provider-cronjoborg/client"
)

func dataSourceJobHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobHistoryRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the job",
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of job execution history entries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the job",
						},
						"date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The execution date and time",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The execution status (OK, FAILED, etc.)",
						},
						"http_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The HTTP status code returned",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The execution duration in milliseconds",
						},
					},
				},
			},
		},
	}
}

func dataSourceJobHistoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, ok := m.(*client.Client)
	if !ok {
		return diag.Errorf("expected *client.Client, got %T", m)
	}

	jobIDVal, ok := d.Get("job_id").(int)
	if !ok {
		return diag.Errorf("job_id must be an integer")
	}
	jobIDStr := strconv.Itoa(jobIDVal)

	history, err := c.GetJobHistory(jobIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set a composite ID based on the job ID and number of history entries
	d.SetId(fmt.Sprintf("job-%s-history-%d", jobIDStr, len(history)))

	// Convert history to the expected format
	historyList := make([]interface{}, len(history))
	for i, entry := range history {
		historyMap := map[string]interface{}{
			"job_id":      entry.JobID,
			"date":        entry.Date,
			"status":      entry.Status,
			"http_status": entry.HttpStatus,
			"duration":    entry.Duration,
		}
		historyList[i] = historyMap
	}

	if err := d.Set("history", historyList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
