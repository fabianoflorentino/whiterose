.PHONY: build test test-cover lint fmt vet clean install run help deps

# Build variables
BINARY_NAME=whiterose
GO=go
OUTPUT_DIR=bin

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo "  lint       - Run linters"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install binary to GOBIN"
	@echo "  run        - Build and run CLI"

# Build the binary
build:
	$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME) .

# Build for multiple platforms
build-all: clean
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 $(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Run tests
test:
	$(GO) test -short ./...

# Run tests with coverage
test-cover:
	$(GO) test -short -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -func=coverage.out

# Run linters (requires golangci-lint)
lint:
	golangci-lint run ./...

# Format code
fmt:
	$(GO) fmt ./...
	gofumpt -w .

# Run go vet
vet:
	$(GO) vet ./...

# Clean build artifacts
clean:
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.out

# Install binary to GOBIN
install: build
	$(GO) install .

# Build and run CLI
run: build
	./$(OUTPUT_DIR)/$(BINARY_NAME)

# Run specific command
run-cmd: build
	./$(OUTPUT_DIR)/$(BINARY_NAME) $(CMD)