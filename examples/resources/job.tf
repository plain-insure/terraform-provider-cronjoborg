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