# Contributing to terraform-provider-cronjoborg

Thank you for your interest in contributing to the Terraform Provider for cron-job.org!

## Development Environment Setup

### Prerequisites

- [Go](https://golang.org/doc/install) 1.20+ (to build the provider plugin)
- [Terraform](https://www.terraform.io/downloads.html) 1.0+ (to test the provider)
- [Git](https://git-scm.com/downloads)

### Setting up the development environment

1. Clone the repository:
   ```bash
   git clone https://github.com/plain-insure/terraform-provider-cronjoborg.git
   cd terraform-provider-cronjoborg
   ```

2. Install development tools:
   ```bash
   make dev-setup
   ```

3. Build the provider:
   ```bash
   make build
   ```

## Development Workflow

### Building

To compile the provider:
```bash
make build
```

### Testing

Run unit tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-cover
```

Run acceptance tests (requires API key):
```bash
export CRON_JOB_API_KEY="your-api-key"
make test-acc
```

### Code Quality

Format code:
```bash
make fmt
```

Run linters:
```bash
make lint
```

Run go vet:
```bash
make vet
```

### Local Development

To install the provider locally for testing:
```bash
make install
```

This will install the provider to your local Terraform plugins directory.

### Testing Local Changes

1. Build and install the provider locally:
   ```bash
   make install
   ```

2. Create a test configuration using the local provider:
   ```hcl
   terraform {
     required_providers {
       cronjoborg = {
         source  = "registry.terraform.io/plain-insure/cronjoborg"
         version = "dev"
       }
     }
   }
   ```

3. Run terraform commands to test your changes:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Contributing Code

### Pull Request Process

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Add tests for your changes
5. Run the test suite and ensure all tests pass
6. Run linters and fix any issues
7. Update documentation if needed
8. Commit your changes with clear commit messages
9. Push to your fork and submit a pull request

### Commit Message Guidelines

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add support for cron job scheduling

- Implement schedule field in job resource
- Add validation for cron expressions
- Update documentation and examples

Fixes #123
```

### Code Style

- Follow Go best practices and idioms
- Use `gofmt` to format your code
- Ensure your code passes all linters
- Add comments for exported functions and types
- Write clear, readable code

### Adding New Resources

When adding a new resource:

1. Create the resource file in the `provider/` directory
2. Add the resource to the provider's ResourcesMap
3. Write comprehensive tests
4. Add documentation and examples
5. Update the main README if needed

### Documentation

- Update relevant documentation for any changes
- Add examples for new features
- Ensure all public APIs are documented
- Use clear, concise language

## Issue Reporting

When reporting issues:

1. Use the issue templates if available
2. Provide a clear description of the problem
3. Include steps to reproduce
4. Add relevant logs and error messages
5. Specify your environment (OS, Terraform version, provider version)

## Getting Help

- Check existing issues and documentation first
- Use GitHub issues for bug reports and feature requests
- Be respectful and constructive in all interactions

Thank you for contributing!