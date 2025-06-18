terraform {
  required_providers {
    cronjob = {
      source  = "plain-insure/cronjob"
      version = "~> 1.0"
    }
  }
}

provider "cronjob" {
  api_url = var.api_url
  # api_key set via CRON_JOB_API_KEY environment variable
}

# Create folders for organization
resource "cronjob_folder" "monitoring" {
  title = "Monitoring"
}

resource "cronjob_folder" "backups" {
  title = "Backup Jobs"
}

resource "cronjob_folder" "maintenance" {
  title = "Maintenance Tasks"
}

# Monitoring jobs
resource "cronjob_job" "website_health" {
  title = "Website Health Check"
  url   = "https://mywebsite.com/health"
}

resource "cronjob_job" "api_health" {
  title = "API Health Check"
  url   = "https://api.mywebsite.com/health"
}

resource "cronjob_job" "database_health" {
  title = "Database Health Check"
  url   = "https://monitoring.mywebsite.com/db-health"
}

# Backup jobs
resource "cronjob_job" "daily_backup" {
  title = "Daily Database Backup"
  url   = "https://backup.mywebsite.com/trigger/daily"
}

resource "cronjob_job" "weekly_backup" {
  title = "Weekly Full Backup"
  url   = "https://backup.mywebsite.com/trigger/weekly"
}

# Maintenance jobs
resource "cronjob_job" "cleanup_logs" {
  title = "Cleanup Old Logs"
  url   = "https://maintenance.mywebsite.com/cleanup/logs"
}

resource "cronjob_job" "update_cache" {
  title = "Update Application Cache"
  url   = "https://api.mywebsite.com/cache/refresh"
}

# Status pages for monitoring
resource "cronjob_status_page" "website_status" {
  title = "Website Status"
}

resource "cronjob_status_page" "api_status" {
  title = "API Status"
}

# Dynamic jobs from a list
locals {
  endpoints_to_monitor = [
    {
      name = "User Service"
      url  = "https://user-service.mywebsite.com/health"
    },
    {
      name = "Payment Service"
      url  = "https://payment-service.mywebsite.com/health"
    },
    {
      name = "Notification Service"
      url  = "https://notification-service.mywebsite.com/health"
    }
  ]
}

resource "cronjob_job" "microservice_health" {
  for_each = { for idx, endpoint in local.endpoints_to_monitor : idx => endpoint }

  title = "${each.value.name} Health Check"
  url   = each.value.url
}