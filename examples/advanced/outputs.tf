output "monitoring_folder_id" {
  description = "ID of the monitoring folder"
  value       = cronjob_folder.monitoring.id
}

output "all_job_ids" {
  description = "IDs of all created jobs"
  value = merge(
    {
      "website_health"  = cronjob_job.website_health.id
      "api_health"      = cronjob_job.api_health.id
      "database_health" = cronjob_job.database_health.id
      "daily_backup"    = cronjob_job.daily_backup.id
      "weekly_backup"   = cronjob_job.weekly_backup.id
      "cleanup_logs"    = cronjob_job.cleanup_logs.id
      "update_cache"    = cronjob_job.update_cache.id
    },
    { for idx, job in cronjob_job.microservice_health : "microservice_${idx}" => job.id }
  )
}

output "status_page_ids" {
  description = "IDs of created status pages"
  value = {
    website = cronjob_status_page.website_status.id
    api     = cronjob_status_page.api_status.id
  }
}

output "microservice_jobs" {
  description = "Details of microservice monitoring jobs"
  value = {
    for idx, job in cronjob_job.microservice_health : idx => {
      id    = job.id
      title = job.title
      url   = job.url
    }
  }
}