package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"cronjob/client"
)

func resourceJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobCreate,
		Read:   resourceJobRead,
		Update: resourceJobUpdate,
		Delete: resourceJobDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Add more fields as needed
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
	var result struct{ JobId int `json:"jobId"` }
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