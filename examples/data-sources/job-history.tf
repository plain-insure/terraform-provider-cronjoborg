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

# Read job execution history
data "cronjob_job_history" "example" {
  job_id = 123 # Replace with actual job ID
}

output "job_history" {
  value = data.cronjob_job_history.example.history
}

output "last_execution" {
  value = length(data.cronjob_job_history.example.history) > 0 ? data.cronjob_job_history.example.history[0] : null
}

# Example of filtering successful executions
locals {
  successful_executions = [
    for execution in data.cronjob_job_history.example.history : execution
    if execution.status == "OK"
  ]
}

output "successful_executions" {
  value = local.successful_executions
}

output "success_rate" {
  value = length(data.cronjob_job_history.example.history) > 0 ? length(local.successful_executions) / length(data.cronjob_job_history.example.history) * 100 : 0
}