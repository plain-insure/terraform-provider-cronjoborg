# Data Sources Examples

This directory contains examples of using the cron-job.org provider's data sources to read job information.

## Overview

The cron-job.org API only supports read operations for jobs, so this provider provides data sources to retrieve job information.

## Available Data Sources

- `cronjoborg_job` - Read a single job by ID
- `cronjoborg_jobs` - Read all jobs
- `cronjoborg_job_history` - Read execution history for a job

## Usage

See the individual `.tf` files for examples of each data source.