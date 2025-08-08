// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/plain-insure/terraform-provider-cronjoborg/client"
)

func resourceJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobCreate,
		Read:   resourceJobRead,
		Update: resourceJobUpdate,
		Delete: resourceJobDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The title of the cron job",
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL to be called by the cron job",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value, ok := v.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("%q must be a string", k))
						return
					}
					if _, err := url.ParseRequestURI(value); err != nil {
						errors = append(errors, fmt.Errorf("%q must be a valid URL: %s", k, err))
					}
					return
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the job is enabled (i.e. being executed) or not",
			},
			"save_responses": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to save job response header/body or not",
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				Description:  "Job timeout in seconds (-1 = use default timeout)",
				ValidateFunc: validation.IntAtLeast(-1),
			},
			"redirect_success": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to treat 3xx HTTP redirect status codes as success or not",
			},
			"folder_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "The identifier of the folder this job resides in (0 = root folder)",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"request_method": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "HTTP request method (0=GET, 1=POST, 2=OPTIONS, 3=HEAD, 4=PUT, 5=DELETE, 6=TRACE, 7=CONNECT, 8=PATCH)",
				ValidateFunc: validation.IntBetween(0, 8),
			},
			"schedule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Job schedule configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timezone": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "UTC",
							Description: "Schedule time zone",
						},
						"expires_at": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							Description:  "Date/time after which the job expires (format: YYYYMMDDhhmmss, 0 = does not expire)",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"hours": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Hours in which to execute the job (0-23; [-1] = every hour)",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(-1, 23),
							},
						},
						"mdays": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Days of month in which to execute the job (1-31; [-1] = every day of month)",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(-1, 31),
							},
						},
						"minutes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Minutes in which to execute the job (0-59; [-1] = every minute)",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(-1, 59),
							},
						},
						"months": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Months in which to execute the job (1-12; [-1] = every month)",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(-1, 12),
							},
						},
						"wdays": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Days of week in which to execute the job (0=Sunday-6=Saturday; [-1] = every day of week)",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(-1, 6),
							},
						},
					},
				},
			},
			"auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "HTTP authentication settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable HTTP basic authentication or not",
						},
						"user": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "HTTP basic auth username",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Sensitive:   true,
							Description: "HTTP basic auth password",
						},
					},
				},
			},
			"notification": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Notification settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"on_failure": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to send a notification on job failure or not",
						},
						"on_success": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to send a notification when the job succeeds after a prior failure or not",
						},
						"on_disable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to send a notification when the job has been disabled automatically or not",
						},
					},
				},
			},
			"extended_data": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Extended request data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Request headers (key-value dictionary)",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"body": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Request body data",
						},
					},
				},
			},
			// Computed fields for read-only values
			"job_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the job",
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
		},
	}
}

