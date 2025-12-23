# Dual Mode Operation

## 🎯 Two Distinct Modes

### Mode 1: Web Server Only
**No config file needed!**
```bash
./gfluxgo --web --output ./my_gallery
```
- Starts web server immediately
- Serves existing images
- Upload interface for `./images/` directory
- Perfect for viewing results or preparing images

### Mode 2: Image Generation
**Config file required**
```bash
./gfluxgo --config prompts.json --output ./results
```
- Generates new images using AI
- Supports text-to-image and image-to-image
- Works with FLUX, Stable Diffusion, and Qwen models

## 🚀 Quick Start Examples

### Just want to view/share images?
```bash
# Start web server on port 8080
./gfluxgo --web --output ./shared_gallery

# Or custom port
./gfluxgo --web --web-port 9000 --output ./team_view
```

### Want to generate images?
```bash
# Create config file first
echo '{"prompts": ["a beautiful landscape"]}' > my_prompts.json

# Generate images
./gfluxgo --config my_prompts.json --output ./generated

# Or generate and serve immediately
./gfluxgo --config my_prompts.json --web --output ./generated
```

## 🔄 Typical Workflow

### Separate Upload and Processing:
```bash
# Terminal 1: Upload phase
./gfluxgo --web --output ./upload_phase

# Browser: Upload images to http://localhost:8080

# Terminal 2: Processing phase (when ready)
./gfluxgo --config processing.json --img2img --output ./upload_phase
```

### All-in-One:
```bash
# Generate and serve in one command
./gfluxgo --config all_in_one.json --web --output ./results
```

## 📁 Directory Management

### Web Server Creates Directories:
```bash
# This works even if ./new_gallery doesn't exist
./gfluxgo --web --output ./new_gallery
```

### Upload Directory:
- Always uses `./images/` (relative to where you run the command)
- Created automatically when you upload first image
- Images here are used for img2img processing

## ⚡ Performance Tips

### Lightweight Web Mode:
```bash
# Minimal memory usage, just serving files
./gfluxgo --web --output ./lightweight
```

### Resource-Intensive Generation:
```bash
# Uses GPU, more memory
./gfluxgo --config heavy.json --use-qwen --img2img --output ./processed
```

## 🎨 Use Case Examples

### 1. Artist Portfolio
```bash
# Just serve existing artwork
./gfluxgo --web --output ./portfolio --web-port 80
```

### 2. Batch Processing Service
```bash
# Upload endpoint for clients
./gfluxgo --web --output ./client_uploads --web-port 8080

# Processing service (runs separately)
./gfluxgo --config client_jobs.json --img2img --output ./client_results
```

### 3. Personal AI Studio
```bash
# Quick experiments
./gfluxgo --config experiments.json --output ./experiments

# Share with friends
./gfluxgo --web --output ./experiments --web-port 8081
```

## ❓ Common Questions

### Q: Do I need a config file for web-only mode?
**A: NO!** Web-only mode doesn't need any config file.

### Q: Can I run web server and generation at the same time?
**A: YES!** Run web server in one terminal, generation in another.

### Q: Where do uploaded images go?
**A:** To `./images/` directory in your current working directory.

### Q: Can I change the upload directory?
**A:** Currently fixed to `./images/` for img2img compatibility.

### Q: What if output directory doesn't exist?
**A:** Web server creates it automatically.

## 🛠️ Troubleshooting

### Web server won't start:
```bash
# Check if port is in use
./gfluxgo --web --web-port 8081 --output ./test

# Check permissions
./gfluxgo --web --output /tmp/test_gallery
```

### No images in gallery:
- Directory might be empty
- Generate images first or upload via web interface
- Check browser console for errors

### Upload fails:
- Check file size (max 10MB)
- Supported formats: PNG, JPG, WebP, GIF, BMP
- Check browser console for error messages

## 📊 Monitoring

### Web Server Logs:
```
🌐 Web UI started at http://localhost:8080
📁 Serving images from: ./my_gallery
📤 Upload images to ./images/ directory for img2img processing
```

### Generation Logs:
```
Starting FLUX generation for 4 images...
✓ Model loaded in 15.2s
✓ Saved to output1.png in 8.1s
```

## 🔒 Security Notes

### Web Server:
- Localhost only by default
- No authentication (for local use)
- File upload limits (10MB)
- File type restrictions

### For Production:
- Use reverse proxy (nginx, Apache)
- Add authentication layer
- Monitor disk usage
- Regular backups

## 🚀 Next Steps

1. **Try web-only mode first**: `./gfluxgo --web --output ./test`
2. **Upload some images** via the web interface
3. **Create a config file** for generation
4. **Experiment** with different models and modes