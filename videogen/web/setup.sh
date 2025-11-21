#!/bin/bash
# Quick setup script for videogen/web

echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║      Wan2.1 Video Generator - Web UI Setup                      ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+"
    exit 1
fi
echo "✅ Go $(go version | awk '{print $3}')"

# Install dependencies
echo ""
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies installed"

# Create .env if it doesn't exist
if [ ! -f ".env" ]; then
    echo ""
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env created"
else
    echo "⚠️  .env already exists, skipping"
fi

# Create directories
echo ""
echo "📁 Creating directories..."
mkdir -p static/css static/js static/images
mkdir -p templates/layouts templates/pages templates/components
mkdir -p handlers middleware
mkdir -p outputs
echo "✅ Directories created"

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║      Setup Complete!                                             ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
echo "🎯 Next steps:"
echo ""
echo "1. Make sure the main API server is running:"
echo "   cd ../../  # Go to main project"
echo "   # Start Python backend and Go API server"
echo ""
echo "2. Start this web interface:"
echo "   go run main.go"
echo ""
echo "3. Open your browser:"
echo "   http://localhost:3000"
echo ""
echo "📖 Read README.md for more information"
echo ""
