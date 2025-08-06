BINARY_NAME=terraform-provider-cronjob
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

.PHONY: help
help: ## Show this help message
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the provider binary
	go build ${LDFLAGS} -o bin/${BINARY_NAME}

.PHONY: install
install: build ## Install the provider locally
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/plain-insure/cronjob/dev/linux_amd64/
	cp bin/${BINARY_NAME} ~/.terraform.d/plugins/registry.terraform.io/plain-insure/cronjob/dev/linux_amd64/

.PHONY: test
test: ## Run unit tests
	go test -v ./...

.PHONY: test-acc
test-acc: ## Run acceptance tests
	TF_ACC=1 go test -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## Run linters
	golangci-lint run

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f ${BINARY_NAME}

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: docs
docs: ## Generate documentation
	go generate ./...

.PHONY: release-dry-run
release-dry-run: ## Run goreleaser in dry-run mode
	goreleaser release --snapshot --skip-publish --rm-dist

.PHONY: dev-setup
dev-setup: ## Set up development environment
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	@echo "Development tools installed!"

.DEFAULT_GOAL := help