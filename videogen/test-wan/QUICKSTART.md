# Quick Start Guide

Get up and running with Wan2.1 Video Server in 5 minutes!

## Prerequisites

- Go 1.21+
- Python 3.9+
- CUDA-capable GPU (optional but recommended)
- 16GB+ RAM

## Installation (5 steps)

### 1. Clone & Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd wan2-video-server

# Run automated setup
chmod +x setup.sh
./setup.sh
```

### 2. Configure

```bash
# Edit configuration (optional)
nano .env

# Key settings:
# - ENABLE_GPU=true (if you have GPU)
# - HUGGINGFACE_TOKEN=your_token (for private models)
```

### 3. Download Model

```bash
./wan2-video-server download
```

This downloads the LTX-Video model (~10GB). Takes 5-10 minutes depending on your internet speed.

### 4. Start Python Backend

```bash
cd python_backend
source venv/bin/activate
python server.py
```

Keep this terminal running.

### 5. Start Go Server

In a new terminal:

```bash
./wan2-video-server
```

Server will start on http://localhost:8080

## First Video Generation

### Using curl

```bash
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A cat playing with a ball of yarn",
    "num_frames": 64,
    "fps": 24,
    "width": 512,
    "height": 512
  }'
```

Response:
```json
{
  "job_id": "job_1234567890_abc",
  "status": "processing"
}
```

### Check Status

```bash
curl http://localhost:8080/api/v1/job/job_1234567890_abc
```

### Download Video

Once completed, video is at:
```
http://localhost:8080/outputs/<filename>
```

## Quick Examples

### Fast Generation (for testing)

```bash
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A spinning cube",
    "num_frames": 16,
    "width": 256,
    "height": 256,
    "num_inference_steps": 20
  }'
```

### High Quality Generation

```bash
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A beautiful sunset over mountains",
    "num_frames": 128,
    "width": 768,
    "height": 768,
    "num_inference_steps": 100,
    "guidance_scale": 9.0
  }'
```

### Image-to-Video

```bash
curl -X POST http://localhost:8080/api/v1/generate/image-to-video \
  -F "image=@my_photo.jpg" \
  -F "prompt=Make this photo come alive" \
  -F "num_frames=64"
```

## Common Parameters

| Parameter | Fast | Balanced | High Quality |
|-----------|------|----------|--------------|
| num_frames | 16 | 64 | 128 |
| width/height | 256 | 512 | 768 |
| num_inference_steps | 20 | 50 | 100 |
| guidance_scale | 5.0 | 7.5 | 9.0 |

## Troubleshooting

### "Connection refused"
- Ensure Python backend is running
- Check port 5000 is not in use

### "Model not loaded"
- Run: `./wan2-video-server download`
- Check Python backend logs

### "Out of memory"
- Reduce resolution: `width: 256, height: 256`
- Reduce frames: `num_frames: 32`
- Lower concurrent requests in `.env`: `MAX_CONCURRENT_REQUESTS=1`

### Slow generation
- Enable GPU: `ENABLE_GPU=true` in `.env`
- Check GPU is being used: `nvidia-smi`
- Reduce quality settings

## Next Steps

1. **Read full documentation:** `README.md`
2. **Explore API:** `docs/API.md`
3. **Tune configuration:** `docs/CONFIGURATION.md`
4. **Run tests:** `docs/TESTING.md`
5. **Try examples:** `./examples/api_examples.sh`

## Useful Commands

```bash
# Build
make build

# Run tests
make test

# Clean outputs
make clean

# View logs (Go server)
./wan2-video-server --log-level debug

# View logs (Python backend)
cd python_backend && python server.py

# Check GPU usage
watch -n 1 nvidia-smi

# Monitor outputs
ls -lh outputs/
```

## Tips

1. **Use consistent seeds** for reproducible results
2. **Start with low settings** to test, then increase quality
3. **Monitor GPU memory** with `nvidia-smi`
4. **Use negative prompts** to avoid unwanted elements
5. **Experiment with guidance_scale** (7.5 is good default)

## Need Help?

- Check logs for errors
- Read `docs/API.md` for detailed API reference
- See `docs/CONFIGURATION.md` for all settings
- Create an issue on GitHub

Happy video generating! 🎥✨
