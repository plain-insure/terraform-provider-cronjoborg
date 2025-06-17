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
					value := v.(string)
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
	c := m.(*client.Client)
	job := map[string]interface{}{
		"title": d.Get("title").(string),
		"url":   d.Get("url").(string),
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
	c := m.(*client.Client)
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
	d.Set("title", result.JobDetails["title"])
	d.Set("url", result.JobDetails["url"])
	return nil
}

func resourceJobUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	job := map[string]interface{}{
		"title": d.Get("title").(string),
		"url":   d.Get("url").(string),
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
	c := m.(*client.Client)
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
