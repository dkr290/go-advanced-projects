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

.PHONY: generate blog
gen-blog: ## Generate blog pb go files 
	protoc --go_out=. --go_opt=paths=source_relative \
		     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		     blog/proto/blog.proto
	go mod tidy


.PHONY: run greet server
run-greet-server: ## Run greet server
	go run greet/server/*.go

.PHONY: run blog server
run-blog-server: ## Run blog server
	go run blog/server/*.go


.PHONY: run greet client
run-greet-client: ## Run greet client 
	go run greet/client/*.go

.PHONY: run generate calculate
gen-calculate: ## Run generate of calculator pb go files
	protoc  -Icalculator/proto \
		  --go_out=. --go_opt=module=github.com/dkr290-go-advanced-projects/grpc-go-course \
			--go-grpc_out=. --go-grpc_opt=module=github.com/dkr290-go-advanced-projects/grpc-go-course \
		     calculator/proto/calculate.proto \
	       calculator/proto/primes.proto \
	       calculator/proto/avg.proto \
				 calculator/proto/max.proto \
				 calculator/proto/sqrt.proto
	go mod tidy
 

.PHONY: run calculator server
run-calculate-server: ## Run calculator server
	go run calculator/server/*.go

run-calculate-client: ## Run the calculator server with optional ARGS="..."
	go run calculator/client/*.go  $(ARGS)

.PHONY: run generate calculate clean-proto
clean-calculator-proto: ## Clean calculator go files 
	rm -f calculator/proto/*.pb.go calculator/proto/*_grpc.pb.go


SSL_DIR   := ssl
DAYS      := 3650
BITS      := 4096
SERVER_CN := localhost
CLIENT_CN := client

.PHONY: certs clean-certs

certs: $(SSL_DIR) 
	# CA
	openssl genrsa -out $(SSL_DIR)/ca.key $(BITS)
	openssl req -new -x509 -key $(SSL_DIR)/ca.key -out $(SSL_DIR)/ca.crt -days $(DAYS) \
		-subj "/CN=MyGRPC-CA/O=MyOrg/C=US"
	# Server
	openssl genrsa -out $(SSL_DIR)/server.key $(BITS)
	openssl req -new -key $(SSL_DIR)/server.key -out $(SSL_DIR)/server.csr \
		-subj "/CN=$(SERVER_CN)/O=MyOrg/C=US"
	printf "[req]\ndistinguished_name=req\n[SAN]\nsubjectAltName=DNS:$(SERVER_CN),IP:127.0.0.1\n" \
		> $(SSL_DIR)/server_ext.cnf
	openssl x509 -req -in $(SSL_DIR)/server.csr -CA $(SSL_DIR)/ca.crt -CAkey $(SSL_DIR)/ca.key \
		-CAcreateserial -out $(SSL_DIR)/server.crt -days $(DAYS) \
		-extfile $(SSL_DIR)/server_ext.cnf -extensions SAN
	# Client
	openssl genrsa -out $(SSL_DIR)/client.key $(BITS)
	openssl req -new -key $(SSL_DIR)/client.key -out $(SSL_DIR)/client.csr \
		-subj "/CN=$(CLIENT_CN)/O=MyOrg/C=US"
	openssl x509 -req -in $(SSL_DIR)/client.csr -CA $(SSL_DIR)/ca.crt -CAkey $(SSL_DIR)/ca.key \
		-CAcreateserial -out $(SSL_DIR)/client.crt -days $(DAYS)
	# Cleanup
	rm -f $(SSL_DIR)/*.csr $(SSL_DIR)/*.srl $(SSL_DIR)/server_ext.cnf

$(SSL_DIR):
	mkdir -p $(SSL_DIR)

clean-certs:  ## Clean certificates 
	rm -rf $(SSL_DIR)

