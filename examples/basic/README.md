# Basic Example

This example demonstrates the basic usage of the cronjoborg provider, including both resource management and data sources.

## Overview

The cronjoborg provider supports full CRUD operations for cron jobs, including:
- Creating new jobs with `cronjoborg_job` resource
- Reading job details with `cronjoborg_job` data source
- Listing all jobs with `cronjoborg_jobs` data source
- Viewing job execution history with `cronjoborg_job_history` data source

## Usage

1. Set your API key:
   ```bash
   export CRON_JOB_API_KEY="your-api-key-here"
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Plan the deployment:
   ```bash
   terraform plan
   ```

4. Apply the configuration:
   ```bash
   terraform apply
   ```

## What this does

- Creates a new cron job with the title "Example Terraform Job"
- Retrieves all jobs from your cron-job.org account
- Shows details of the job we just created
- Retrieves execution history for the created job
- Outputs job count and summary information

## Resources and Data Sources

This example uses:
- `cronjoborg_job` resource: Creates and manages a cron job
- `cronjoborg_jobs` data source: Lists all jobs in your account
- `cronjoborg_job` data source: Gets details of a specific job
- `cronjoborg_job_history` data source: Gets execution history for a job