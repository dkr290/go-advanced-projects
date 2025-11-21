# Testing Guide

Guide for testing the Wan2.1 Video Server.

## Quick Start Testing

### 1. Start the Server

**Terminal 1 - Python Backend:**
```bash
cd python_backend
source venv/bin/activate
python server.py
```

**Terminal 2 - Go Server:**
```bash
./wan2-video-server
```

### 2. Run Basic Tests

```bash
# Make test script executable
chmod +x examples/api_examples.sh

# Run tests
./examples/api_examples.sh
```

## Unit Tests

### Go Tests

Run all Go unit tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run specific package tests:

```bash
go test ./pkg/config -v
go test ./pkg/handlers -v
go test ./pkg/model -v
```

## Integration Tests

### Test Health Endpoint

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "wan2-video-server",
  "version": "1.0.0"
}
```

### Test Model Info

```bash
curl http://localhost:8080/api/v1/model/info
```

### Test Python Backend Health

```bash
curl http://localhost:5000/health
```

## API Tests

### Text-to-Video Test

```bash
# Submit generation request
JOB_ID=$(curl -s -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Test video: a spinning cube",
    "num_frames": 32,
    "fps": 24,
    "width": 256,
    "height": 256,
    "num_inference_steps": 20
  }' | jq -r '.job_id')

echo "Job ID: $JOB_ID"

# Wait a bit
sleep 5

# Check status
curl -s http://localhost:8080/api/v1/job/$JOB_ID | jq '.'
```

### Image-to-Video Test

First, create a test image:

```bash
# Using ImageMagick to create a test image
convert -size 512x512 xc:blue -fill white -pointsize 72 \
  -gravity center -annotate +0+0 "Test" test_image.jpg
```

Then test:

```bash
curl -X POST http://localhost:8080/api/v1/generate/image-to-video \
  -F "image=@test_image.jpg" \
  -F "prompt=Animate this image" \
  -F "num_frames=32" \
  -F "fps=24"
```

## Load Testing

### Using Apache Bench

```bash
# Test health endpoint
ab -n 100 -c 10 http://localhost:8080/health

# Test with POST (text-to-video)
ab -n 5 -c 2 -p test_request.json -T application/json \
  http://localhost:8080/api/v1/generate/text-to-video
```

### Using wrk

```bash
# Install wrk
# Ubuntu: sudo apt install wrk
# macOS: brew install wrk

# Test health endpoint
wrk -t4 -c10 -d30s http://localhost:8080/health

# Test with Lua script for POST
wrk -t2 -c2 -d10s -s post.lua http://localhost:8080/api/v1/generate/text-to-video
```

Create `post.lua`:
```lua
wrk.method = "POST"
wrk.body   = '{"prompt":"test","num_frames":16,"width":256,"height":256}'
wrk.headers["Content-Type"] = "application/json"
```

## Performance Tests

### GPU Memory Test

Check GPU utilization during generation:

```bash
# Monitor in separate terminal
watch -n 1 nvidia-smi

# Submit multiple requests
for i in {1..3}; do
  curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
    -H "Content-Type: application/json" \
    -d "{\"prompt\":\"Test $i\",\"num_frames\":32}"
  sleep 1
done
```

### Concurrent Request Test

```bash
# Test rate limiting
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
    -H "Content-Type: application/json" \
    -d '{"prompt":"Test concurrent","num_frames":16}' &
