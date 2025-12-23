# 🚀 Quick Reference Card - FLUX Go Backend with Web UI

## Installation & Setup

```bash
# 1. Make sure you have Go and Python3 installed
go version
python3 --version

# 2. Install Python dependencies
pip3 install torch diffusers pillow

# 3. Build the project
go build -o flux-gen
```

## Command Cheat Sheet

### Text-to-Image (Basic)
```bash
go run main.go -config prompts.json
```

### Text-to-Image + Web UI
```bash
go run main.go -config prompts.json --web
# Open: http://localhost:8080
```

### Image-to-Image
```bash
# 1. Place input images in ./images/
# 2. Run:
go run main.go -config prompts.json --img2img --strength 0.75
```

### Image-to-Image + Web UI
```bash
go run main.go -config prompts.json --img2img --web
```

### Web Server Only (Browse Existing)
```bash
go run main.go -config empty.json --web --output ./your_images
```

### Custom Port
```bash
go run main.go -config prompts.json --web --web-port 9000
# Open: http://localhost:9000
```

## All Available Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config` | string | *required* | Path to JSON config file |
| `--gguf-model-url` | string | "" | URL to download GGUF model |
| `--hf-model` | string | FLUX.1-dev | HuggingFace model ID |
| `--lora-url` | string | "" | URL to LoRA safetensors file |
| `--seed` | int | 42 | Random seed |
| `--output` | string | ./output | Output directory |
| `--resolution` | string | 1024x1024 | Image resolution |
| `--num_steps` | int | 28 | Inference steps |
| `--guidence_scale` | float | 7.0 | Guidance scale |
| `--model-down-path` | string | ./models | Model download path |
| `--lora-down-path` | string | ./models/lora | LoRA download path |
| `--low-vram` | bool | false | Enable CPU offload |
| `--debug` | bool | false | Debug mode |
| **`--img2img`** | **bool** | **false** | **Image-to-image mode** |
| **`--strength`** | **float** | **0.75** | **Img2img strength (0.0-1.0)** |
| **`--web`** | **bool** | **false** | **Enable web server** |
| **`--web-port`** | **int** | **8080** | **Web server port** |

## Environment Variables

```bash
export MODEL_URL="https://example.com/model.gguf"
export LORA_URL="https://example.com/lora.safetensors"
export HF_MODEL="black-forest-labs/FLUX.1-schnell"
export SEED=1234
export OUTPUT="./my_images"
export STEPS=50
export RESOLUTION="1920x1080"
export GUIDANCE_SCALE=7.5
export STRENGTH=0.8
export WEB_PORT=9000
```

## Config File Format (prompts.json)

```json
{
  "prompts": [
    "a mystical forest with glowing mushrooms",
    "a cyberpunk city at night",
    "a dragon flying over mountains"
  ],
  "style_suffix": ", highly detailed, 8k, masterpiece",
  "negative_prompt": "blurry, low quality, distorted"
}
```

## Empty Config (for web-only mode)

```json
{
  "prompts": [],
  "style_suffix": "",
  "negative_prompt": ""
}
```

## Common Workflows

### Workflow 1: Generate and Browse
```bash
# Step 1: Generate images
go run main.go -config prompts.json

# Step 2: View in browser
go run main.go -config empty.json --web --output ./output
```

### Workflow 2: All-in-One
```bash
# Generate and immediately browse
go run main.go -config prompts.json --web
```

### Workflow 3: Image-to-Image Pipeline
```bash
# 1. Prepare input images
mkdir -p ./images
cp your_photo.jpg ./images/

# 2. Generate transformations
go run main.go -config prompts.json --img2img --strength 0.6 --web

# 3. Browse results at http://localhost:8080
```

### Workflow 4: Production Build
```bash
# Build binary
go build -o flux-gen

# Run in production
./flux-gen -config config.json --web --web-port 8080
```

## Web UI Features

### Access
```
http://localhost:8080
```

