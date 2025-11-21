# Wan2.1 Video Server - Project Overview

## 📋 Project Summary

A production-ready, high-performance video generation server built in Go that serves the LTX-Video (Wan2.1) model from Hugging Face. The application supports text-to-video, image-to-video, and video-to-video generation with GPU acceleration.

## 🏗️ Architecture

### Two-Tier Design

```
┌──────────────────────────────────────────┐
│         Client Applications              │
│     (Browser, Mobile, API Clients)       │
└──────────────┬───────────────────────────┘
               │ HTTP/REST
               ▼
┌──────────────────────────────────────────┐
│          Go Server (Port 8080)           │
│  ┌────────────────────────────────────┐  │
│  │  • HTTP API (Gin Framework)       │  │
│  │  • Request Validation             │  │
│  │  • File Upload Handling           │  │
│  │  • Job Queue Management           │  │
│  │  • Rate Limiting                  │  │
│  │  • Logging & Monitoring           │  │
│  └────────────────────────────────────┘  │
└──────────────┬───────────────────────────┘
               │ HTTP
               ▼
┌──────────────────────────────────────────┐
│      Python Backend (Port 5000)          │
│  ┌────────────────────────────────────┐  │
│  │  • PyTorch Model Loading          │  │
│  │  • GPU/CUDA Management            │  │
│  │  • Diffusers Pipeline             │  │
│  │  • Video Generation               │  │
│  │  • Frame Processing               │  │
│  └────────────────────────────────────┘  │
└──────────────┬───────────────────────────┘
               │
               ▼
         ┌─────────────┐
         │  NVIDIA GPU │
         │   (CUDA)    │
         └─────────────┘
```

## 📁 Project Structure

```
wan2-video-server/
├── cmd/                           # CLI commands
│   ├── root.go                   # Main server command
│   └── download.go               # Model download command
│
├── pkg/                          # Go packages
│   ├── config/                   # Configuration management
│   │   └── config.go            # Config loader with viper
│   ├── server/                   # HTTP server
│   │   └── server.go            # Gin server setup
│   ├── handlers/                 # HTTP request handlers
│   │   ├── health.go            # Health checks
│   │   ├── video.go             # Video generation endpoints
│   │   ├── model.go             # Model info
│   │   └── model_management.go  # Model operations
│   ├── middleware/               # HTTP middleware
│   │   └── middleware.go        # CORS, logging, rate limiting
│   ├── model/                    # Model engines
│   │   ├── engine.go            # Engine interface
│   │   ├── python_engine.go     # Python backend client
│   │   ├── local_engine.go      # Local inference (future)
│   │   └── huggingface.go       # HF model downloader
│   ├── types/                    # Type definitions
│   │   └── types.go             # Request/response structs
│   ├── logger/                   # Logging
│   │   └── logger.go            # Logrus setup
│   └── utils/                    # Utilities
│       └── utils.go             # Helper functions
│
├── python_backend/               # Python inference server
│   ├── server.py                # Flask API server
│   ├── requirements.txt         # Python dependencies
│   └── setup.sh                 # Setup script
│
├── docs/                         # Documentation
│   ├── API.md                   # API reference
│   ├── CONFIGURATION.md         # Config guide
│   ├── TESTING.md               # Testing guide
│   └── DEPLOYMENT.md            # Deployment guide
│
├── examples/                     # Example requests
│   ├── api_examples.sh          # Shell script examples
│   └── postman_collection.json  # Postman collection
│
├── main.go                       # Application entry point
├── go.mod                        # Go dependencies
├── .env.example                  # Example configuration
├── Dockerfile                    # Docker image definition
├── docker-compose.yml            # Docker Compose config
├── Makefile                      # Build automation
├── setup.sh                      # Quick setup script
├── README.md                     # Main documentation
├── QUICKSTART.md                 # Quick start guide
└── LICENSE                       # MIT license
```

## 🎯 Key Features

### Video Generation Capabilities
- ✅ **Text-to-Video**: Generate videos from text descriptions
- ✅ **Image-to-Video**: Animate static images
- ✅ **Video-to-Video**: Transform existing videos

### Technical Features
- ✅ **GPU Acceleration**: CUDA support for faster inference
- ✅ **Async Processing**: Non-blocking job queue
- ✅ **Rate Limiting**: Prevent resource exhaustion
- ✅ **Model Caching**: Efficient model storage
- ✅ **File Management**: Upload and output handling
- ✅ **Health Monitoring**: Ready for production monitoring
- ✅ **CORS Support**: Cross-origin requests enabled
- ✅ **Structured Logging**: JSON logs for easy parsing

### Model Features
- 🔧 Adjustable frame count (up to 128 frames)
- 🔧 Configurable resolution (256x256 to 1024x1024+)
- 🔧 Custom FPS settings
- 🔧 Seed support for reproducibility
- 🔧 Guidance scale tuning
- 🔧 Inference step control
- 🔧 Negative prompt support

