package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/plain-insure/terraform-provider-cron-job.org/client"
)

func resourceStatusPage() *schema.Resource {
	return &schema.Resource{
		Create: resourceStatusPageCreate,
		Read:   resourceStatusPageRead,
		Delete: resourceStatusPageDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true, // Status pages cannot be updated
				Description:  "The title of the status page",
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
		},
	}
}

func resourceStatusPageCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	reqBody := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	resp, err := c.DoRequest("PUT", "/statuspages", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		StatusPageId int `json:"statusPageId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", result.StatusPageId))
	return resourceStatusPageRead(d, m)
}

func resourceStatusPageRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	resp, err := c.DoRequest("GET", fmt.Sprintf("/statuspages/%s", d.Id()), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		StatusPage map[string]interface{} `json:"statusPage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if err := d.Set("title", result.StatusPage["title"]); err != nil {
		return fmt.Errorf("error setting title: %w", err)
	}
	return nil
}

func resourceStatusPageDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	reqBody := map[string]interface{}{
		"statusPageId": d.Id(),
	}
	_, err := c.DoRequest("DELETE", fmt.Sprintf("/statuspages/%s", d.Id()), reqBody)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
