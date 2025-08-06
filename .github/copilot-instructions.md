# Terraform Provider for cron-job.org

**ALWAYS follow these instructions first and only fallback to additional search and context gathering if the information here is incomplete or found to be in error.**

This is a Terraform provider for cron-job.org written in Go using the Terraform Plugin SDK v2. The provider allows management of cron jobs, folders, and status pages using the cron-job.org API.

## Working Effectively

### Initial Setup
Install development tools and dependencies:
```bash
# CRITICAL: Add Go tools to PATH - required for all linting/docs commands
export PATH=$PATH:~/go/bin

# Install development tools - NEVER CANCEL: Takes ~3-4 minutes
make dev-setup

# Download dependencies - Takes ~30 seconds
make deps
```

### Core Development Commands
Build and test the provider:
```bash
# Build the provider binary - Takes ~10 seconds
make build

# Run unit tests - Takes ~7 seconds, NEVER CANCEL
make test

# Run tests with coverage - Takes ~2 seconds
make test-cover

# Format Go code - Takes <1 second
make fmt

# Run go vet - Takes ~4 seconds
make vet

# Run linters - Takes ~8 seconds, REQUIRES PATH=$PATH:~/go/bin
make lint

# Generate documentation - Takes ~1 second, REQUIRES terraform installed
make docs

# Install provider locally for testing - Takes <1 second
make install

# Clean build artifacts
make clean
```

### Prerequisites and Installation Requirements
- **Go 1.20+** (tested with 1.24.5)
- **Terraform 1.0+** for documentation generation
- **golangci-lint** and **tfplugindocs** (installed via `make dev-setup`)
- **CRITICAL PATH REQUIREMENT**: Export `PATH=$PATH:~/go/bin` before running lint/docs commands

### Install Terraform (Required for Documentation Generation)
```bash
cd /tmp
wget https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip
unzip terraform_1.5.7_linux_amd64.zip
sudo mv terraform /usr/local/bin/
terraform version
```

### Local Provider Testing
Install the provider locally for testing with Terraform:
```bash
# Build and install to local terraform plugins directory
make install

# Provider will be available as:
# registry.terraform.io/plain-insure/cronjob version "dev"
```

## Validation and Testing

### Before Committing Changes
**ALWAYS run these commands before committing to ensure CI passes:**
```bash
# CRITICAL: Set PATH first
export PATH=$PATH:~/go/bin

# Format, lint, and test - NEVER CANCEL any of these
make fmt
make vet  
make lint
make test
make docs
```

### Manual Validation Requirements
- **ALWAYS** run through the complete command sequence above after making changes
- Verify that `make build` produces a working binary in `bin/terraform-provider-cronjob`
- Confirm that `make install` completes without errors
- Test that documentation generation works with `make docs`
- **NOTE**: Acceptance tests require a valid CRON_JOB_API_KEY environment variable but are currently disabled in CI

### Unit Test Coverage
- Client package: ~83% coverage
- Provider package: ~9% coverage (mostly unit tests, acceptance tests require API key)
- Tests run quickly (~6-7 seconds total)

## Critical Timing and Timeout Information

| Command | Time | Timeout Needed | Notes |
|---------|------|----------------|--------|
| `make dev-setup` | ~3-4 minutes | 10+ minutes | NEVER CANCEL - Downloads many dependencies |
| `make deps` | ~30 seconds | 5 minutes | Downloads Go modules |
| `make build` | ~10 seconds | 2 minutes | Very fast compilation |
| `make test` | ~7 seconds | 2 minutes | Unit tests only |
| `make lint` | ~8 seconds | 2 minutes | Requires PATH setup |
| `make docs` | ~1 second | 1 minute | Requires terraform |

## Repository Structure

### Key Directories
- `provider/` - Provider implementation and resource definitions
  - `provider.go` - Main provider configuration
  - `resource_*.go` - Resource implementations (job, folder, status_page)
  - `*_test.go` - Unit and acceptance tests
- `client/` - API client for cron-job.org service
- `examples/` - Terraform configuration examples
  - `basic/` - Simple examples
  - `advanced/` - Complex configurations
- `docs/` - Generated documentation (auto-generated via `make docs`)
- `templates/` - Documentation templates for tfplugindocs

### Key Files
- `main.go` - Provider entry point with go:generate directives
- `Makefile` - All build targets and commands
- `go.mod` - Go module definition (Go 1.21)
- `.github/workflows/test.yml` - CI pipeline with unit tests and linting
- `.github/workflows/release.yml` - Release pipeline using goreleaser
- `.golangci.yml` - Linter configuration
- `.goreleaser.yml` - Release configuration

## Common Issues and Solutions

### PATH Issues
If `make lint` or `make docs` fails with "command not found":
```bash
# SOLUTION: Add Go tools to PATH
export PATH=$PATH:~/go/bin
# Then retry the command
```

### Documentation Generation Fails
If `make docs` fails with "terraform: executable file not found":
```bash
# SOLUTION: Install Terraform (see installation commands above)
```

### Linting Failures
If golangci-lint reports issues:
```bash
# Fix formatting first
make fmt
# Then re-run linting
export PATH=$PATH:~/go/bin && make lint
```

## CI/CD Pipeline
- **Tests**: Run on every PR and push to main
- **Build timeout**: 5 minutes in CI
- **Test timeout**: 15 minutes in CI (unit tests only, acceptance tests disabled)
- **Go versions tested**: 1.20, 1.21
- **Linting**: Uses golangci-lint with extensive rule set
- **Documentation**: Validated that `go generate ./...` produces no changes

## API and Testing Notes
- Acceptance tests exist but require `CRON_JOB_API_KEY` environment variable
- Unit tests mock the API client and test provider logic
- Provider supports both cron-job.org and self-hosted compatible APIs
- API key can be set via environment variable `CRON_JOB_API_KEY` or provider configuration

## Development Workflow Summary
1. **Setup**: `export PATH=$PATH:~/go/bin && make dev-setup`
2. **Build**: `make build`
3. **Test**: `make test`
4. **Lint**: `make fmt && make vet && make lint`
5. **Document**: `make docs`
6. **Install**: `make install` (for local testing)