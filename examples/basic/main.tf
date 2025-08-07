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

# Get all jobs from your cron-job.org account
data "cronjob_jobs" "all" {
}

# Example: Get a specific job by ID (uncomment and set actual job ID)
# data "cronjob_job" "example" {
#   job_id = 123
# }

# Example: Get job history (uncomment and set actual job ID)
# data "cronjob_job_history" "example_history" {
#   job_id = 123
# }