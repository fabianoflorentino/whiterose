.PHONY: build test test-cover lint fmt vet clean install run help deps \
	docker-build docker-run docker-dev docker-prod docker-clean docker-push docker-tag

# Build variables
BINARY_NAME=whiterose
GO=go
OUTPUT_DIR=bin
DOCKER_REGISTRY=docker.io
DOCKER_IMAGE=$(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_TAG=$(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build:"
	@echo "  build        - Build the binary"
	@echo "  build-all   - Build for multiple platforms"
	@echo ""
	@echo "Test:"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint        - Run linters"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run container"
	@echo "  docker-dev    - Run in development mode"
	@echo "  docker-prod   - Run in production mode"
	@echo "  docker-push   - Push image to registry"
	@echo "  docker-clean - Remove local images"
	@echo ""
	@echo "Utility:"
	@echo "  clean      - Clean build artifacts"
	@echo "  install   - Install binary to GOBIN"
	@echo "  run       - Build and run CLI"

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

# ============================================
# Docker targets
# ============================================

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -t $(DOCKER_IMAGE):latest .

# Run container (production)
docker-run:
	docker run --rm -it $(DOCKER_IMAGE):latest

# Run in development mode
docker-dev:
	docker build --target development -t $(BINARY_NAME)-dev .
	docker run --rm -it -v "$$(pwd):/whiterose" -w /whiterose $(BINARY_NAME)-dev

# Run production container
docker-prod:
	docker build --target production -t $(DOCKER_IMAGE)-prod .
	docker run --rm -it $(DOCKER_IMAGE)-prod

# Push image to registry
docker-push: docker-build
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

# Tag image
docker-tag:
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):$(DOCKER_TAG)

# Clean Docker images
docker-clean:
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest $(DOCKER_IMAGE)-prod 2>/dev/null || true