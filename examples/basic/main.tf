terraform {
  required_providers {
    cronjob = {
      source  = "plain-insure/cronjob"
      version = "~> 1.0"
    }
  }
}

provider "cronjob" {
  api_url = "https://api.cron-job.org/"
  # api_key can be set via CRON_JOB_API_KEY environment variable
  # api_key = var.cron_job_api_key
}

# Create a new cron job
resource "cronjob_job" "example" {
  title = "Example Terraform Job"
  url   = "https://httpbin.org/post"
}

# Get all jobs from your cron-job.org account
data "cronjob_jobs" "all" {
}

# Get details of the job we just created
data "cronjob_job" "created_job" {
  job_id = cronjob_job.example.id
}

# Get job history for the created job
data "cronjob_job_history" "example_history" {
  job_id = cronjob_job.example.id
}