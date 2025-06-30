.PHONY: test build clean run help

# Variables
BINARY_NAME=voice-notify-mcp

# Default target
all: test build

# Build the binary
build:
	go build -o $(BINARY_NAME) -v

# Run tests
test:
	go test -v ./...

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
	@echo "  test   - Run tests"
	@echo "  build  - Build the binary"
	@echo "  clean  - Remove build artifacts"
	@echo "  run    - Build and run the application"
	@echo "  help   - Display this help message"