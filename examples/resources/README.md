# Job Resource Example

This example demonstrates how to create and manage cron jobs using the `cronjob_job` resource.

## Configuration

```terraform
terraform {
  required_providers {
    cronjob = {
      source = "plain-insure/cronjob"
    }
  }
}

provider "cronjob" {
  # API key can be set via CRON_JOB_API_KEY environment variable
  # api_key = "your-api-key-here"
}

# Create a cron job
resource "cronjob_job" "example" {
  title = "Example Job"
  url   = "https://example.com/webhook"
}

# Output the job ID
output "job_id" {
  value = cronjob_job.example.id
}
```

## Usage

1. Set your API key:
   ```bash
   export CRON_JOB_API_KEY="your-api-key"
   ```

2. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

This will create a new cron job with the specified title and URL.