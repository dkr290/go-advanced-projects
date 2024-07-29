# Introduction

Internal pypi repository written in golang and also idea is to make packages stored in azure file share and make pods scaled at least 3

# Getting Started

To start locally yhe app

```
go run cmd/web/*.go
```

# Build and Test

To build an executable

- Linux

env GOOS=linux GOARCH=amd64 go build -o bin/internal-pypi-linux-0.0.1

- Windows

env GOOS=windows GOARCH=amd64 go build -o bin/internal-pypi-windows-0.0.1.exe

- to test the application

## Build

hatch build

## Publish

The handler for publishing is `/upload`. Use the following command to publish:  
 hatch publish -u username -a password -r http://internal-pypi-url/upload

## pip install

pip install --index-url http://username:password@internal-pypi-url/simple/ --no-cache-dir my_package==0.2.0

- In browser you can browse all packages or upload one manually

# Contribute
