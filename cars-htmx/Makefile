.PHONY: build test run

default: run

build:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css
	@templ generate view
	@go build -o bin/gocarshtmx main.go 

test:
	@go test -v ./...
	
run: build
	@./bin/gocarshtmx


install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
#	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss@3
	@npm install -D daisyui@latest

tailwind:
	@npx tailwindcss -i views/css/styles.css -o ./public/styles.css --watch


templ:
	@templ generate -watch -proxy=http://localhost:3000

