terraform {
  required_providers {
    cronjob = {
      source = "registry.terraform.io/plain-insure/cronjob"
    }
  }
}

provider "cronjob" {
  # API key can be set via CRON_JOB_API_KEY environment variable
  # or specified here (not recommended for production)
  # api_key = "your-api-key-here"
}

# Read a single job by ID
data "cronjob_job" "example" {
  job_id = 123 # Replace with actual job ID
}

output "job_details" {
  value = {
    title         = data.cronjob_job.example.title
    url           = data.cronjob_job.example.url
    enabled       = data.cronjob_job.example.enabled
    save_responses = data.cronjob_job.example.save_responses
    schedule      = data.cronjob_job.example.schedule
  }
}