# Usage Modes

The application now supports two distinct modes of operation:

## 1. Web Server Only Mode
**Purpose**: Serve existing images and upload new images for later processing
**Config File**: NOT required
**When to use**: When you only want to view images or upload images for future processing

### Commands:
```bash
# Basic web server
./gfluxgo --web --output ./my_gallery

# Custom port
./gfluxgo --web --web-port 9000 --output ./my_gallery

# With debug logging
./gfluxgo --web --output ./my_gallery --debug
```

### What it does:
- Starts web server on specified port (default: 8080)
- Serves images from output directory
- Provides upload interface for `./images/` directory
- **Does NOT** generate any images
- **Does NOT** require config file

### Use Case Examples:
1. **View existing results**: `./gfluxgo --web --output ./previous_results`
2. **Upload images for later**: `./gfluxgo --web --output ./to_process_later`
3. **Share gallery with team**: `./gfluxgo --web --web-port 8080 --output ./team_shared`

## 2. Image Generation Mode
**Purpose**: Generate new images using AI models
**Config File**: REQUIRED
**When to use**: When you want to generate new images

### Commands:
```bash
# Text-to-image generation
./gfluxgo --config config.json --output ./results

# Image-to-image generation
./gfluxgo --config config.json --img2img --output ./results

# With Qwen model
./gfluxgo --config config.json --use-qwen --img2img --output ./results

# With web server (generate + serve)
./gfluxgo --config config.json --web --output ./results
```

### What it does:
- Reads prompts from config file
- Generates images using specified model
- Saves results to output directory
- Optionally starts web server after generation

## 3. Combined Mode (Generate + Serve)
**Purpose**: Generate images and immediately serve them
**Config File**: REQUIRED
**When to use**: When you want to generate and view results immediately

### Command:
```bash
./gfluxgo --config config.json --web --output ./results
```

### What it does:
1. Generates images first
2. Then starts web server with results
3. One command does both

## Directory Structure

### For Web-Only Mode:
```
your-project/
├── images/           # Upload via web interface
│   └── uploaded.jpg
├── my_gallery/       # Output directory (--output)
│   └── existing.png
└── gfluxgo           # Executable
```

### For Image Generation:
```
your-project/
├── config.json       # REQUIRED: Prompt configuration
├── images/           # Input images for img2img
│   └── input.jpg
├── results/          # Output directory (--output)
│   └── generated.png
└── gfluxgo           # Executable
```

## Workflow Examples

### Example 1: Separate Upload and Generation
```bash
# Terminal 1: Start web server for uploading
./gfluxgo --web --output ./batch_processing

# Browser: Upload images to ./images/ via http://localhost:8080

# Terminal 2: Generate images (different terminal)
./gfluxgo --config prompts.json --img2img --output ./batch_processing
```

### Example 2: Quick Test
```bash
# Just view what you have
./gfluxgo --web --output ./test_view

# Generate some test images
echo '{"prompts": ["test image"]}' > test.json
./gfluxgo --config test.json --output ./test_view
```

### Example 3: Production Pipeline
```bash
# Phase 1: Upload phase (team uploads images)
./gfluxgo --web --web-port 8080 --output ./phase1_uploads

# Phase 2: Processing phase (AI generates variations)
./gfluxgo --config production_prompts.json --use-qwen --img2img --output ./phase2_results

# Phase 3: Review phase (team reviews results)
./gfluxgo --web --web-port 8081 --output ./phase2_results
```

## Environment Variables

### Web-Only Mode:
```bash
export OUTPUT="./web_gallery"
export WEB_PORT="9000"
./gfluxgo --web
```

### Image Generation Mode:
```bash
export CONFIG="config.json"
export OUTPUT="./results"
export HF_MODEL="Qwen/Qwen-Image-Edit"
export USE_QWEN="true"
./gfluxgo
```

## Error Messages

### Common Errors and Solutions:

1. **"Error: Configuration file is required for image generation"**
   - You're trying to generate images without a config file
   - Solution: Add `--config config.json` or use `--web` for web-only mode

2. **"Error reading config file"**
   - Config file doesn't exist or can't be read
   - Solution: Check file path and permissions

3. **Web server starts but no images shown**
   - Output directory might be empty
   - Solution: Generate images first or upload via web interface

## Quick Reference

| Mode | Command | Config Required | Generates Images | Serves Web UI |
|------|---------|----------------|------------------|---------------|
| Web Only | `./gfluxgo --web` | ❌ No | ❌ No | ✅ Yes |
| Generate Only | `./gfluxgo --config file.json` | ✅ Yes | ✅ Yes | ❌ No |
| Generate + Serve | `./gfluxgo --config file.json --web` | ✅ Yes | ✅ Yes | ✅ Yes |

## Tips

1. **Use different ports** for multiple web servers:
   ```bash
   ./gfluxgo --web --web-port 8080 --output ./server1
   ./gfluxgo --web --web-port 8081 --output ./server2
   ```

2. **Default output directory** is `./output` if not specified

3. **Web server auto-creates** output directory if it doesn't exist

4. **Upload directory** is always `./images/` relative to where you run the command

5. **For production**, run web server in background:
   ```bash
   nohup ./gfluxgo --web --output /var/www/gallery > server.log 2>&1 &
   ```