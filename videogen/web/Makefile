# Makefile for videogen/web

.PHONY: help run build clean dev test

# Default target
help: ## Show this help message
	@echo "Wan2.1 Video Generator - Web UI"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the web server
	@echo "🚀 Starting web UI on http://localhost:3000"
	@go run main.go

build: ## Build the binary
	@echo "🔨 Building binary..."
	@go build -o videogen-web main.go
	@echo "✅ Built: videogen-web"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -f videogen-web
	@rm -rf outputs/*
	@echo "✅ Clean complete"

dev: ## Run with auto-reload (requires air)
	@echo "🔥 Running in development mode..."
	@air || (echo "Install air: go install github.com/cosmtrek/air@latest" && exit 1)

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test ./...

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

setup: ## Initial setup
	@echo "⚙️  Running setup..."
	@chmod +x setup.sh
	@./setup.sh

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t videogen-web:latest .

docker-run: ## Run in Docker
	@echo "🐳 Running Docker container..."
	@docker run -p 3000:3000 -e API_BASE_URL=http://host.docker.internal:8080 videogen-web:latest

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run || (echo "Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)

fmt: ## Format code
	@echo "💅 Formatting code..."
	@go fmt ./...

.DEFAULT_GOAL := help
