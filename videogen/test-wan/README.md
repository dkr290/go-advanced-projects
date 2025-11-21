# Wan2.1 Video Generation Server

A high-performance video generation server built in Go for serving the LTX-Video (Wan2.1) model with support for text-to-video, image-to-video, and video-to-video generation.

## Features

- 🎥 **Multiple Generation Modes**
  - Text-to-video generation
  - Image-to-video generation
  - Video-to-video generation

- 🚀 **High Performance**
  - GPU acceleration support (CUDA)
  - Concurrent request handling
  - Asynchronous job processing
  - Request rate limiting

- 🤗 **Model Management**
  - Hugging Face model integration
  - Automatic model downloading
  - Model caching

- 🔧 **Flexible Configuration**
  - Environment-based configuration
  - Command-line flags
  - Customizable model parameters

- 📊 **Production Ready**
  - RESTful API
  - Job status tracking
  - Health check endpoints
  - Structured logging

## Architecture

The application consists of two main components:

1. **Go Server** - HTTP API server that handles requests, file uploads, and job management
2. **Python Backend** - GPU-accelerated model inference using PyTorch and Diffusers

```
┌─────────────────┐
│   Client/API    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Go Server     │  - Request handling
│   (Port 8080)   │  - File management
│                 │  - Job tracking
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Python Backend  │  - Model inference
│   (Port 5000)   │  - GPU processing
│                 │  - Video generation
└─────────────────┘
```

## Prerequisites

- Go 1.21 or higher
- Python 3.9 or higher
- **GPU with 12GB+ VRAM** (16GB+ recommended) OR use CPU mode
  - ⚠️ **4GB GPUs will NOT work** with LTX-Video - see [Low Memory Guide](docs/LOW_MEMORY_GUIDE.md)
  - For 4GB-8GB GPUs, consider using ModelScope or CPU mode
- 16GB+ RAM (32GB+ recommended)
- 20GB+ free disk space for models

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/wan2-video-server
cd wan2-video-server
```

### 2. Setup Go Server

```bash
# Install Go dependencies
go mod download

# Build the application
go build -o wan2-video-server main.go
```

### 3. Setup Python Backend

```bash
cd python_backend

# Make setup script executable
chmod +x setup.sh

# Run setup script
./setup.sh

# Or manually:
python3 -m venv venv
source venv/bin/activate
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118
pip install -r requirements.txt
```

### 4. Configure Environment

```bash
# Copy example configuration
cp .env.example .env

# Edit configuration
nano .env
```

Key configuration options:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Model
HUGGINGFACE_MODEL_ID=Lightricks/LTX-Video
MODEL_CACHE_DIR=./models

# GPU
ENABLE_GPU=true
GPU_DEVICE_ID=0

# Hugging Face (optional, for private models)
HUGGINGFACE_TOKEN=your_token_here
```

## Usage

### Starting the Server

#### Option 1: Start Both Services

**Terminal 1 - Python Backend:**
```bash
cd python_backend
source venv/bin/activate
python server.py
```

**Terminal 2 - Go Server:**
```bash
./wan2-video-server
# or
go run main.go
```

#### Option 2: Using Docker Compose (Future)

```bash
docker-compose up
```

### Download Model

Before generating videos, download the model:

```bash
./wan2-video-server download --model-id Lightricks/LTX-Video
```

### API Endpoints

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Model Information
```bash
curl http://localhost:8080/api/v1/model/info
```

#### Text-to-Video Generation
```bash
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A cat playing with a ball of yarn",
    "num_frames": 64,
    "fps": 24,
    "width": 512,
    "height": 512,
    "guidance_scale": 7.5,
    "num_inference_steps": 50
  }'
```

Response:
```json
{
  "job_id": "job_1234567890_abcdef",
  "status": "processing",
  "message": "Video generation started"
}
```

#### Check Job Status
```bash
curl http://localhost:8080/api/v1/job/{job_id}
```

#### Image-to-Video Generation
```bash
curl -X POST http://localhost:8080/api/v1/generate/image-to-video \
  -F "image=@/path/to/image.jpg" \
  -F "prompt=A beautiful sunset over the ocean" \
  -F "num_frames=64" \
  -F "fps=24"
```

#### Video-to-Video Generation
```bash
curl -X POST http://localhost:8080/api/v1/generate/video-to-video \
  -F "video=@/path/to/video.mp4" \
  -F "prompt=Transform into an anime style" \
  -F "strength=0.8"
```

## API Documentation

