# Makefile for Wan2.1 Video Server

.PHONY: help build run test clean install setup-python docker

# Variables
BINARY_NAME=wan2-video-server
VERSION?=1.0.0
BUILD_DIR=./build
GO_FILES=$(shell find . -name '*.go' -type f)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Display this help message
	@echo "Wan2.1 Video Server - Makefile commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Build the Go server
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-linux: ## Build for Linux
	@echo "$(GREEN)Building for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64$(NC)"

build-all: build build-linux ## Build for all platforms
	@echo "$(GREEN)All builds complete$(NC)"

run: ## Run the Go server
	@echo "$(GREEN)Starting $(BINARY_NAME)...$(NC)"
	@go run main.go

run-python: ## Run the Python backend
	@echo "$(GREEN)Starting Python backend...$(NC)"
	@cd python_backend && source venv/bin/activate && python server.py

install: ## Install Go dependencies
	@echo "$(GREEN)Installing Go dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies installed$(NC)"

setup-python: ## Setup Python backend
	@echo "$(GREEN)Setting up Python backend...$(NC)"
	@cd python_backend && chmod +x setup.sh && ./setup.sh
	@echo "$(GREEN)Python backend setup complete$(NC)"

setup: install setup-python ## Setup both Go and Python environments
	@echo "$(GREEN)Full setup complete!$(NC)"
	@cp .env.example .env
	@echo "$(YELLOW)Don't forget to edit .env with your configuration$(NC)"

test: ## Run Go tests
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Tests complete. Coverage report: coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(RED)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -rf uploads/* outputs/*
	@echo "$(GREEN)Clean complete$(NC)"

clean-all: clean ## Clean everything including models and Python venv
	@echo "$(RED)Cleaning everything...$(NC)"
	@rm -rf models/*
	@rm -rf python_backend/venv
	@echo "$(GREEN)Deep clean complete$(NC)"

download-model: ## Download the model from Hugging Face
	@echo "$(GREEN)Downloading model...$(NC)"
	@go run main.go download

lint: ## Run linters
	@echo "$(GREEN)Running linters...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed$(NC)" && exit 1)
	@golangci-lint run ./...

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@gofmt -s -w .

dev: ## Run in development mode
	@echo "$(GREEN)Starting development environment...$(NC)"
	@echo "$(YELLOW)Make sure Python backend is running in another terminal$(NC)"
	@echo "$(YELLOW)Run: make run-python$(NC)"
	@echo ""
	@go run main.go --log-level debug

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t wan2-video-server:$(VERSION) .

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run -p 8080:8080 -p 5000:5000 --gpus all wan2-video-server:$(VERSION)

.DEFAULT_GOAL := help
