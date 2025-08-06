output "folder_id" {
  description = "ID of the created folder"
  value       = cronjob_folder.monitoring.id
}

output "job_id" {
  description = "ID of the created job"
  value       = cronjob_job.health_check.id
}

output "status_page_id" {
  description = "ID of the created status page"
  value       = cronjob_status_page.example.id
}