### Features
- 🖼️ **Gallery View** - Grid of all images
- 🔍 **Search** - Filter by filename
- 👁️ **Full Screen** - Click image to enlarge
- ⬇️ **Download** - Individual or batch
- 🗑️ **Delete** - Remove unwanted images
- 🔄 **Auto-Refresh** - Updates every 5 seconds

### API Endpoints
```bash
# Get images list (JSON)
curl http://localhost:8080/api/images

# Download image
curl -O http://localhost:8080/api/download/image.png

# Delete image
curl -X DELETE http://localhost:8080/api/delete/image.png
```

## Image-to-Image Strength Guide

| Strength | Effect | Use Case |
|----------|--------|----------|
| 0.1-0.3 | Minimal change | Minor tweaks, style transfer |
| 0.4-0.6 | Moderate change | Artistic interpretation |
| 0.7-0.8 | Significant change | Major transformation |
| 0.9-1.0 | Maximum creativity | Almost new image |

## Troubleshooting

### Port in Use
```bash
# Error: bind: address already in use
# Solution: Use different port
go run main.go -config config.json --web --web-port 9000
```

### No Images in Gallery
```bash
# Check output directory exists and has images
ls -la ./output

# Make sure paths match
go run main.go -config config.json --web --output ./output
```

### Python Script Not Found
```bash
# Make sure scripts exist
ls scripts/
# Should show: python_generate.py, python_img2img.py
```

### CUDA Out of Memory
```bash
# Use low VRAM mode
go run main.go -config config.json --low-vram
```

## Performance Tips

### Speed Up Generation
```bash
# Use GGUF quantized models
--gguf-model-url "https://example.com/flux-q4.gguf"

# Reduce steps (faster, lower quality)
--num_steps 20

# Lower resolution
--resolution 768x768
```

### Save Memory
```bash
# Enable low VRAM mode
--low-vram

# Close web browser when generating
# Or run generation and web server separately
```

## Directory Structure

```
your-project/
├── main.go
├── scripts/
│   ├── python_generate.py       # Text-to-image
│   └── python_img2img.py        # Image-to-image
├── pkg/
│   ├── config/
│   ├── generate/
│   ├── logging/
│   ├── utils/
│   ├── download/
│   └── webserver/               # Web UI
├── models/                      # Downloaded models
│   └── lora/                    # LoRA files
├── images/                      # Input images (img2img)
├── output/                      # Generated images
└── prompts.json                 # Your prompts config
```

## Quick Start Examples

### Example 1: First Time Setup
```bash
# 1. Create config
cat > prompts.json << 'EOF'
{
  "prompts": ["a beautiful sunset"],
  "style_suffix": "",
  "negative_prompt": ""
}
EOF

# 2. Generate with web UI
go run main.go -config prompts.json --web

# 3. Open http://localhost:8080
```

### Example 2: Batch Generation
```bash
# Generate 100 variations
cat > batch.json << 'EOF'
{
  "prompts": [
    "astronaut on mars",
    "underwater city",
    "floating islands"
  ],
  "style_suffix": ", cinematic, 8k",
  "negative_prompt": "blurry"
}
EOF

go run main.go -config batch.json --seed 42
```

### Example 3: Transform Photos
```bash
# 1. Put photos in ./images/
cp vacation/*.jpg ./images/

# 2. Transform to art
cat > art.json << 'EOF'
{
  "prompts": [
    "oil painting style",
    "watercolor painting",
    "digital art"
  ],
  "style_suffix": ", masterpiece",
  "negative_prompt": "realistic"
}
EOF

go run main.go -config art.json --img2img --strength 0.7 --web
```

## Keyboard Shortcuts (Web UI)

| Key | Action |
|-----|--------|
| `Esc` | Close full-screen image |
| `Ctrl+C` | Stop server (terminal) |

## Support

- 📖 Full docs: `docs/WEB_UI_GUIDE.md`
- 🚀 Quick start: `docs/QUICKSTART_WEB.md`
- 📝 Summary: `IMPLEMENTATION_SUMMARY.md`

---

**Quick Tip**: Bookmark this page for easy reference! 📌
