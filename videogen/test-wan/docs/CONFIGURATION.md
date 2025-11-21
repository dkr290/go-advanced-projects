# Configuration Guide

This guide explains all configuration options for the Wan2.1 Video Server.

## Configuration Methods

Configuration can be provided through:

1. **Environment variables** (in `.env` file or system environment)
2. **Command-line flags**
3. **Default values**

Priority: Command-line flags > Environment variables > Defaults

## Configuration File

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

## Server Configuration

### SERVER_HOST
- **Type:** String
- **Default:** `0.0.0.0`
- **Description:** Server bind address. Use `0.0.0.0` for all interfaces, `127.0.0.1` for localhost only.
- **Example:** `SERVER_HOST=0.0.0.0`

### SERVER_PORT
- **Type:** Integer
- **Default:** `8080`
- **Description:** HTTP server port.
- **Example:** `SERVER_PORT=8080`

### SERVER_MODE
- **Type:** String
- **Default:** `release`
- **Options:** `release`, `debug`
- **Description:** Gin framework mode. Use `debug` for development.
- **Example:** `SERVER_MODE=release`

## Model Configuration

### MODEL_NAME
- **Type:** String
- **Default:** `Wan2.1`
- **Description:** Display name for the model.
- **Example:** `MODEL_NAME=Wan2.1`

### HUGGINGFACE_MODEL_ID
- **Type:** String
- **Default:** `Lightricks/LTX-Video`
- **Description:** Hugging Face model repository ID.
- **Example:** `HUGGINGFACE_MODEL_ID=Lightricks/LTX-Video`

### MODEL_CACHE_DIR
- **Type:** String
- **Default:** `./models`
- **Description:** Directory to cache downloaded models.
- **Example:** `MODEL_CACHE_DIR=/data/models`

### USE_HUGGINGFACE
- **Type:** Boolean
- **Default:** `true`
- **Description:** Enable Hugging Face model provider.
- **Example:** `USE_HUGGINGFACE=true`

### USE_OLLAMA
- **Type:** Boolean
- **Default:** `false`
- **Description:** Enable Ollama model provider (future feature).
- **Example:** `USE_OLLAMA=false`

### MAX_FRAMES
- **Type:** Integer
- **Default:** `128`
- **Description:** Maximum number of frames allowed per video.
- **Example:** `MAX_FRAMES=128`

### DEFAULT_FPS
- **Type:** Integer
- **Default:** `24`
- **Description:** Default frames per second for generated videos.
- **Example:** `DEFAULT_FPS=24`

### DEFAULT_WIDTH
- **Type:** Integer
- **Default:** `512`
- **Description:** Default video width in pixels.
- **Example:** `DEFAULT_WIDTH=512`

### DEFAULT_HEIGHT
- **Type:** Integer
- **Default:** `512`
- **Description:** Default video height in pixels.
- **Example:** `DEFAULT_HEIGHT=512`

### MAX_VIDEO_DURATION
- **Type:** Integer
- **Default:** `10`
- **Description:** Maximum video duration in seconds.
- **Example:** `MAX_VIDEO_DURATION=10`

## GPU Configuration

### ENABLE_GPU
- **Type:** Boolean
- **Default:** `true`
- **Description:** Enable GPU acceleration. Set to `false` for CPU-only mode.
- **Example:** `ENABLE_GPU=true`

### GPU_DEVICE_ID
- **Type:** Integer
- **Default:** `0`
- **Description:** CUDA device ID to use (for multi-GPU systems).
- **Example:** `GPU_DEVICE_ID=0`

### GPU_MEMORY_FRACTION
- **Type:** Float
- **Default:** `0.9`
- **Description:** Fraction of GPU memory to use (0.0 - 1.0).
- **Example:** `GPU_MEMORY_FRACTION=0.9`

## Hugging Face Configuration

### HUGGINGFACE_TOKEN
- **Type:** String
- **Default:** `your_hf_token_here`
- **Description:** Hugging Face API token for private models or faster downloads.
- **How to get:** Visit https://huggingface.co/settings/tokens
- **Example:** `HUGGINGFACE_TOKEN=hf_xxxxxxxxxxxxxxxxxxxxx`

### HUGGINGFACE_API_URL
- **Type:** String
- **Default:** `https://huggingface.co`
- **Description:** Hugging Face API base URL (rarely needs changing).
- **Example:** `HUGGINGFACE_API_URL=https://huggingface.co`

## Python Backend Configuration

### PYTHON_BACKEND_URL
- **Type:** String
- **Default:** `http://localhost:5000`
- **Description:** URL of the Python inference backend.
- **Example:** `PYTHON_BACKEND_URL=http://localhost:5000`

