# Advanced Example

This example demonstrates advanced usage patterns of the cronjoborg provider including:

- Multiple folders for organization
- Various types of monitoring and maintenance jobs
- Status pages for different services
- Dynamic job creation using `for_each`
- Local values for reusable configurations

## Features Demonstrated

### Folder Organization
- Monitoring folder for health checks
- Backup folder for backup jobs
- Maintenance folder for routine tasks

### Job Types
- **Health Checks**: Website, API, and database monitoring
- **Backups**: Daily and weekly backup triggers
- **Maintenance**: Log cleanup and cache updates
- **Microservices**: Dynamic monitoring of multiple services

### Advanced Patterns
- Using `locals` for reusable endpoint configurations
- Dynamic resource creation with `for_each`
- Organized outputs for easy reference

## Usage

1. Set your API key:
   ```bash
   export CRON_JOB_API_KEY="your-api-key-here"
   ```

2. Customize the variables:
   ```bash
   # Optional: customize the environment
   export TF_VAR_environment="prod"
   
   # Optional: use a different API URL (for self-hosted instances)
   export TF_VAR_api_url="https://your-cron-job-api.com/"
   ```

3. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Customization

### Adding More Microservices

Edit the `locals.endpoints_to_monitor` list in `main.tf`:

```hcl
locals {
  endpoints_to_monitor = [
    {
      name = "User Service"
      url  = "https://user-service.mywebsite.com/health"
    },
    {
      name = "Your New Service"
      url  = "https://your-new-service.mywebsite.com/health"
    }
    # Add more services here
  ]
}
```

### Environment-Specific Configurations

You can use the `environment` variable to customize behavior:

```hcl
resource "cronjoborg_job" "health_check" {
  title = "${var.environment}-website-health"
  url   = "https://${var.environment}.mywebsite.com/health"
}
```

## Cleanup

```bash
terraform destroy
```