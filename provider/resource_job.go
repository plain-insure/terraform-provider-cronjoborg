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
	title, ok := d.Get("title").(string)
	if !ok {
		return fmt.Errorf("title must be a string")
	}
	job["title"] = title

	url, ok := d.Get("url").(string)
	if !ok {
		return fmt.Errorf("url must be a string")
	}
	job["url"] = url

	// Optional fields with defaults
	enabled, ok := d.Get("enabled").(bool)
	if !ok {
		return fmt.Errorf("enabled must be a boolean")
	}
	job["enabled"] = enabled

	saveResponses, ok := d.Get("save_responses").(bool)
	if !ok {
		return fmt.Errorf("save_responses must be a boolean")
	}
	job["saveResponses"] = saveResponses

	requestTimeout, ok := d.Get("request_timeout").(int)
	if !ok {
		return fmt.Errorf("request_timeout must be an integer")
	}
	job["requestTimeout"] = requestTimeout

	redirectSuccess, ok := d.Get("redirect_success").(bool)
	if !ok {
		return fmt.Errorf("redirect_success must be a boolean")
	}
	job["redirectSuccess"] = redirectSuccess

	folderId, ok := d.Get("folder_id").(int)
	if !ok {
		return fmt.Errorf("folder_id must be an integer")
	}
	job["folderId"] = folderId

	requestMethod, ok := d.Get("request_method").(int)
	if !ok {
		return fmt.Errorf("request_method must be an integer")
	}
	job["requestMethod"] = requestMethod

	// Schedule - always include schedule with default values
	schedule, err := buildScheduleFromResourceData(d)
	if err != nil {
		return err
	}
	job["schedule"] = schedule

	// Auth
	if authListRaw := d.Get("auth"); authListRaw != nil {
		authList, ok := authListRaw.([]interface{})
		if !ok {
			return fmt.Errorf("auth must be a list")
		}
		if len(authList) > 0 {
			authMap, ok := authList[0].(map[string]interface{})
			if !ok {
				return fmt.Errorf("auth element must be a map")
			}
			enable, ok := authMap["enable"].(bool)
			if !ok {
				return fmt.Errorf("auth.enable must be a boolean")
			}
			user, ok := authMap["user"].(string)
			if !ok {
				return fmt.Errorf("auth.user must be a string")
			}
			password, ok := authMap["password"].(string)
			if !ok {
				return fmt.Errorf("auth.password must be a string")
			}
			auth := map[string]interface{}{
				"enable":   enable,
				"user":     user,
				"password": password,
			}
			job["auth"] = auth
		}
	}

	// Notification
	if notificationListRaw := d.Get("notification"); notificationListRaw != nil {
		notificationList, ok := notificationListRaw.([]interface{})
		if !ok {
			return fmt.Errorf("notification must be a list")
		}
		if len(notificationList) > 0 {
			notificationMap, ok := notificationList[0].(map[string]interface{})
			if !ok {
				return fmt.Errorf("notification element must be a map")
			}
			onFailure, ok := notificationMap["on_failure"].(bool)
			if !ok {
				return fmt.Errorf("notification.on_failure must be a boolean")
			}
			onSuccess, ok := notificationMap["on_success"].(bool)
			if !ok {
				return fmt.Errorf("notification.on_success must be a boolean")
			}
			onDisable, ok := notificationMap["on_disable"].(bool)
			if !ok {
				return fmt.Errorf("notification.on_disable must be a boolean")
			}
			notification := map[string]interface{}{
				"onFailure": onFailure,
				"onSuccess": onSuccess,
				"onDisable": onDisable,
			}
			job["notification"] = notification
		}
	}

	// Extended Data
	if extendedDataListRaw := d.Get("extended_data"); extendedDataListRaw != nil {
		extendedDataList, ok := extendedDataListRaw.([]interface{})
		if !ok {
			return fmt.Errorf("extended_data must be a list")
		}
		if len(extendedDataList) > 0 {
			extendedDataMap, ok := extendedDataList[0].(map[string]interface{})
			if !ok {
				return fmt.Errorf("extended_data element must be a map")
			}
			extendedData := make(map[string]interface{})

			if headersRaw, exists := extendedDataMap["headers"]; exists {
				headers, ok := headersRaw.(map[string]interface{})
				if !ok {
					return fmt.Errorf("extended_data.headers must be a map")
				}
				if len(headers) > 0 {
					stringHeaders := make(map[string]string)
					for k, v := range headers {
						vStr, ok := v.(string)
						if !ok {
							return fmt.Errorf("extended_data.headers values must be strings")
						}
						stringHeaders[k] = vStr
					}
					extendedData["headers"] = stringHeaders
				} else {
					extendedData["headers"] = map[string]string{}
				}
			} else {
				extendedData["headers"] = map[string]string{}
			}

			if bodyRaw, exists := extendedDataMap["body"]; exists {
				body, ok := bodyRaw.(string)
				if !ok {
					return fmt.Errorf("extended_data.body must be a string")
				}
				extendedData["body"] = body
			}
			job["extendedData"] = extendedData
		}
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
		title, ok := d.Get("title").(string)
		if !ok {
			return fmt.Errorf("title must be a string")
		}
		job["title"] = title
	}
	if d.HasChange("url") {
		url, ok := d.Get("url").(string)
		if !ok {
			return fmt.Errorf("url must be a string")
		}
		job["url"] = url
	}
	if d.HasChange("enabled") {
		enabled, ok := d.Get("enabled").(bool)
		if !ok {
			return fmt.Errorf("enabled must be a boolean")
		}
		job["enabled"] = enabled
	}
	if d.HasChange("save_responses") {
		saveResponses, ok := d.Get("save_responses").(bool)
		if !ok {
			return fmt.Errorf("save_responses must be a boolean")
		}
		job["saveResponses"] = saveResponses
	}
	if d.HasChange("request_timeout") {
		requestTimeout, ok := d.Get("request_timeout").(int)
		if !ok {
			return fmt.Errorf("request_timeout must be an integer")
		}
		job["requestTimeout"] = requestTimeout
	}
	if d.HasChange("redirect_success") {
		redirectSuccess, ok := d.Get("redirect_success").(bool)
		if !ok {
			return fmt.Errorf("redirect_success must be a boolean")
		}
		job["redirectSuccess"] = redirectSuccess
	}
	if d.HasChange("folder_id") {
		folderId, ok := d.Get("folder_id").(int)
		if !ok {
			return fmt.Errorf("folder_id must be an integer")
		}
		job["folderId"] = folderId
	}
	if d.HasChange("request_method") {
		requestMethod, ok := d.Get("request_method").(int)
		if !ok {
			return fmt.Errorf("request_method must be an integer")
		}
		job["requestMethod"] = requestMethod
	}

	// Schedule
	if d.HasChange("schedule") {
		schedule, err := buildScheduleFromResourceData(d)
		if err != nil {
			return err
		}
		job["schedule"] = schedule
	}

	// Auth
	if d.HasChange("auth") {
		if authListRaw := d.Get("auth"); authListRaw != nil {
			authList, ok := authListRaw.([]interface{})
			if !ok {
				return fmt.Errorf("auth must be a list")
			}
			if len(authList) > 0 {
				authMap, ok := authList[0].(map[string]interface{})
				if !ok {
					return fmt.Errorf("auth element must be a map")
				}
				enable, ok := authMap["enable"].(bool)
				if !ok {
					return fmt.Errorf("auth.enable must be a boolean")
				}
				user, ok := authMap["user"].(string)
				if !ok {
					return fmt.Errorf("auth.user must be a string")
				}
				password, ok := authMap["password"].(string)
				if !ok {
					return fmt.Errorf("auth.password must be a string")
				}
				auth := map[string]interface{}{
					"enable":   enable,
					"user":     user,
					"password": password,
				}
				job["auth"] = auth
			}
		}
	}

	// Notification
	if d.HasChange("notification") {
		if notificationListRaw := d.Get("notification"); notificationListRaw != nil {
			notificationList, ok := notificationListRaw.([]interface{})
			if !ok {
				return fmt.Errorf("notification must be a list")
			}
			if len(notificationList) > 0 {
				notificationMap, ok := notificationList[0].(map[string]interface{})
				if !ok {
					return fmt.Errorf("notification element must be a map")
				}
				onFailure, ok := notificationMap["on_failure"].(bool)
				if !ok {
					return fmt.Errorf("notification.on_failure must be a boolean")
				}
				onSuccess, ok := notificationMap["on_success"].(bool)
				if !ok {
					return fmt.Errorf("notification.on_success must be a boolean")
				}
				onDisable, ok := notificationMap["on_disable"].(bool)
				if !ok {
					return fmt.Errorf("notification.on_disable must be a boolean")
				}
				notification := map[string]interface{}{
					"onFailure": onFailure,
					"onSuccess": onSuccess,
					"onDisable": onDisable,
				}
				job["notification"] = notification
			}
		}
	}

	// Extended Data
	if d.HasChange("extended_data") {
		if extendedDataListRaw := d.Get("extended_data"); extendedDataListRaw != nil {
			extendedDataList, ok := extendedDataListRaw.([]interface{})
			if !ok {
				return fmt.Errorf("extended_data must be a list")
			}
			if len(extendedDataList) > 0 {
				extendedDataMap, ok := extendedDataList[0].(map[string]interface{})
				if !ok {
					return fmt.Errorf("extended_data element must be a map")
				}
				extendedData := make(map[string]interface{})

				if headersRaw, exists := extendedDataMap["headers"]; exists {
					headers, ok := headersRaw.(map[string]interface{})
					if !ok {
						return fmt.Errorf("extended_data.headers must be a map")
					}
					if len(headers) > 0 {
						stringHeaders := make(map[string]string)
						for k, v := range headers {
							vStr, ok := v.(string)
							if !ok {
								return fmt.Errorf("extended_data.headers values must be strings")
							}
							stringHeaders[k] = vStr
						}
						extendedData["headers"] = stringHeaders
					} else {
						extendedData["headers"] = map[string]string{}
					}
				} else {
					extendedData["headers"] = map[string]string{}
				}

				if bodyRaw, exists := extendedDataMap["body"]; exists {
					body, ok := bodyRaw.(string)
					if !ok {
						return fmt.Errorf("extended_data.body must be a string")
					}
					extendedData["body"] = body
				}
				job["extendedData"] = extendedData
			}
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

// buildScheduleFromResourceData extracts schedule configuration from resource data and applies defaults.
func buildScheduleFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	schedule := make(map[string]interface{})

	scheduleListRaw := d.Get("schedule")
	if scheduleListRaw != nil {
		scheduleList, ok := scheduleListRaw.([]interface{})
		if !ok {
			return nil, fmt.Errorf("schedule must be a list")
		}
		if len(scheduleList) > 0 {
			scheduleMap, ok := scheduleList[0].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("schedule element must be a map")
			}

			timezone, ok := scheduleMap["timezone"].(string)
			if !ok {
				return nil, fmt.Errorf("schedule.timezone must be a string")
			}
			schedule["timezone"] = timezone

			expiresAt, ok := scheduleMap["expires_at"].(int)
			if !ok {
				return nil, fmt.Errorf("schedule.expires_at must be an integer")
			}
			schedule["expiresAt"] = expiresAt

			if hoursRaw, exists := scheduleMap["hours"]; exists {
				hours, ok := hoursRaw.([]interface{})
				if !ok {
					return nil, fmt.Errorf("schedule.hours must be a list")
				}
				if len(hours) > 0 {
					hoursInt := make([]int, len(hours))
					for i, v := range hours {
						vInt, ok := v.(int)
						if !ok {
							return nil, fmt.Errorf("schedule.hours values must be integers")
						}
						hoursInt[i] = vInt
					}
					schedule["hours"] = hoursInt
				} else {
					schedule["hours"] = []int{-1}
				}
			} else {
				schedule["hours"] = []int{-1}
			}

			if mdaysRaw, exists := scheduleMap["mdays"]; exists {
				mdays, ok := mdaysRaw.([]interface{})
				if !ok {
					return nil, fmt.Errorf("schedule.mdays must be a list")
				}
				if len(mdays) > 0 {
					mdaysInt := make([]int, len(mdays))
					for i, v := range mdays {
						vInt, ok := v.(int)
						if !ok {
							return nil, fmt.Errorf("schedule.mdays values must be integers")
						}
						mdaysInt[i] = vInt
					}
					schedule["mdays"] = mdaysInt
				} else {
					schedule["mdays"] = []int{-1}
				}
			} else {
				schedule["mdays"] = []int{-1}
			}

			if minutesRaw, exists := scheduleMap["minutes"]; exists {
				minutes, ok := minutesRaw.([]interface{})
				if !ok {
					return nil, fmt.Errorf("schedule.minutes must be a list")
				}
				if len(minutes) > 0 {
					minutesInt := make([]int, len(minutes))
					for i, v := range minutes {
						vInt, ok := v.(int)
						if !ok {
							return nil, fmt.Errorf("schedule.minutes values must be integers")
						}
						minutesInt[i] = vInt
					}
					schedule["minutes"] = minutesInt
				} else {
					schedule["minutes"] = []int{-1}
				}
			} else {
				schedule["minutes"] = []int{-1}
			}

			if monthsRaw, exists := scheduleMap["months"]; exists {
				months, ok := monthsRaw.([]interface{})
				if !ok {
					return nil, fmt.Errorf("schedule.months must be a list")
				}
				if len(months) > 0 {
					monthsInt := make([]int, len(months))
					for i, v := range months {
						vInt, ok := v.(int)
						if !ok {
							return nil, fmt.Errorf("schedule.months values must be integers")
						}
						monthsInt[i] = vInt
					}
					schedule["months"] = monthsInt
				} else {
					schedule["months"] = []int{-1}
				}
			} else {
				schedule["months"] = []int{-1}
			}

			if wdaysRaw, exists := scheduleMap["wdays"]; exists {
				wdays, ok := wdaysRaw.([]interface{})
				if !ok {
					return nil, fmt.Errorf("schedule.wdays must be a list")
				}
				if len(wdays) > 0 {
					wdaysInt := make([]int, len(wdays))
					for i, v := range wdays {
						vInt, ok := v.(int)
						if !ok {
							return nil, fmt.Errorf("schedule.wdays values must be integers")
						}
						wdaysInt[i] = vInt
					}
					schedule["wdays"] = wdaysInt
				} else {
					schedule["wdays"] = []int{-1}
				}
			} else {
				schedule["wdays"] = []int{-1}
			}
		} else {
			// Empty schedule list - use defaults
			schedule["timezone"] = "UTC"
			schedule["expiresAt"] = 0
			schedule["hours"] = []int{-1}
			schedule["mdays"] = []int{-1}
			schedule["minutes"] = []int{-1}
			schedule["months"] = []int{-1}
			schedule["wdays"] = []int{-1}
		}
	} else {
		// No schedule block provided - use defaults
		schedule["timezone"] = "UTC"
		schedule["expiresAt"] = 0
		schedule["hours"] = []int{-1}
		schedule["mdays"] = []int{-1}
		schedule["minutes"] = []int{-1}
		schedule["months"] = []int{-1}
		schedule["wdays"] = []int{-1}
	}

	return schedule, nil
}
