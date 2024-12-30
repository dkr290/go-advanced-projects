
run: build
	@./bin/todo-api

templ:
	@templ generate view

build:
	@templ generate view
	@go build -tags dev -o bin/todo-api main.go 
