# Terraform Module: cron-job

This module provides a Terraform provider for working with the [cron-job.org](https://www.cron-job.org/) API and self-hosted compatible installations.

## Provider Configuration

```hcl
provider "cron-job" {
  api_url = "https://api.cron-job.org/" # Optional. Use your own URL for self-hosted.
  api_key = var.cron_job_api_key        # Or set via environment variable CRON_JOB_API_KEY
}