### Text-to-Video Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `prompt` | string | required | Text description of the video |
| `negative_prompt` | string | "" | What to avoid in generation |
| `num_frames` | int | 64 | Number of frames to generate |
| `fps` | int | 24 | Frames per second |
| `width` | int | 512 | Video width |
| `height` | int | 512 | Video height |
| `seed` | int | -1 | Random seed (-1 for random) |
| `guidance_scale` | float | 7.5 | How closely to follow prompt |
| `num_inference_steps` | int | 50 | Generation quality/speed trade-off |

### Image-to-Video Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `image` | file | required | Input image file |
| `prompt` | string | "" | Additional guidance |
| `num_frames` | int | 64 | Number of frames |
| `fps` | int | 24 | Frames per second |
| ... | ... | ... | (same as text-to-video) |

## Project Structure

```
wan2-video-server/
├── cmd/                      # Command-line interface
│   ├── root.go              # Main command
│   └── download.go          # Model download command
├── pkg/                      # Application packages
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP handlers
│   ├── logger/              # Logging utilities
│   ├── middleware/          # HTTP middleware
│   ├── model/               # Model interfaces and engines
│   ├── server/              # HTTP server
│   ├── types/               # Type definitions
│   └── utils/               # Utility functions
├── python_backend/          # Python inference backend
│   ├── server.py           # Flask server
│   ├── requirements.txt    # Python dependencies
│   └── setup.sh            # Setup script
├── main.go                  # Application entry point
├── go.mod                   # Go dependencies
├── .env.example            # Example configuration
└── README.md               # This file
```

## Performance Tips

### GPU Optimization

1. **Enable xformers** for memory efficiency:
   ```bash
   pip install xformers
   ```

2. **Adjust GPU memory fraction** in `.env`:
   ```env
   GPU_MEMORY_FRACTION=0.9
   ```

3. **Use mixed precision** (automatically enabled on CUDA)

### Generation Tips

1. **Faster generation**: Reduce `num_inference_steps` (e.g., 25-30)
2. **Better quality**: Increase `num_inference_steps` (e.g., 50-100)
3. **Consistency**: Use the same `seed` value
4. **Memory usage**: Reduce `width`, `height`, or `num_frames`

## Troubleshooting

### Common Issues

**"Model not loaded" error:**
- Ensure Python backend is running
- Check Python backend logs for errors
- Verify model was downloaded successfully

**Out of memory errors:**
- **CRITICAL: LTX-Video needs 12GB+ VRAM minimum**
- For 4GB GPUs: See [Low Memory Guide](docs/LOW_MEMORY_GUIDE.md)
- Switch to CPU mode: `ENABLE_GPU=false` in `.env`
- Use smaller model like ModelScope: `HUGGINGFACE_MODEL_ID=damo-vilab/text-to-video-ms-1.7b`
- Reduce batch size or frame count
- Lower resolution (width/height)
- Enable gradient checkpointing
- Use cloud GPU (Google Colab is free!)

**Slow generation:**
- Ensure GPU is being used (check logs)
- Install xformers
- Reduce inference steps
- Check GPU utilization: `nvidia-smi`

**Connection refused:**
- Verify Python backend is running on port 5000
- Check firewall settings
- Ensure PYTHON_BACKEND_URL is correct

**AMD GPU issues:**
- See [AMD GPU Guide](docs/AMD_GPU_GUIDE.md)
- Linux only (ROCm not available on Windows)
- Run `./setup_rocm.sh` for automated setup
- Or use CPU mode: `ENABLE_GPU=false`

**GPU Compatibility:**
- Check `GPU_COMPATIBILITY.txt` for your GPU
- NVIDIA: Native CUDA support ✅
- AMD: ROCm on Linux only ⚠️
- Intel/Apple Silicon: CPU mode recommended
- 4GB VRAM: Not sufficient for LTX-Video

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o wan2-video-server main.go

# Or with version info
VERSION=1.0.0
go build -ldflags="-s -w -X main.Version=${VERSION}" -o wan2-video-server main.go
```

### Adding Custom Models

1. Update model ID in `.env`
2. Implement custom pipeline in `python_backend/server.py`
3. Adjust parameters in `pkg/config/config.go`

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Lightricks](https://www.lightricks.com/) for the LTX-Video model
- [Hugging Face](https://huggingface.co/) for model hosting and diffusers library
- [LocalAI](https://github.com/mudler/LocalAI) for inspiration

## Support

- 📧 Email: support@example.com
- 💬 Discord: [Join our community](https://discord.gg/example)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/wan2-video-server/issues)

## Roadmap

- [ ] Docker and Docker Compose support
- [ ] Ollama integration
- [ ] Web UI dashboard
- [ ] Batch processing
- [ ] Video upscaling
- [ ] Frame interpolation
- [ ] Multiple model support
- [ ] Cloud storage integration (S3, GCS)
- [ ] Authentication and API keys
- [ ] Usage analytics
