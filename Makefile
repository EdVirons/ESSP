.PHONY: tidy test lint build-all docker-build docker-push clean help dashboard build-dashboard

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Services
SERVICES=ims-api ssot-school ssot-devices ssot-parts sync-worker

# Docker registry (override with your own)
REGISTRY?=ghcr.io/yourusername/essp

# Version from git tag or commit sha
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tidy: ## Run go mod tidy for all modules
	go work sync
	(cd shared && go mod tidy)
	(cd services/ims-api && go mod tidy)
	(cd services/ssot-school && go mod tidy)
	(cd services/ssot-devices && go mod tidy)
	(cd services/ssot-parts && go mod tidy)
	(cd services/sync-worker && go mod tidy)

test: ## Run tests for all services
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests completed successfully"

lint: ## Run golangci-lint for all services
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...
	@echo "Linting completed successfully"

build-all: ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		(cd services/$$service && $(GOBUILD) -o bin/$$service ./cmd/...) || exit 1; \
	done
	@echo "All services built successfully"

build-%: ## Build a specific service (e.g., make build-ims-api)
	@echo "Building $*..."
	@(cd services/$* && $(GOBUILD) -o bin/$* ./cmd/...)
	@echo "$* built successfully"

docker-build: ## Build Docker images for all services
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -t $(REGISTRY)/$$service:$(VERSION) -f services/$$service/Dockerfile services/$$service || exit 1; \
		docker tag $(REGISTRY)/$$service:$(VERSION) $(REGISTRY)/$$service:latest; \
	done
	@echo "All Docker images built successfully"

docker-build-%: ## Build Docker image for a specific service (e.g., make docker-build-ims-api)
	@echo "Building Docker image for $*..."
	docker build -t $(REGISTRY)/$*:$(VERSION) -f services/$*/Dockerfile services/$*
	docker tag $(REGISTRY)/$*:$(VERSION) $(REGISTRY)/$*:latest
	@echo "Docker image for $* built successfully"

docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images to $(REGISTRY)..."
	@for service in $(SERVICES); do \
		echo "Pushing $$service..."; \
		docker push $(REGISTRY)/$$service:$(VERSION); \
		docker push $(REGISTRY)/$$service:latest; \
	done
	@echo "All Docker images pushed successfully"

docker-push-%: ## Push Docker image for a specific service (e.g., make docker-push-ims-api)
	@echo "Pushing Docker image for $*..."
	docker push $(REGISTRY)/$*:$(VERSION)
	docker push $(REGISTRY)/$*:latest
	@echo "Docker image for $* pushed successfully"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		echo "Cleaning $$service..."; \
		(cd services/$$service && $(GOCLEAN) && rm -rf bin); \
	done
	@rm -f coverage.out
	@echo "Clean completed successfully"

# Dashboard targets
dashboard: ## Start dashboard development server
	@echo "Starting dashboard dev server..."
	cd dashboard && npm run dev

build-dashboard: ## Build dashboard for production
	@echo "Building dashboard..."
	cd dashboard && npm ci && npm run build
	@echo "Copying dashboard dist to ims-api embed directory..."
	rm -rf services/ims-api/internal/admin/dashboard/dist/*
	cp -r dashboard/dist/* services/ims-api/internal/admin/dashboard/dist/
	@echo "Dashboard built successfully"

dashboard-install: ## Install dashboard dependencies
	@echo "Installing dashboard dependencies..."
	cd dashboard && npm ci
	@echo "Dashboard dependencies installed"

dashboard-lint: ## Lint dashboard code
	@echo "Linting dashboard..."
	cd dashboard && npm run lint

dashboard-typecheck: ## Typecheck dashboard code
	@echo "Type-checking dashboard..."
	cd dashboard && npm run typecheck

# Build ims-api with embedded dashboard
build-ims-api-with-dashboard: build-dashboard build-ims-api ## Build ims-api with embedded dashboard

ci: lint test build-all ## Run CI pipeline locally (lint, test, build)
	@echo "CI pipeline completed successfully"
