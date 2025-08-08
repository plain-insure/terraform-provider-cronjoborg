# Terraform Provider for cron-job.org

[![Tests](https://github.com/plain-insure/terraform-provider-cronjoborg/workflows/Tests/badge.svg)](https://github.com/plain-insure/terraform-provider-cronjoborg/actions?query=workflow%3ATests)
[![Release](https://github.com/plain-insure/terraform-provider-cronjoborg/workflows/Release/badge.svg)](https://github.com/plain-insure/terraform-provider-cronjoborg/actions?query=workflow%3ARelease)

This provider allows you to manage jobs, folders, and status pages using the [cron-job.org](https://www.cron-job.org/) API or a compatible self-hosted service.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20 (for building from source)

## Installation

### From the Terraform Registry

```hcl
terraform {
  required_providers {
    cronjoborg = {
      source  = "plain-insure/cronjoborg"
      version = "~> 1.0"
    }
  }
}
```

### From Source

```bash
git clone https://github.com/plain-insure/terraform-provider-cronjoborg.git
cd terraform-provider-cronjoborg
make build
make install
```

## Usage

### Provider Configuration

```hcl
provider "cronjoborg" {
  api_url = "https://api.cron-job.org/"  # Optional, defaults to cron-job.org API
  api_key = var.cron_job_api_key         # Required, or set CRON_JOB_API_KEY env var
}
```

### Basic Example

```hcl
# Create a folder
resource "cronjoborg_folder" "monitoring" {
  title = "Monitoring Jobs"
}

# Create a cron job
resource "cronjoborg_job" "health_check" {
  title = "Health Check"
  url   = "https://example.com/health"
}

# Create a status page
resource "cronjoborg_status_page" "uptime" {
  title = "Service Uptime"
}
```

For more examples, see the [examples/](./examples/) directory.

## Authentication

The provider requires an API key from cron-job.org. You can obtain one by:

1. Creating an account at [cron-job.org](https://www.cron-job.org/)
2. Generating an API key in your account settings

Set the API key using one of these methods:

### Environment Variable (Recommended)
```bash
export CRON_JOB_API_KEY="your-api-key-here"
```

### Provider Configuration
```hcl
provider "cronjoborg" {
  api_key = "your-api-key-here"
}
```

### Variable
```hcl
variable "cron_job_api_key" {
  description = "API key for cron-job.org"
  type        = string
  sensitive   = true
}

provider "cronjoborg" {
  api_key = var.cron_job_api_key
}
```

## Resources

- `cronjoborg_job` - Manages cron jobs
- `cronjoborg_folder` - Manages folders for organizing jobs
- `cronjoborg_status_page` - Manages status pages

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/plain-insure/cronjob/latest/docs).

## Development

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup and contribution guidelines.

### Quick Start for Development

```bash
# Clone the repository
git clone https://github.com/plain-insure/terraform-provider-cronjoborg.git
cd terraform-provider-cronjoborg

# Install development tools
make dev-setup

# Build the provider
make build

# Run tests
make test

# Install locally for testing
make install
```

## Available Make Targets

Run `make help` to see all available targets:

```bash
make help
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 For bug reports and feature requests, please use [GitHub Issues](https://github.com/plain-insure/terraform-provider-cronjoborg/issues)
- 💬 For questions and discussions, please use [GitHub Discussions](https://github.com/plain-insure/terraform-provider-cronjoborg/discussions)
- 📖 Check out the [examples](./examples/) for common use cases