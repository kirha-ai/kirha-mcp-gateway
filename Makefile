# Kirha MCP Gateway Makefile

# Build variables
BINARY_NAME=kirha-mcp-gateway
BINARY_DIR=bin
DIST_DIR=dist

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION ?= $(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS=-ldflags "-s -w \
	-X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.version=$(VERSION)' \
	-X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.commit=$(COMMIT)' \
	-X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.date=$(DATE)' \
	-X 'go.kirha.ai/kirha-mcp-gateway/cmd/cli.goVersion=$(GO_VERSION)'"

# Go build flags
BUILD_FLAGS=-trimpath $(LDFLAGS)

# Platform targets
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build test clean install wire dev docker help

# Default target
all: test build

## Build commands

# Build for current platform
build: wire
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BINARY_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd

# Build for all platforms
build-all: wire
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BINARY_NAME); \
		if [ "$$OS" = "windows" ]; then OUTPUT_NAME=$(BINARY_NAME).exe; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=0 go build $(BUILD_FLAGS) \
			-o $(DIST_DIR)/$$OS-$$ARCH/$$OUTPUT_NAME ./cmd; \
	done

# Install wire and generate code
wire:
	@echo "Installing Wire and generating dependency injection code..."
	@go install github.com/google/wire/cmd/wire@latest
	@cd di && wire

## Development commands

# Run the server locally
dev: build
	@echo "Starting development server..."
	@$(BINARY_DIR)/$(BINARY_NAME) server

# Run tests
test:
	@echo "Running tests..."
	@go test -race -coverprofile=coverage.out ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -race -coverprofile=coverage.out -v ./...

# Run tests and show coverage
test-coverage: test
	@echo "Test coverage:"
	@go tool cover -func=coverage.out

# Run tests and open coverage report in browser
test-coverage-html: test
	@echo "Opening coverage report in browser..."
	@go tool cover -html=coverage.out

# Run linting
lint:
	@echo "Running linters..."
	@go vet ./...
	@if command -v staticcheck > /dev/null; then staticcheck ./...; fi
	@if command -v golangci-lint > /dev/null; then golangci-lint run; fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

## Docker commands

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t kirha-mcp-gateway:$(VERSION) \
		-t kirha-mcp-gateway:latest .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run --rm -it \
		-e KIRHA_API_KEY="$$KIRHA_API_KEY" \
		-e KIRHA_VERTICAL="$$KIRHA_VERTICAL" \
		-p 8080:8080 \
		kirha-mcp-gateway:latest

# Run with Docker Compose
docker-compose-up:
	@echo "Starting with Docker Compose..."
	@docker-compose up --build

# Stop Docker Compose
docker-compose-down:
	@echo "Stopping Docker Compose..."
	@docker-compose down

## Utility commands

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR) $(DIST_DIR) coverage.out

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install $(BUILD_FLAGS) ./cmd

# Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"
	@echo "Go Version: $(GO_VERSION)"

# Run health check
health: build
	@echo "Running health check..."
	@$(BINARY_DIR)/$(BINARY_NAME) health

# Create release archives
release: build-all
	@echo "Creating release archives..."
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		DIR=$(DIST_DIR)/$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then \
			cd $$DIR && zip -r ../$(BINARY_NAME)-$$OS-$$ARCH.zip .; \
		else \
			cd $$DIR && tar -czf ../$(BINARY_NAME)-$$OS-$$ARCH.tar.gz .; \
		fi; \
		cd - > /dev/null; \
	done

# Show help
help:
	@echo "Kirha MCP Gateway - Available commands:"
	@echo ""
	@echo "Build commands:"
	@echo "  build          Build for current platform"
	@echo "  build-all      Build for all platforms"
	@echo "  wire           Install Wire and generate code"
	@echo ""
	@echo "Development commands:"
	@echo "  dev            Run development server"
	@echo "  test           Run tests"
	@echo "  test-verbose   Run tests with verbose output"
	@echo "  test-coverage  Run tests and show coverage"
	@echo "  test-coverage-html  Open coverage report in browser"
	@echo "  lint           Run linters"
	@echo "  fmt            Format code"
	@echo "  tidy           Tidy dependencies"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo "  docker-compose-up    Start with Docker Compose"
	@echo "  docker-compose-down  Stop Docker Compose"
	@echo ""
	@echo "Utility commands:"
	@echo "  deps           Install dependencies"
	@echo "  clean          Clean build artifacts"
	@echo "  install        Install binary to GOPATH/bin"
	@echo "  version        Show version information"
	@echo "  health         Run health check"
	@echo "  release        Create release archives"
	@echo "  help           Show this help message"