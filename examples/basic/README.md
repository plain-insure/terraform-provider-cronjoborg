# Basic Example

This example demonstrates the basic usage of the cronjob provider.

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

## What this creates

- A folder named "Monitoring Jobs"
- A cron job for health checking
- A status page

## Cleanup

To destroy the created resources:
```bash
terraform destroy
```