func resourceJobCreate(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}

	// Build job object from schema
	job := make(map[string]interface{})

	// Required fields
	job["title"] = d.Get("title").(string)
	job["url"] = d.Get("url").(string)

	// Optional fields with defaults
	job["enabled"] = d.Get("enabled").(bool)
	job["saveResponses"] = d.Get("save_responses").(bool)
	job["requestTimeout"] = d.Get("request_timeout").(int)
	job["redirectSuccess"] = d.Get("redirect_success").(bool)
	job["folderId"] = d.Get("folder_id").(int)
	job["requestMethod"] = d.Get("request_method").(int)

	// Schedule
	if scheduleList := d.Get("schedule").([]interface{}); len(scheduleList) > 0 {
		scheduleMap := scheduleList[0].(map[string]interface{})
		schedule := make(map[string]interface{})

		schedule["timezone"] = scheduleMap["timezone"].(string)
		schedule["expiresAt"] = scheduleMap["expires_at"].(int)

		if hours := scheduleMap["hours"].([]interface{}); len(hours) > 0 {
			hoursInt := make([]int, len(hours))
			for i, v := range hours {
				hoursInt[i] = v.(int)
			}
			schedule["hours"] = hoursInt
		} else {
			schedule["hours"] = []int{}
		}

		if mdays := scheduleMap["mdays"].([]interface{}); len(mdays) > 0 {
			mdaysInt := make([]int, len(mdays))
			for i, v := range mdays {
				mdaysInt[i] = v.(int)
			}
			schedule["mdays"] = mdaysInt
		} else {
			schedule["mdays"] = []int{}
		}

		if minutes := scheduleMap["minutes"].([]interface{}); len(minutes) > 0 {
			minutesInt := make([]int, len(minutes))
			for i, v := range minutes {
				minutesInt[i] = v.(int)
			}
			schedule["minutes"] = minutesInt
		} else {
			schedule["minutes"] = []int{}
		}

		if months := scheduleMap["months"].([]interface{}); len(months) > 0 {
			monthsInt := make([]int, len(months))
			for i, v := range months {
				monthsInt[i] = v.(int)
			}
			schedule["months"] = monthsInt
		} else {
			schedule["months"] = []int{}
		}

		if wdays := scheduleMap["wdays"].([]interface{}); len(wdays) > 0 {
			wdaysInt := make([]int, len(wdays))
			for i, v := range wdays {
				wdaysInt[i] = v.(int)
			}
			schedule["wdays"] = wdaysInt
		} else {
			schedule["wdays"] = []int{}
		}

		job["schedule"] = schedule
	}

	// Auth
	if authList := d.Get("auth").([]interface{}); len(authList) > 0 {
		authMap := authList[0].(map[string]interface{})
		auth := map[string]interface{}{
			"enable":   authMap["enable"].(bool),
			"user":     authMap["user"].(string),
			"password": authMap["password"].(string),
		}
		job["auth"] = auth
	}

	// Notification
	if notificationList := d.Get("notification").([]interface{}); len(notificationList) > 0 {
		notificationMap := notificationList[0].(map[string]interface{})
		notification := map[string]interface{}{
			"onFailure": notificationMap["on_failure"].(bool),
			"onSuccess": notificationMap["on_success"].(bool),
			"onDisable": notificationMap["on_disable"].(bool),
		}
		job["notification"] = notification
	}

	// Extended Data
	if extendedDataList := d.Get("extended_data").([]interface{}); len(extendedDataList) > 0 {
		extendedDataMap := extendedDataList[0].(map[string]interface{})
		extendedData := make(map[string]interface{})

		if headers := extendedDataMap["headers"].(map[string]interface{}); len(headers) > 0 {
			stringHeaders := make(map[string]string)
			for k, v := range headers {
				stringHeaders[k] = v.(string)
			}
			extendedData["headers"] = stringHeaders
		} else {
			extendedData["headers"] = map[string]string{}
		}

		extendedData["body"] = extendedDataMap["body"].(string)
		job["extendedData"] = extendedData
	}

	jobID, err := c.CreateJob(job)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", jobID))
	return resourceJobRead(d, m)
}

