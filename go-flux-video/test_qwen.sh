#!/bin/bash

# Test script for Qwen-Image-Edit integration
echo "Testing Qwen-Image-Edit integration..."

# Create test directories
mkdir -p test_output
mkdir -p test_images

# Create a simple test image (placeholder)
echo "Creating test image..."
convert -size 512x512 xc:white -fill blue -draw 'circle 256,256 256,400' test_images/test_input.png 2>/dev/null || \
echo "Note: ImageMagick not available, using placeholder image"

# Create test config
cat > test_config.json << EOF
{
  "style_suffix": "digital art, vibrant colors",
  "negative_prompt": "blurry, distorted",
  "prompts": [
    "a blue circle transformed into a red square",
    "abstract geometric patterns"
  ]
}
EOF

echo "Test setup complete!"
echo ""
echo "To test Qwen-Image-Edit, run:"
echo ""
echo "  ./gfluxgo \\"
echo "    --config test_config.json \\"
echo "    --hf-model \"Qwen/Qwen-Image-Edit\" \\"
echo "    --use-qwen \\"
echo "    --img2img \\"
echo "    --strength 0.7 \\"
echo "    --output test_output \\"
echo "    --seed 42 \\"
echo "    --resolution 512x512 \\"
echo "    --steps 20 \\"
echo "    --guidence_scale 3.0"
echo ""
echo "Or for text-to-image mode (without --img2img flag):"
echo ""
echo "  ./gfluxgo \\"
echo "    --config test_config.json \\"
echo "    --hf-model \"Qwen/Qwen-Image-Edit\" \\"
echo "    --use-qwen \\"
echo "    --output test_output \\"
echo "    --seed 42 \\"
echo "    --resolution 512x512 \\"
echo "    --steps 20 \\"
echo "    --guidence_scale 3.0"
echo ""
echo "Note: This will download the Qwen-Image-Edit model (~15GB) on first run."