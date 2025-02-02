.PHONY: build test run

default: run

build:
	@templ generate view
	@go build -o bin/gocarshtmx main.go 

test:
	@go test -v ./...
	
run: build
	@./bin/gocarshtmx

