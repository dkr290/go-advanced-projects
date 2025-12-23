#!/bin/bash

echo "=== Testing Web Server ==="
echo ""

# Create test directory
mkdir -p test_web_output

echo "1. Checking template file..."
if [ ! -f "templates/gallery.html" ]; then
    echo "❌ Template file not found: templates/gallery.html"
    exit 1
fi
echo "✓ Template file exists"

echo ""
echo "2. Testing web-only mode (no config required)..."
echo "Command: ./gfluxgo --web --output ./test_web_output --web-port 8082"
echo ""
echo "Open browser to: http://localhost:8082"
echo "You should see:"
echo "  - Web interface with upload section"
echo "  - Empty gallery (no images yet)"
echo "  - Ability to upload images to ./images/"
echo ""
echo "Press Ctrl+C to stop the web server"
echo ""
echo "3. Testing with existing images..."
echo "Add some images to test_web_output directory:"
echo "  cp /path/to/any/image.png test_web_output/"
echo "Then restart web server to see them in gallery"
echo ""
echo "✅ Web server should work without any config file!"