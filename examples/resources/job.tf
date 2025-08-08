terraform {
  required_providers {
    cronjoborg = {
      source = "plain-insure/cronjoborg"
    }
  }
}

provider "cronjoborg" {
  # API key can be set via CRON_JOB_API_KEY environment variable
  # api_key = "your-api-key-here"
}

# Create a cron job
resource "cronjoborg_job" "example" {
  title = "Example Job"
  url   = "https://example.com/webhook"
}

# Output the job ID
output "job_id" {
  value = cronjoborg_job.example.id
}