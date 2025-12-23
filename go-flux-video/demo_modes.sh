#!/bin/bash

echo "=== GFluxGo Usage Modes Demo ==="
echo ""

# Create test directories
mkdir -p demo_web_only
mkdir -p demo_generation

echo "1. Testing Web-Only Mode (no config required)"
echo "--------------------------------------------"
echo "Starting web server on port 8081..."
echo "Open browser to: http://localhost:8081"
echo "Upload images to ./images/ directory via web interface"
echo ""
echo "Command: ./gfluxgo --web --web-port 8081 --output ./demo_web_only"
echo ""
echo "Press Ctrl+C in the terminal to stop the web server"
echo ""

echo "2. Testing Image Generation Mode (requires config)"
echo "--------------------------------------------------"
echo "First, create a config file:"
cat > demo_config.json << 'EOF'
{
  "style_suffix": "digital art",
  "negative_prompt": "blurry, distorted",
  "prompts": [
    "a beautiful sunset",
    "a futuristic city"
  ]
}
EOF

echo "Config created: demo_config.json"
echo ""
echo "Then generate images:"
echo "Command: ./gfluxgo --config demo_config.json --output ./demo_generation"
echo ""
echo "Or generate and serve:"
echo "Command: ./gfluxgo --config demo_config.json --web --output ./demo_generation"
echo ""

echo "3. Testing Combined Workflow"
echo "---------------------------"
echo "Step 1: Start web server for uploading"
echo "  ./gfluxgo --web --web-port 8082 --output ./demo_workflow"
echo ""
echo "Step 2: Upload images via browser"
echo "  http://localhost:8082"
echo ""
echo "Step 3: Generate images (in another terminal)"
echo "  ./gfluxgo --config demo_config.json --img2img --output ./demo_workflow"
echo ""
echo "Step 4: View results in browser"
echo "  http://localhost:8082 (auto-refreshes)"
echo ""

echo "Quick Commands Reference:"
echo "-------------------------"
echo "Web only:          ./gfluxgo --web --output ./gallery"
echo "Generate only:     ./gfluxgo --config prompts.json --output ./results"
echo "Generate + serve:  ./gfluxgo --config prompts.json --web --output ./results"
echo "Custom port:       ./gfluxgo --web --web-port 9000 --output ./gallery"
echo "With Qwen model:   ./gfluxgo --config prompts.json --use-qwen --img2img --output ./results"
echo ""

echo "Note: For web-only mode, NO config file is needed!"
echo "Note: For image generation, config file IS required!"