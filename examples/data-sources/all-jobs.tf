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

# Read all jobs
data "cronjoborg_jobs" "all" {
}

output "all_jobs" {
  value = data.cronjoborg_jobs.all.jobs
}

output "job_count" {
  value = length(data.cronjoborg_jobs.all.jobs)
}

# Example of filtering enabled jobs
locals {
  enabled_jobs = [
    for job in data.cronjoborg_jobs.all.jobs : job
    if job.enabled
  ]
}

output "enabled_jobs" {
  value = local.enabled_jobs
}

output "enabled_job_count" {
  value = length(local.enabled_jobs)
}