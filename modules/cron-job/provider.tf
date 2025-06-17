variable "api_url" {
  description = "Base URL for the cron-job API. Defaults to https://api.cron-job.org/"
  type        = string
  default     = "https://api.cron-job.org/"
}

variable "api_key" {
  description = "API key for the cron-job API. Can also be set via the CRON_JOB_API_KEY environment variable."
  type        = string
  default     = ""
  sensitive   = true
}

provider "cron-job" {
  api_url = var.api_url
  api_key = var.api_key != "" ? var.api_key : (try(env("CRON_JOB_API_KEY"), ""))
}