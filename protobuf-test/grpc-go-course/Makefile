.PHONY: help run-calculator-server build clean
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?##' Makefile | \
        sed -E 's/:.*##/: /' | \
        column -t -s ':'

gen: 
	protoc -Igreet/proto \
		  --go_out=. --go_opt=module=github.com/dkr290-go-advanced-projects/grpc-go-course \
			--go-grpc_out=. --go-grpc_opt=module=github.com/dkr290-go-advanced-projects/grpc-go-course \
			greet/proto/dummy.proto

.PHONY: generate greet
gen-greet: ## Generate greet pb go files
	protoc --go_out=. --go_opt=paths=source_relative \
		     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		     greet/proto/greet.proto
	go mod tidy

.PHONY: run greet server
run-greet-server: ## Run greet server
	go run greet/server/*.go

.PHONY: run greet client
run-greet-client: ## Run greet client 
	go run greet/client/*.go

.PHONY: run generate calculate
gen-calculate: ## Run generate of calculator pb go files
	protoc --go_out=. --go_opt=paths=source_relative \
		     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		     calculator/proto/calculate.proto
	go mod tidy


.PHONY: run calculator server
run-calculate-server: ## Run calculator server 
	go run calculator/server/*.go 

run-calculate-client: ## Run the calculator server with optional ARGS="..."
	go run calculator/client/*.go  $(ARGS)
