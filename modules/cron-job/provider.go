package cronjob

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.cron-job.org/",
				Description: "Base URL for the cron-job API. Defaults to https://api.cron-job.org/.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CRON_JOB_API_KEY", nil),
				Description: "API key for the cron-job API. Can also be set via CRON_JOB_API_KEY env variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// Define your resources here (e.g., "cron-job_job": resourceJob()),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type ProviderConfig struct {
	APIUrl string
	APIKey string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiUrl := d.Get("api_url").(string)
	apiKey := d.Get("api_key").(string)
	if apiKey == "" {
		apiKey = os.Getenv("CRON_JOB_API_KEY")
	}

	if apiKey == "" {
		return nil, diag.Errorf("API key must be provided via provider configuration or CRON_JOB_API_KEY environment variable")
	}

	return &ProviderConfig{
		APIUrl: apiUrl,
		APIKey: apiKey,
	}, nil
}