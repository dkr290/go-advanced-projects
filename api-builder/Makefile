# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api-builder
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=docker-image-builder
DOCKER_TAG=latest

.PHONY: all build clean test deps run docker-build docker-run docker-stop help

all: test build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOCMD) run main.go

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run with compose
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Docker logs
docker-logs:
	docker-compose logs -f api-builder

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code
lint:
	golangci-lint run

# Security scan
security:
	gosec ./...

# Test API endpoints
test-api:
	@echo "Testing health endpoint..."
	curl -f http://localhost:8080/health || echo "Health check failed"
	@echo "\nTesting API info..."
	curl -s http://localhost:8080/api/v1/ | jq '.' || echo "API info failed"

# Build sample Flask image
test-build-flask:
	curl -X POST http://localhost:8080/api/v1/build \
		-H "Content-Type: application/json" \
		-d '{"model_version":"python-flask","version":"1.0.0","name":"test-flask","tag":"latest","description":"Test Flask app"}'

# Build sample FastAPI image  
test-build-fastapi:
	curl -X POST http://localhost:8080/api/v1/build \
		-H "Content-Type: application/json" \
		-d '{"model_version":"python-fastapi","version":"1.0.0","name":"test-fastapi","tag":"latest","description":"Test FastAPI app"}'

# Build sample Node.js image
test-build-nodejs:
	curl -X POST http://localhost:8080/api/v1/build \
		-H "Content-Type: application/json" \
		-d '{"model_version":"nodejs","version":"1.0.0","name":"test-nodejs","tag":"latest","description":"Test Node.js app"}'

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the binary"
	@echo "  clean          - Clean build files"
	@echo "  test           - Run tests"
	@echo "  deps           - Download dependencies"
	@echo "  run            - Run the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose"
	@echo "  docker-logs    - Show Docker logs"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  security       - Security scan"
	@echo "  test-api       - Test API endpoints"
	@echo "  test-build-*   - Test building different image types"
	@echo "  help           - Show this help"
