# Terraform Provider: cron-job

This provider allows you to manage jobs, folders, and status pages using the cron-job.org API (or a compatible self-hosted service).

## Provider Example

```hcl
provider "cronjob" {
  api_url = "https://api.cron-job.org/"
  api_key = var.cron_job_api_key # or set CRON_JOB_API_KEY env variable
}
```

## Resource Example

```hcl
resource "cronjob_folder" "example" {
  title = "My Folder"
}

resource "cronjob_job" "example" {
  title = "My Job"
  url   = "https://example.com"
}

resource "cronjob_status_page" "example" {
  title = "My Status Page"
}
```