func resourceJobRead(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}

	jobDetails, err := c.GetJobDetails(d.Id())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	// Set all fields from the detailed job
	if err := d.Set("job_id", jobDetails.JobID); err != nil {
		return fmt.Errorf("error setting job_id: %w", err)
	}
	if err := d.Set("title", jobDetails.Title); err != nil {
		return fmt.Errorf("error setting title: %w", err)
	}
	if err := d.Set("url", jobDetails.URL); err != nil {
		return fmt.Errorf("error setting url: %w", err)
	}
	if err := d.Set("enabled", jobDetails.Enabled); err != nil {
		return fmt.Errorf("error setting enabled: %w", err)
	}
	if err := d.Set("save_responses", jobDetails.SaveResponses); err != nil {
		return fmt.Errorf("error setting save_responses: %w", err)
	}
	if err := d.Set("last_status", jobDetails.LastStatus); err != nil {
		return fmt.Errorf("error setting last_status: %w", err)
	}
	if err := d.Set("last_duration", jobDetails.LastDuration); err != nil {
		return fmt.Errorf("error setting last_duration: %w", err)
	}
	if err := d.Set("last_execution", jobDetails.LastExecution); err != nil {
		return fmt.Errorf("error setting last_execution: %w", err)
	}
	if jobDetails.NextExecution != nil {
		if err := d.Set("next_execution", *jobDetails.NextExecution); err != nil {
			return fmt.Errorf("error setting next_execution: %w", err)
		}
	}
	if err := d.Set("type", jobDetails.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("request_timeout", jobDetails.RequestTimeout); err != nil {
		return fmt.Errorf("error setting request_timeout: %w", err)
	}
	if err := d.Set("redirect_success", jobDetails.RedirectSuccess); err != nil {
		return fmt.Errorf("error setting redirect_success: %w", err)
	}
	if err := d.Set("folder_id", jobDetails.FolderID); err != nil {
		return fmt.Errorf("error setting folder_id: %w", err)
	}
	if err := d.Set("request_method", jobDetails.RequestMethod); err != nil {
		return fmt.Errorf("error setting request_method: %w", err)
	}

	// Set schedule
	schedule := []interface{}{
		map[string]interface{}{
			"timezone":   jobDetails.Schedule.Timezone,
			"expires_at": jobDetails.Schedule.ExpiresAt,
			"hours":      jobDetails.Schedule.Hours,
			"mdays":      jobDetails.Schedule.MDays,
			"minutes":    jobDetails.Schedule.Minutes,
			"months":     jobDetails.Schedule.Months,
			"wdays":      jobDetails.Schedule.WDays,
		},
	}
	if err := d.Set("schedule", schedule); err != nil {
		return fmt.Errorf("error setting schedule: %w", err)
	}

	// Set auth
	auth := []interface{}{
		map[string]interface{}{
			"enable":   jobDetails.Auth.Enable,
			"user":     jobDetails.Auth.User,
			"password": jobDetails.Auth.Password,
		},
	}
	if err := d.Set("auth", auth); err != nil {
		return fmt.Errorf("error setting auth: %w", err)
	}

	// Set notification
	notification := []interface{}{
		map[string]interface{}{
			"on_failure": jobDetails.Notification.OnFailure,
			"on_success": jobDetails.Notification.OnSuccess,
			"on_disable": jobDetails.Notification.OnDisable,
		},
	}
	if err := d.Set("notification", notification); err != nil {
		return fmt.Errorf("error setting notification: %w", err)
	}

	// Set extended_data
	extendedData := []interface{}{
		map[string]interface{}{
			"headers": jobDetails.ExtendedData.Headers,
			"body":    jobDetails.ExtendedData.Body,
		},
	}
	if err := d.Set("extended_data", extendedData); err != nil {
		return fmt.Errorf("error setting extended_data: %w", err)
	}

	return nil
}

