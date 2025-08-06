SHELL := /bin/bash

.PHONY: help build test clean install fmt vet

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the provider
	@echo "Building terraform-provider-uptimekuma..."
	go build -o terraform-provider-uptimekuma

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-acc: ## Run acceptance tests
	@echo "Running acceptance tests..."
	TF_ACC=1 go test -v ./internal/provider/

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f terraform-provider-uptimekuma
	rm -rf dist/
	go clean

install: build ## Build and install the provider locally
	@echo "Installing provider locally..."
	mkdir -p ~/.terraform.d/plugins/j0r15.local/provider/uptimekuma/1.0.0/linux_amd64/
	cp terraform-provider-uptimekuma ~/.terraform.d/plugins/j0r15.local/provider/uptimekuma/1.0.0/linux_amd64/

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	go mod tidy

release: clean ## Build release artifacts
	@echo "Building release artifacts..."
	./build.sh

docs: ## Generate documentation
	@echo "Generating documentation..."
	go generate

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

all: fmt vet test build ## Run fmt, vet, test, and build
