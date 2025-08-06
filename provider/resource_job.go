// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/plain-insure/terraform-provider-cron-job.org/client"
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
		},
	}
}

func resourceJobCreate(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}
	title, ok := d.Get("title").(string)
	if !ok {
		return fmt.Errorf("title must be a string")
	}
	url, ok := d.Get("url").(string)
	if !ok {
		return fmt.Errorf("url must be a string")
	}
	job := map[string]interface{}{
		"title": title,
		"url":   url,
	}
	reqBody := map[string]interface{}{
		"job": job,
	}
	resp, err := c.DoRequest("PUT", "/jobs", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		JobId int `json:"jobId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", result.JobId))
	return resourceJobRead(d, m)
}

func resourceJobRead(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}
	resp, err := c.DoRequest("GET", fmt.Sprintf("/jobs/%s", d.Id()), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		JobDetails map[string]interface{} `json:"jobDetails"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if err := d.Set("title", result.JobDetails["title"]); err != nil {
		return fmt.Errorf("error setting title: %w", err)
	}
	if err := d.Set("url", result.JobDetails["url"]); err != nil {
		return fmt.Errorf("error setting url: %w", err)
	}
	return nil
}

func resourceJobUpdate(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}
	title, ok := d.Get("title").(string)
	if !ok {
		return fmt.Errorf("title must be a string")
	}
	url, ok := d.Get("url").(string)
	if !ok {
		return fmt.Errorf("url must be a string")
	}
	job := map[string]interface{}{
		"title": title,
		"url":   url,
	}
	reqBody := map[string]interface{}{
		"jobId": d.Id(),
		"job":   job,
	}
	_, err := c.DoRequest("PATCH", fmt.Sprintf("/jobs/%s", d.Id()), reqBody)
	if err != nil {
		return err
	}
	return resourceJobRead(d, m)
}

func resourceJobDelete(d *schema.ResourceData, m interface{}) error {
	c, ok := m.(*client.Client)
	if !ok {
		return fmt.Errorf("expected *client.Client, got %T", m)
	}
	reqBody := map[string]interface{}{
		"jobId": d.Id(),
	}
	_, err := c.DoRequest("DELETE", fmt.Sprintf("/jobs/%s", d.Id()), reqBody)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
