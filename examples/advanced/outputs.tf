output "monitoring_folder_id" {
  description = "ID of the monitoring folder"
  value       = cronjoborg_folder.monitoring.id
}

output "all_job_ids" {
  description = "IDs of all created jobs"
  value = merge(
    {
      "website_health"  = cronjoborg_job.website_health.id
      "api_health"      = cronjoborg_job.api_health.id
      "database_health" = cronjoborg_job.database_health.id
      "daily_backup"    = cronjoborg_job.daily_backup.id
      "weekly_backup"   = cronjoborg_job.weekly_backup.id
      "cleanup_logs"    = cronjoborg_job.cleanup_logs.id
      "update_cache"    = cronjoborg_job.update_cache.id
    },
    { for idx, job in cronjoborg_job.microservice_health : "microservice_${idx}" => job.id }
  )
}

output "status_page_ids" {
  description = "IDs of created status pages"
  value = {
    website = cronjoborg_status_page.website_status.id
    api     = cronjoborg_status_page.api_status.id
  }
}

output "microservice_jobs" {
  description = "Details of microservice monitoring jobs"
  value = {
    for idx, job in cronjoborg_job.microservice_health : idx => {
      id    = job.id
      title = job.title
      url   = job.url
    }
  }
}