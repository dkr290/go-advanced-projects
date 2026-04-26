
.PHONY: help run-calculator-server build clean
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?##' Makefile | \
        sed -E 's/:.*##/: /' | \
        column -t -s ':'



SSL_DIR   := cert
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



.PHONY: generate
gen-proto: ## Generate pb go files
	protoc -I=proto \
	       --go_out=. --go_opt=module=github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api \
	       --go-grpc_out=. --go-grpc_opt=module=github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api \
	       proto/main.proto proto/students.proto proto/execs.proto
	go mod tidy

.PHONY: run-server
run-server: ## run the server
	go run cmd/grpcapi/*.go
