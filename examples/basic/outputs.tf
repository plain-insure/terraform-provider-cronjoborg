output "created_job_id" {
  description = "ID of the created job"
  value       = cronjoborg_job.example.id
}

output "all_jobs" {
  description = "List of all jobs"
  value       = data.cronjoborg_jobs.all.jobs
}

output "job_count" {
  description = "Total number of jobs"
  value       = length(data.cronjoborg_jobs.all.jobs)
}

output "enabled_jobs" {
  description = "List of enabled jobs"
  value = [
    for job in data.cronjoborg_jobs.all.jobs : job
    if job.enabled
  ]
}

output "job_titles" {
  description = "List of job titles"
  value = [
    for job in data.cronjoborg_jobs.all.jobs : job.title
  ]
}

output "created_job_details" {
  description = "Details of the job we created"
  value       = data.cronjoborg_job.created_job
}

output "created_job_history" {
  description = "Execution history of the created job"
  value       = data.cronjoborg_job_history.example_history.history
}