# Quick Start Guide: Web UI

## Step 1: Generate Some Images

First, generate some images using your existing workflow:

```bash
go run main.go -config config.json
```

This will create images in your output directory (default: `./output`)

## Step 2: Start the Web Server

Now start the web server to view your images:

```bash
go run main.go -config config.json --web
```

You'll see:
```
🌐 Web UI started at http://localhost:8080
📁 Serving images from: ./output
Press Ctrl+C to stop the server
```

## Step 3: Open Your Browser

Navigate to:
```
http://localhost:8080
```

You'll see a beautiful gallery with all your generated images!

## One-Command Workflow

Generate images and start web server in one go:

```bash
go run main.go -config config.json --web
```

The web server will start automatically after image generation completes.

## Custom Port

If port 8080 is already in use:

```bash
go run main.go -config config.json --web --web-port 9000
```

Then access at: `http://localhost:9000`

## Web-Only Mode (No Generation)

To just browse existing images without generating new ones:

```bash
# Create minimal config
echo '{"prompts":[],"style_suffix":"","negative_prompt":""}' > empty.json

# Start web server
go run main.go -config empty.json --web --output ./your_images_folder
```

## Features You Can Try

1. **Click an image** - Opens full-screen view
2. **Search box** - Type to filter images by name
3. **Download button** - Save individual images
4. **Download All** - Batch download all images
5. **Delete button** - Remove unwanted images
6. **Auto-refresh** - Gallery updates every 5 seconds

## Example Commands

### Text-to-Image with Web UI
```bash
go run main.go \
  -config prompts.json \
  --output ./my_gallery \
  --web \
  --web-port 8080
```

### Image-to-Image with Web UI
```bash
go run main.go \
  -config prompts.json \
  --img2img \
  --output ./img2img_results \
  --web \
  --web-port 8080
```

### Web Server Only
```bash
go run main.go \
  -config empty.json \
  --web \
  --output ./existing_images
```

## Troubleshooting

**Problem**: No images appear in gallery
- Make sure images exist in the output directory
- Check the output path matches: `--output ./your_folder`
- Try clicking the "🔄 Refresh" button

**Problem**: Port already in use
- Use a different port: `--web-port 9000`

**Problem**: Can't access from another device
- Find your IP: `ip addr show` or `ifconfig`
- Access at: `http://YOUR_IP:8080`
- Make sure firewall allows the port

## Next Steps

- Read the full guide: `docs/WEB_UI_GUIDE.md`
- Customize the web UI styling in `pkg/webserver/webserver.go`
- Build the binary for production: `go build -o flux-gen`

Enjoy! 🎨✨
