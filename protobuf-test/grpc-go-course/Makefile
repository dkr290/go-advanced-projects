gen:
	protoc -Igreet/proto \
		  --go_out=. --go_opt=module=github.com/dkr290-go-advanced-projects/grpc-go-course \
			--go-grpc_out=. --go-grpc_opt=module=github.com/dkr290-go-advanced-projects/grpc-go-course \
			greet/proto/dummy.proto

.PHONY: generate greet
gen-greet:
	protoc --go_out=. --go_opt=paths=source_relative \
		     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		     greet/proto/greet.proto
	go mod tidy

.PHONY: run greet server
run-greet-server:
	go run greet/server/*.go

.PHONY: run greet client
run-greet-client:
	go run greet/client/*.go


