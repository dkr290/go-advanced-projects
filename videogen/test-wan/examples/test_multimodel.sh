#!/bin/bash
# Test all available video models

API_URL="http://localhost:8080"
PYTHON_URL="http://localhost:5000"

echo "=================================="
echo "  Multi-Model Test Suite"
echo "=================================="
echo ""

# Check server health
echo "1. Checking server health..."
HEALTH=$(curl -s "$PYTHON_URL/health")
echo "$HEALTH" | jq '.'
echo ""

# List available models
echo "2. Listing available models..."
MODELS=$(curl -s "$PYTHON_URL/api/models")
echo "$MODELS" | jq '.'
echo ""

# Test each model
MODELS_LIST=("modelscope" "zeroscope" "ltx-video")

for MODEL in "${MODELS_LIST[@]}"; do
    echo "=================================="
    echo "  Testing: $MODEL"
    echo "=================================="
    
    # Switch model
    echo "Switching to $MODEL..."
    SWITCH_RESULT=$(curl -s -X POST "$PYTHON_URL/api/switch-model" \
        -H "Content-Type: application/json" \
        -d "{\"model_name\": \"$MODEL\"}")
    
    echo "$SWITCH_RESULT" | jq '.'
    
    # Wait for model to load
    echo "Waiting for model to load..."
    sleep 5
    
    # Generate test video
    echo "Generating test video with $MODEL..."
    RESULT=$(curl -s -X POST "$API_URL/api/v1/generate/text-to-video" \
        -H "Content-Type: application/json" \
        -d '{
            "prompt": "A cat playing with a ball",
            "num_frames": 16,
            "width": 256,
            "height": 256,
            "num_inference_steps": 20
        }')
    
    echo "$RESULT" | jq '.'
    echo ""
    
    # Small delay between tests
    sleep 2
done

echo "=================================="
echo "  All tests complete!"
echo "=================================="
