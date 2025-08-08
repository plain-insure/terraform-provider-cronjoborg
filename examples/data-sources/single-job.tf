terraform {
  required_providers {
    cronjoborg = {
      source = "registry.terraform.io/plain-insure/cronjoborg"
    }
  }
}

provider "cronjoborg" {
  # API key can be set via CRON_JOB_API_KEY environment variable
  # or specified here (not recommended for production)
  # api_key = "your-api-key-here"
}

# Read a single job by ID
data "cronjoborg_job" "example" {
  job_id = 123 # Replace with actual job ID
}

output "job_details" {
  value = {
    title          = data.cronjoborg_job.example.title
    url            = data.cronjoborg_job.example.url
    enabled        = data.cronjoborg_job.example.enabled
    save_responses = data.cronjoborg_job.example.save_responses
    schedule       = data.cronjoborg_job.example.schedule
  }
}