## 🔌 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/api/v1/model/info` | GET | Model information |
| `/api/v1/generate/text-to-video` | POST | Text-to-video generation |
| `/api/v1/generate/image-to-video` | POST | Image-to-video generation |
| `/api/v1/generate/video-to-video` | POST | Video-to-video generation |
| `/api/v1/job/:id` | GET | Job status |
| `/api/v1/models` | GET | List models |
| `/api/v1/models/download` | POST | Download model |
| `/outputs/*` | GET | Static file serving |

## 🛠️ Technology Stack

### Backend (Go)
- **Web Framework**: Gin v1.9.1
- **Configuration**: Viper v1.18.2
- **CLI**: Cobra v1.8.0
- **Logging**: Logrus v1.9.3
- **Environment**: godotenv v1.5.1

### AI Backend (Python)
- **Web Framework**: Flask 3.0.0
- **ML Framework**: PyTorch 2.0+
- **Model Pipeline**: Diffusers 0.25+
- **Transformers**: Transformers 4.35+
- **Acceleration**: Accelerate 0.25+
- **Optimization**: xformers 0.0.22 (optional)

### Infrastructure
- **Containerization**: Docker, Docker Compose
- **GPU**: NVIDIA CUDA 11.8+
- **Reverse Proxy**: Nginx (optional)
- **Monitoring**: Prometheus, Grafana (optional)

## 📊 Performance Characteristics

### Generation Times (Approximate)

| Configuration | Resolution | Frames | GPU (T4) | GPU (A100) |
|--------------|------------|--------|----------|------------|
| Fast | 256x256 | 16 | ~10s | ~3s |
| Balanced | 512x512 | 64 | ~45s | ~15s |
| High Quality | 768x768 | 128 | ~120s | ~40s |

### Resource Requirements

| Scenario | GPU Memory | RAM | Concurrent Jobs |
|----------|------------|-----|-----------------|
| Minimal | 8GB | 16GB | 1 |
| Recommended | 16GB | 32GB | 2-3 |
| High Load | 24GB+ | 64GB | 4+ |

## 🚀 Quick Start Commands

```bash
# Setup
./setup.sh

# Download model
./wan2-video-server download

# Start services
# Terminal 1:
cd python_backend && source venv/bin/activate && python server.py

# Terminal 2:
./wan2-video-server

# Test
curl http://localhost:8080/health

# Generate video
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{"prompt": "A cat playing", "num_frames": 32}'
```

## 🔐 Security Considerations

- ✅ Input validation on all endpoints
- ✅ File size limits (100MB default)
- ✅ Rate limiting (2 concurrent requests default)
- ✅ Request timeouts (300s default)
- ⚠️ No authentication (add in production)
- ⚠️ No HTTPS (use reverse proxy)
- ⚠️ No API keys (implement as needed)

## 📈 Scalability

### Horizontal Scaling
- Deploy multiple instances behind load balancer
- Use shared storage for models and outputs
- Implement Redis for distributed job queue

### Vertical Scaling
- Add more GPUs to single instance
- Increase concurrent request limit
- Allocate more memory

## 🧪 Testing

```bash
# Unit tests
go test ./...

# Integration tests
./examples/api_examples.sh

# Load testing
make test-load
```

## 📝 Documentation

| Document | Purpose |
|----------|---------|
| `README.md` | Main documentation and overview |
| `QUICKSTART.md` | Get started in 5 minutes |
| `docs/API.md` | Complete API reference |
| `docs/CONFIGURATION.md` | All configuration options |
| `docs/TESTING.md` | Testing guide |
| `docs/DEPLOYMENT.md` | Production deployment |

## 🔄 Workflow

1. **Client** sends generation request to Go server
2. **Go Server** validates request and creates job
3. **Go Server** forwards to Python backend via HTTP
4. **Python Backend** loads model (if not cached)
5. **Python Backend** generates video using GPU
6. **Python Backend** saves video and returns path
7. **Go Server** updates job status
8. **Client** polls for status and downloads result

## 🎨 Use Cases

- **Content Creation**: Automated video content generation
- **Marketing**: Product visualization and ads
- **Education**: Instructional video creation
- **Entertainment**: Story visualization
- **Prototyping**: Rapid video concept testing
- **Research**: AI/ML experimentation

## 🌟 Future Enhancements

- [ ] Web UI dashboard
- [ ] Batch processing
- [ ] Video upscaling
- [ ] Frame interpolation
- [ ] Multiple model support
- [ ] Ollama integration
- [ ] Authentication system
- [ ] Cloud storage (S3, GCS)
- [ ] WebSocket streaming
- [ ] Job queue with Redis
- [ ] Kubernetes deployment
- [ ] API rate limiting per user
- [ ] Video editing features
- [ ] Custom model fine-tuning

## 👥 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

MIT License - see `LICENSE` file

## 🙏 Acknowledgments

- **Lightricks** - LTX-Video model
- **Hugging Face** - Model hosting and Diffusers library
- **LocalAI** - Inspiration for architecture
- **Go Community** - Excellent libraries and tools
- **PyTorch Team** - Deep learning framework

## 📞 Support

- 📖 Documentation: See `docs/` folder
- 🐛 Bug Reports: GitHub Issues
- 💡 Feature Requests: GitHub Issues
- 📧 Email: support@example.com

---

**Built with ❤️ for the AI video generation community**
