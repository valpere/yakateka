# YakaTeka Makefile

# Application name
APP_NAME := yakateka
VERSION := 0.1.0
BUILD_DIR := bin
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# Build flags
LDFLAGS := -ldflags "-X github.com/valpere/yakateka/cmd.version=$(VERSION)"

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: clean fmt lint test build ## Run all: clean, format, lint, test, and build

.PHONY: build
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

.PHONY: install
install: ## Install the application to $GOPATH/bin
	@echo "Installing $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(APP_NAME) main.go
	@echo "Installed to $(GOPATH)/bin/$(APP_NAME)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) -s -w $(GO_FILES)
	@echo "Format complete"

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@test -z "$$($(GOFMT) -l $(GO_FILES))" || (echo "Code is not formatted, run 'make fmt'"; exit 1)

.PHONY: lint
lint: ## Run linters (go vet)
	@echo "Running linters..."
	$(GOVET) ./...
	@echo "Lint complete"

.PHONY: lint-full
lint-full: ## Run full linting with golangci-lint (requires golangci-lint)
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install from https://golangci-lint.run/usage/install/"; exit 1)
	golangci-lint run ./...

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies downloaded"

.PHONY: verify
verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GOMOD) verify

.PHONY: run
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run main.go

.PHONY: dev
dev: ## Run in development mode (rebuild on changes, requires entr)
	@which entr > /dev/null || (echo "entr not found, install with 'apt install entr' or 'brew install entr'"; exit 1)
	@echo "Running in development mode (watching for changes)..."
	find . -name '*.go' | entr -r make run

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest

.PHONY: version
version: ## Show version
	@echo "$(APP_NAME) version $(VERSION)"
