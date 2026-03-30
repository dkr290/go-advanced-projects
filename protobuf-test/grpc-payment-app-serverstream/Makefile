.PHONY: help 
help:


proto-gen:
	protoc -Iproto \
		--go_out=. --go_opt=module=github.com/dkr290/go-advanced-projects/protobuf-test/grpc-payment-app-serverstream \
	  --go-grpc_out=. --go-grpc_opt=module=github.com/dkr290/go-advanced-projects/protobuf-test/grpc-payment-app-serverstream \
		proto/payment.proto 
	go mod tidy


run-client:
	go run client/*.go 

run-server:
	go run server/*.go
