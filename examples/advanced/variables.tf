variable "api_url" {
  description = "The cron-job.org API URL"
  type        = string
  default     = "https://api.cron-job.org/"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}