func resourceJobUpdate(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}

	// Build job object with only changed fields
	job := make(map[string]interface{})

	// Check each field for changes
	if d.HasChange("title") {
		job["title"] = d.Get("title").(string)
	}
	if d.HasChange("url") {
		job["url"] = d.Get("url").(string)
	}
	if d.HasChange("enabled") {
		job["enabled"] = d.Get("enabled").(bool)
	}
	if d.HasChange("save_responses") {
		job["saveResponses"] = d.Get("save_responses").(bool)
	}
	if d.HasChange("request_timeout") {
		job["requestTimeout"] = d.Get("request_timeout").(int)
	}
	if d.HasChange("redirect_success") {
		job["redirectSuccess"] = d.Get("redirect_success").(bool)
	}
	if d.HasChange("folder_id") {
		job["folderId"] = d.Get("folder_id").(int)
	}
	if d.HasChange("request_method") {
		job["requestMethod"] = d.Get("request_method").(int)
	}

	// Schedule
	if d.HasChange("schedule") {
		if scheduleList := d.Get("schedule").([]interface{}); len(scheduleList) > 0 {
			scheduleMap := scheduleList[0].(map[string]interface{})
			schedule := make(map[string]interface{})

			schedule["timezone"] = scheduleMap["timezone"].(string)
			schedule["expiresAt"] = scheduleMap["expires_at"].(int)

			if hours := scheduleMap["hours"].([]interface{}); len(hours) > 0 {
				hoursInt := make([]int, len(hours))
				for i, v := range hours {
					hoursInt[i] = v.(int)
				}
				schedule["hours"] = hoursInt
			} else {
				schedule["hours"] = []int{}
			}

			if mdays := scheduleMap["mdays"].([]interface{}); len(mdays) > 0 {
				mdaysInt := make([]int, len(mdays))
				for i, v := range mdays {
					mdaysInt[i] = v.(int)
				}
				schedule["mdays"] = mdaysInt
			} else {
				schedule["mdays"] = []int{}
			}

			if minutes := scheduleMap["minutes"].([]interface{}); len(minutes) > 0 {
				minutesInt := make([]int, len(minutes))
				for i, v := range minutes {
					minutesInt[i] = v.(int)
				}
				schedule["minutes"] = minutesInt
			} else {
				schedule["minutes"] = []int{}
			}

			if months := scheduleMap["months"].([]interface{}); len(months) > 0 {
				monthsInt := make([]int, len(months))
				for i, v := range months {
					monthsInt[i] = v.(int)
				}
				schedule["months"] = monthsInt
			} else {
				schedule["months"] = []int{}
			}

			if wdays := scheduleMap["wdays"].([]interface{}); len(wdays) > 0 {
				wdaysInt := make([]int, len(wdays))
				for i, v := range wdays {
					wdaysInt[i] = v.(int)
				}
				schedule["wdays"] = wdaysInt
			} else {
				schedule["wdays"] = []int{}
			}

			job["schedule"] = schedule
		}
	}

	// Auth
	if d.HasChange("auth") {
		if authList := d.Get("auth").([]interface{}); len(authList) > 0 {
			authMap := authList[0].(map[string]interface{})
			auth := map[string]interface{}{
				"enable":   authMap["enable"].(bool),
				"user":     authMap["user"].(string),
				"password": authMap["password"].(string),
			}
			job["auth"] = auth
		}
	}

	// Notification
	if d.HasChange("notification") {
		if notificationList := d.Get("notification").([]interface{}); len(notificationList) > 0 {
			notificationMap := notificationList[0].(map[string]interface{})
			notification := map[string]interface{}{
				"onFailure": notificationMap["on_failure"].(bool),
				"onSuccess": notificationMap["on_success"].(bool),
				"onDisable": notificationMap["on_disable"].(bool),
			}
			job["notification"] = notification
		}
	}

	// Extended Data
	if d.HasChange("extended_data") {
		if extendedDataList := d.Get("extended_data").([]interface{}); len(extendedDataList) > 0 {
			extendedDataMap := extendedDataList[0].(map[string]interface{})
			extendedData := make(map[string]interface{})

			if headers := extendedDataMap["headers"].(map[string]interface{}); len(headers) > 0 {
				stringHeaders := make(map[string]string)
				for k, v := range headers {
					stringHeaders[k] = v.(string)
				}
				extendedData["headers"] = stringHeaders
			} else {
				extendedData["headers"] = map[string]string{}
			}

			extendedData["body"] = extendedDataMap["body"].(string)
			job["extendedData"] = extendedData
		}
	}

	// Only update if there are changes
	if len(job) > 0 {
		err := c.UpdateJob(d.Id(), job)
		if err != nil {
			return err
		}
	}

	return resourceJobRead(d, m)
}

func resourceJobDelete(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}

	err := c.DeleteJob(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
