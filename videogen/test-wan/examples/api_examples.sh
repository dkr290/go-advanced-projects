#!/bin/bash
# Example API requests for Wan2.1 Video Server

API_URL="http://localhost:8080"

echo "=== Wan2.1 Video Server API Examples ==="
echo ""

# 1. Health Check
echo "1. Health Check"
echo "GET $API_URL/health"
curl -s "$API_URL/health" | jq '.'
echo ""
echo ""

# 2. Model Info
echo "2. Model Information"
echo "GET $API_URL/api/v1/model/info"
curl -s "$API_URL/api/v1/model/info" | jq '.'
echo ""
echo ""

# 3. Text-to-Video Generation
echo "3. Text-to-Video Generation"
echo "POST $API_URL/api/v1/generate/text-to-video"
JOB_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/generate/text-to-video" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A beautiful sunset over the ocean with waves crashing on the shore",
    "negative_prompt": "blurry, low quality",
    "num_frames": 64,
    "fps": 24,
    "width": 512,
    "height": 512,
    "guidance_scale": 7.5,
    "num_inference_steps": 50,
    "seed": 42
  }')
echo "$JOB_RESPONSE" | jq '.'
JOB_ID=$(echo "$JOB_RESPONSE" | jq -r '.job_id')
echo ""
echo ""

# 4. Check Job Status
if [ ! -z "$JOB_ID" ] && [ "$JOB_ID" != "null" ]; then
  echo "4. Check Job Status"
  echo "GET $API_URL/api/v1/job/$JOB_ID"
  sleep 2
  curl -s "$API_URL/api/v1/job/$JOB_ID" | jq '.'
  echo ""
  echo ""
fi

# 5. Image-to-Video (requires an image file)
echo "5. Image-to-Video Generation (example command)"
echo "curl -X POST $API_URL/api/v1/generate/image-to-video \\"
echo "  -F 'image=@/path/to/your/image.jpg' \\"
echo "  -F 'prompt=A cat playing with a ball' \\"
echo "  -F 'num_frames=64' \\"
echo "  -F 'fps=24'"
echo ""
echo ""

# 6. Video-to-Video (requires a video file)
echo "6. Video-to-Video Generation (example command)"
echo "curl -X POST $API_URL/api/v1/generate/video-to-video \\"
echo "  -F 'video=@/path/to/your/video.mp4' \\"
echo "  -F 'prompt=Transform into anime style' \\"
echo "  -F 'strength=0.8' \\"
echo "  -F 'fps=24'"
echo ""
echo ""

# 7. List Models
echo "7. List Available Models"
echo "GET $API_URL/api/v1/models"
curl -s "$API_URL/api/v1/models" | jq '.'
echo ""
echo ""

echo "=== Examples Complete ==="
echo ""
echo "For actual image/video uploads, create test files and replace paths in commands above."
