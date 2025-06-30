.PHONY: test build clean run lint help

# Variables
BINARY_NAME=voice-notify-mcp
GOLANGCI_LINT=golangci-lint

# Default target
all: lint test build

# Build the binary
build:
	go build -o $(BINARY_NAME) -v

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	@if command -v $(GOLANGCI_LINT) > /dev/null; then \
		$(GOLANGCI_LINT) run; \
	else \
		echo "golangci-lint is not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		$(GOLANGCI_LINT) run; \
	fi

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Run the application
run: build
	./$(BINARY_NAME)

# Display help
help:
	@echo "Available targets:"
	@echo "  lint   - Run golangci-lint"
	@echo "  test   - Run tests"
	@echo "  build  - Build the binary"
	@echo "  clean  - Remove build artifacts"
	@echo "  run    - Build and run the application"
	@echo "  help   - Display this help message"