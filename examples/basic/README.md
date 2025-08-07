# Basic Example

This example demonstrates the basic usage of the cronjob provider data sources.

## Overview

The cron-job.org API only supports read operations, so this provider offers data sources to retrieve job information rather than managing jobs directly.

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

- Retrieves all jobs from your cron-job.org account
- Shows individual job details if you have any jobs
- Outputs job count and summary information

## Note

This provider only provides data sources since the cron-job.org API only supports listing jobs and job history. It does not support creating, updating, or deleting jobs via the API.