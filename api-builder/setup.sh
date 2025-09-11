#!/bin/bash

echo "🚀 Setting up Docker Image Builder API..."

# Clean up any existing module issues
echo "📦 Cleaning up Go modules..."
go clean -modcache
rm -f go.sum

# Download dependencies with specific versions
echo "📥 Downloading dependencies..."
go mod tidy
go mod download

# Verify the build works
echo "🔨 Testing build..."
if go build -o api-builder .; then
    echo "✅ Build successful!"
    ./api-builder --version 2>/dev/null || echo "✅ Binary created successfully"
    rm -f api-builder
else
    echo "❌ Build failed"
    echo "Try running:"
    echo "  go mod tidy"
    echo "  go clean -modcache"
    echo "  go mod download"
    exit 1
fi

echo "🐳 Building Docker image..."
if docker build -t docker-image-builder:latest .; then
    echo "✅ Docker image built successfully!"
else
    echo "❌ Docker build failed"
    exit 1
fi

echo "🎉 Setup complete! You can now run:"
echo "  make docker-run    # Start with Docker Compose"
echo "  make run          # Run locally"
echo "  make test-api     # Test the API endpoints"