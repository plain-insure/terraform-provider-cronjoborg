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
		Description: "Fetch execution history and predictions for a specific cron job.",
		ReadContext: dataSourceJobHistoryRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the job",
			},
			"predictions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Unix timestamps of predicted next executions (up to 3)",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of job execution history entries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_log_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the history log entry",
						},
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the job",
						},
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the history item",
						},
						"date": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp of the actual execution",
						},
						"date_planned": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp of the planned execution",
						},
						"jitter": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Scheduling jitter in milliseconds",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job URL at time of execution",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The execution duration in milliseconds",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of execution",
						},
						"status_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detailed job status description",
						},
						"http_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The HTTP status code returned",
						},
						"headers": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Raw response headers returned by the host",
						},
						"body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Raw response body returned by the host",
						},
						"stats": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Additional timing information for this request",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name_lookup": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time from transfer start until name lookups completed (in microseconds)",
									},
									"connect": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time from transfer start until socket connect completed (in microseconds)",
									},
									"app_connect": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time from transfer start until SSL handshake completed (in microseconds)",
									},
									"pre_transfer": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time from transfer start until beginning of data transfer (in microseconds)",
									},
									"start_transfer": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time from transfer start until the first response byte is received (in microseconds)",
									},
									"total": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total transfer time (in microseconds)",
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

	history, predictions, err := c.GetJobHistory(jobIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set a composite ID based on the job ID and number of history entries
	d.SetId(fmt.Sprintf("job-%s-history-%d", jobIDStr, len(history)))

	// Set predictions
	if err := d.Set("predictions", predictions); err != nil {
		return diag.FromErr(err)
	}

	// Convert history to the expected format
	historyList := make([]interface{}, len(history))
	for i, entry := range history {
		stats := []interface{}{
			map[string]interface{}{
				"name_lookup":    entry.Stats.NameLookup,
				"connect":        entry.Stats.Connect,
				"app_connect":    entry.Stats.AppConnect,
				"pre_transfer":   entry.Stats.PreTransfer,
				"start_transfer": entry.Stats.StartTransfer,
				"total":          entry.Stats.Total,
			},
		}

		historyMap := map[string]interface{}{
			"job_log_id":   entry.JobLogID,
			"job_id":       entry.JobID,
			"identifier":   entry.Identifier,
			"date":         entry.Date,
			"date_planned": entry.DatePlanned,
			"jitter":       entry.Jitter,
			"url":          entry.URL,
			"duration":     entry.Duration,
			"status":       entry.Status,
			"status_text":  entry.StatusText,
			"http_status":  entry.HttpStatus,
			"stats":        stats,
		}

		// Handle nullable headers and body
		if entry.Headers != nil {
			historyMap["headers"] = *entry.Headers
		} else {
			historyMap["headers"] = ""
		}

		if entry.Body != nil {
			historyMap["body"] = *entry.Body
		} else {
			historyMap["body"] = ""
		}

		historyList[i] = historyMap
	}

	if err := d.Set("history", historyList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
