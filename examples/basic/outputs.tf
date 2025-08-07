output "all_jobs" {
  description = "List of all jobs"
  value       = data.cronjob_jobs.all.jobs
}

output "job_count" {
  description = "Total number of jobs"
  value       = length(data.cronjob_jobs.all.jobs)
}

output "enabled_jobs" {
  description = "List of enabled jobs"
  value = [
    for job in data.cronjob_jobs.all.jobs : job
    if job.enabled
  ]
}

output "job_titles" {
  description = "List of job titles"
  value = [
    for job in data.cronjob_jobs.all.jobs : job.title
  ]
}