package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/plain-insure/terraform-provider-cron-job.org/client"
)

func resourceFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceFolderCreate,
		Read:   resourceFolderRead,
		Delete: resourceFolderDelete,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true, // Folders cannot be updated
				Description:  "The title of the folder",
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
		},
	}
}

func resourceFolderCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	reqBody := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	resp, err := c.DoRequest("PUT", "/folders", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		FolderId int `json:"folderId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", result.FolderId))
	return resourceFolderRead(d, m)
}

func resourceFolderRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	resp, err := c.DoRequest("GET", "/folders", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		Folders []map[string]interface{} `json:"folders"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	// Find the folder by ID
	for _, folder := range result.Folders {
		id := fmt.Sprintf("%v", folder["folderId"])
		if id == d.Id() {
			d.Set("title", folder["title"])
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceFolderDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	reqBody := map[string]interface{}{
		"folderId": d.Id(),
	}
	_, err := c.DoRequest("DELETE", fmt.Sprintf("/folders/%s", d.Id()), reqBody)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
