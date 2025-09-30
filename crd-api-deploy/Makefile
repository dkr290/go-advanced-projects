APP_NAME := crd-api-deploy
DOCKER_IMAGE := $(APP_NAME):latest
PORT := 8090

.DEFAULT_GOAL := help

.PHONY: build
build: ## Build the application binary
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/$(APP_NAME) ./cmd/api

.PHONY: run
run: ## Run the application locally
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	go run ./cmd/api/main.go --port $(PORT)


.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: help
help: ## Show this help message
	@echo "$(APP_NAME) - Available targets:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