### PYTHON_BACKEND_ENABLED
- **Type:** Boolean
- **Default:** `true`
- **Description:** Enable Python backend for inference.
- **Example:** `PYTHON_BACKEND_ENABLED=true`

## Processing Configuration

### MAX_CONCURRENT_REQUESTS
- **Type:** Integer
- **Default:** `2`
- **Description:** Maximum number of concurrent video generation requests.
- **Recommendation:** Adjust based on GPU memory (1-2 for 8GB GPU, 2-4 for 16GB+).
- **Example:** `MAX_CONCURRENT_REQUESTS=2`

### REQUEST_TIMEOUT
- **Type:** Integer (seconds)
- **Default:** `300`
- **Description:** Timeout for video generation requests.
- **Example:** `REQUEST_TIMEOUT=300`

### UPLOAD_MAX_SIZE
- **Type:** String
- **Default:** `100MB`
- **Description:** Maximum upload file size.
- **Format:** Number + Unit (KB, MB, GB)
- **Example:** `UPLOAD_MAX_SIZE=200MB`

## Logging Configuration

### LOG_LEVEL
- **Type:** String
- **Default:** `info`
- **Options:** `debug`, `info`, `warn`, `error`
- **Description:** Logging verbosity level.
- **Example:** `LOG_LEVEL=info`

### LOG_FORMAT
- **Type:** String
- **Default:** `json`
- **Options:** `json`, `text`
- **Description:** Log output format.
- **Example:** `LOG_FORMAT=json`

## Command-Line Flags

Override configuration with command-line flags:

```bash
./wan2-video-server \
  --host 0.0.0.0 \
  --port 8080 \
  --gpu true \
  --log-level debug
```

Available flags:
- `--config <file>` - Configuration file path
- `--host <address>` - Server host
- `--port <number>` - Server port
- `--gpu <bool>` - Enable GPU
- `--log-level <level>` - Log level

## Environment-Specific Configurations

### Development

```env
SERVER_MODE=debug
LOG_LEVEL=debug
LOG_FORMAT=text
MAX_CONCURRENT_REQUESTS=1
```

### Production

```env
SERVER_MODE=release
LOG_LEVEL=info
LOG_FORMAT=json
MAX_CONCURRENT_REQUESTS=2
REQUEST_TIMEOUT=600
```

### High-Performance (Multi-GPU)

```env
ENABLE_GPU=true
GPU_DEVICE_ID=0
GPU_MEMORY_FRACTION=0.95
MAX_CONCURRENT_REQUESTS=4
```

### CPU-Only

```env
ENABLE_GPU=false
MAX_CONCURRENT_REQUESTS=1
REQUEST_TIMEOUT=1800
```

## Performance Tuning

### For Better Quality
```env
# Default settings already optimized for quality
DEFAULT_WIDTH=768
DEFAULT_HEIGHT=768
MAX_FRAMES=128
```

### For Faster Generation
```env
DEFAULT_WIDTH=256
DEFAULT_HEIGHT=256
MAX_FRAMES=32
MAX_CONCURRENT_REQUESTS=3
```

### For Memory-Constrained Systems
```env
DEFAULT_WIDTH=384
DEFAULT_HEIGHT=384
MAX_FRAMES=48
GPU_MEMORY_FRACTION=0.8
MAX_CONCURRENT_REQUESTS=1
```

## Security Considerations

### Production Deployment

1. **Restrict network access:**
   ```env
   SERVER_HOST=127.0.0.1  # Localhost only
   ```

2. **Use reverse proxy** (nginx, traefik) for:
   - SSL/TLS termination
   - Authentication
   - Rate limiting

3. **Limit upload sizes:**
   ```env
   UPLOAD_MAX_SIZE=50MB
   ```

4. **Set request timeouts:**
   ```env
   REQUEST_TIMEOUT=300
   ```

## Troubleshooting

### High Memory Usage
- Reduce `MAX_CONCURRENT_REQUESTS`
- Lower `GPU_MEMORY_FRACTION`
- Decrease default resolution

### Slow Performance
- Enable GPU: `ENABLE_GPU=true`
- Increase `GPU_MEMORY_FRACTION`
- Reduce `REQUEST_TIMEOUT` to fail faster

### Connection Issues
- Check `PYTHON_BACKEND_URL` is correct
- Ensure Python backend is running
- Verify firewall rules

## Configuration Validation

The server validates configuration on startup. Check logs for:

```
INFO Using config file: .env
INFO Configuration loaded: GPU=true, Backend=huggingface
```

Invalid configurations will show errors:

```
ERROR Failed to load configuration: invalid port number
```

## Example Configurations

See `examples/` directory for sample configurations:
- `development.env` - Development setup
- `production.env` - Production setup
- `gpu-optimized.env` - Multi-GPU setup
