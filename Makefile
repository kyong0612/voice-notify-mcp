.PHONY: all build test lint fmt clean run install help

# Variables
BINARY_NAME=voice-notify-mcp
GO=go
GOLANGCI_LINT=golangci-lint
GOFILES=$(shell find . -name "*.go" -type f)

# Default target
all: lint test build

# Build the binary
build:
	@echo "Building..."
	@$(GO) build -o $(BINARY_NAME) -v

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests (macOS only)
test-integration:
	@echo "Running integration tests..."
	@$(GO) test -v -tags=integration ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@$(GO) test -bench=. -benchmem ./...

# Run linter
lint:
	@echo "Running linter..."
	@if command -v $(GOLANGCI_LINT) > /dev/null; then \
		$(GOLANGCI_LINT) run ./...; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(GOLANGCI_LINT) run ./...; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@goimports -w $(GOFILES)

# Run go mod tidy
tidy:
	@echo "Tidying modules..."
	@$(GO) mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -f *.tar.gz

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Install the binary to GOPATH/bin
install: build
	@echo "Installing..."
	@$(GO) install

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Tools installed successfully"

# Run security checks
security:
	@echo "Running security checks..."
	@gosec ./...
	@govulncheck ./...

# Run all checks (lint, test, security)
check: lint test security

# Create a new release tag
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Tag created. Push with: git push origin $(VERSION)"

# Pre-commit hooks
pre-commit:
	@echo "Installing pre-commit hooks..."
	@if command -v pre-commit > /dev/null; then \
		pre-commit install; \
	else \
		echo "pre-commit not found. Install with: pip install pre-commit"; \
		exit 1; \
	fi

# Run pre-commit on all files
pre-commit-all:
	@echo "Running pre-commit on all files..."
	@pre-commit run --all-files

# Display help
help:
	@echo "Available targets:"
	@echo "  all            - Run lint, test, and build (default)"
	@echo "  build          - Build the binary"
	@echo "  test           - Run unit tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-integration - Run integration tests (macOS only)"
	@echo "  bench          - Run benchmarks"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Run go mod tidy"
	@echo "  clean          - Remove build artifacts"
	@echo "  run            - Build and run the application"
	@echo "  install        - Install the binary to GOPATH/bin"
	@echo "  install-tools  - Install development tools"
	@echo "  security       - Run security checks"
	@echo "  check          - Run all checks (lint, test, security)"
	@echo "  release        - Create a new release tag (VERSION=v1.0.0)"
	@echo "  pre-commit     - Install pre-commit hooks"
	@echo "  pre-commit-all - Run pre-commit on all files"
	@echo "  help           - Display this help message"