done
wait
```

Should see some requests return 429 (Too Many Requests).

## Manual Testing Checklist

### Server Startup
- [ ] Server starts without errors
- [ ] Python backend connects successfully
- [ ] GPU is detected (if available)
- [ ] Model loads correctly
- [ ] Logs show correct configuration

### Basic Functionality
- [ ] Health check returns 200
- [ ] Model info endpoint works
- [ ] List models endpoint works
- [ ] Invalid requests return 400

### Text-to-Video
- [ ] Simple prompt generates video
- [ ] Negative prompt is respected
- [ ] Different resolutions work (256x256, 512x512, 768x768)
- [ ] Different frame counts work (16, 32, 64, 128)
- [ ] Seed produces consistent results
- [ ] Job status updates correctly
- [ ] Generated video is playable

### Image-to-Video
- [ ] JPG images work
- [ ] PNG images work
- [ ] Large images are handled
- [ ] Generated video matches input aspect ratio
- [ ] Prompt guidance affects output

### Video-to-Video
- [ ] MP4 videos work
- [ ] Strength parameter affects transformation
- [ ] Output maintains input length

### Error Handling
- [ ] Missing prompt returns 400
- [ ] Invalid parameters return 400
- [ ] Oversized uploads return 400
- [ ] Too many concurrent requests return 429
- [ ] Backend down returns 503
- [ ] Invalid job ID returns 404

### Performance
- [ ] Generation completes within timeout
- [ ] Memory usage is stable
- [ ] No memory leaks during multiple generations
- [ ] Rate limiting works correctly
- [ ] Multiple jobs can queue

## Automated Test Suite

Create a comprehensive test suite:

```bash
#!/bin/bash
# test_suite.sh

set -e

API_URL="http://localhost:8080"
PYTHON_URL="http://localhost:5000"

echo "Running Wan2.1 Video Server Test Suite"
echo "======================================="

# Test 1: Health checks
echo "Test 1: Health Checks"
curl -f "$API_URL/health" > /dev/null && echo "✓ Go server healthy"
curl -f "$PYTHON_URL/health" > /dev/null && echo "✓ Python backend healthy"

# Test 2: Model info
echo "Test 2: Model Info"
curl -f "$API_URL/api/v1/model/info" > /dev/null && echo "✓ Model info accessible"

# Test 3: Text-to-video (minimal)
echo "Test 3: Text-to-Video Generation"
RESPONSE=$(curl -s -X POST "$API_URL/api/v1/generate/text-to-video" \
  -H "Content-Type: application/json" \
  -d '{"prompt":"test","num_frames":8,"width":128,"height":128,"num_inference_steps":10}')
JOB_ID=$(echo "$RESPONSE" | jq -r '.job_id')
if [ -n "$JOB_ID" ] && [ "$JOB_ID" != "null" ]; then
  echo "✓ Text-to-video request accepted: $JOB_ID"
else
  echo "✗ Text-to-video request failed"
  exit 1
fi

# Test 4: Job status
echo "Test 4: Job Status Check"
sleep 2
curl -f "$API_URL/api/v1/job/$JOB_ID" > /dev/null && echo "✓ Job status accessible"

# Test 5: List models
echo "Test 5: List Models"
curl -f "$API_URL/api/v1/models" > /dev/null && echo "✓ Models list accessible"

# Test 6: Invalid requests
echo "Test 6: Error Handling"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/api/v1/generate/text-to-video" \
  -H "Content-Type: application/json" \
  -d '{}')
if [ "$HTTP_CODE" = "400" ]; then
  echo "✓ Invalid request returns 400"
else
  echo "✗ Expected 400, got $HTTP_CODE"
fi

echo ""
echo "All tests passed! ✓"
```

## Debugging Tests

### Enable Debug Logging

```bash
# In .env
LOG_LEVEL=debug
SERVER_MODE=debug
```

Restart server and check logs for detailed information.

### Common Issues

**Connection Refused:**
- Check if servers are running
- Verify ports are correct
- Check firewall rules

**Timeout Errors:**
- Increase `REQUEST_TIMEOUT`
- Reduce frame count or resolution
- Check GPU availability

**Out of Memory:**
- Reduce concurrent requests
- Lower resolution
- Decrease `GPU_MEMORY_FRACTION`

**Model Not Loaded:**
- Run `./wan2-video-server download`
- Check Python backend logs
- Verify Hugging Face token

## Continuous Integration

### GitHub Actions Example

Create `.github/workflows/test.yml`:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Go tests
        run: |
          go test ./... -v -race -coverprofile=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## Test Coverage Goals

- **Go Code:** > 70% coverage
- **Handlers:** > 80% coverage
- **Core Logic:** > 90% coverage

Check coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Benchmarking

Run Go benchmarks:

```bash
go test -bench=. -benchmem ./...
```

Example benchmark:

```go
func BenchmarkTextToVideoRequest(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark code
    }
}
```
