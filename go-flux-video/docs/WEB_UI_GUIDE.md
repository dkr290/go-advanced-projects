# Web UI for FLUX Image Generation

## Overview

The web UI provides a browser-based gallery to view, download, and manage generated images.

## Features

✨ **Image Gallery**
- Grid view of all generated images
- Auto-refresh every 5 seconds
- Full-size image modal view
- Responsive design (mobile-friendly)

🔍 **Search & Filter**
- Real-time search by filename
- Filter images instantly

⬇️ **Download Options**
- Download individual images
- Download all images at once
- Direct file serving

🗑️ **Image Management**
- Delete unwanted images
- View image metadata (size, date)
- Organized by generation time

## Usage

### 1. Generate Images First

Run your image generation as normal:

```bash
go run main.go -config config.json
```

### 2. Start Web Server

Add the `--web` flag to enable the web UI:

```bash
go run main.go -config config.json --web
```

Or specify a custom port:

```bash
go run main.go -config config.json --web --web-port 9000
```

### 3. Access in Browser

Open your browser and navigate to:

```
http://localhost:8080
```

Or if you specified a custom port:

```
http://localhost:9000
```

## Combined Workflow

### Generate and View Immediately

```bash
# Generate images and start web server in one command
go run main.go -config config.json --web

# The web server will start after generation completes
# Images will appear in the gallery automatically
```

### Web-Only Mode (View Existing Images)

If you just want to browse existing images without generating new ones, you can create a minimal config:

```json
{
  "prompts": [],
  "style_suffix": "",
  "negative_prompt": ""
}
```

Then run:

```bash
go run main.go -config empty_config.json --web --output ./output
```

## Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--web` | Enable web server mode | `false` |
| `--web-port` | Web server port | `8080` |
| `--output` | Output directory to serve | `./output` |

## Environment Variables

You can also use environment variables:

```bash
export WEB_PORT=9000
go run main.go -config config.json --web
```

## Web UI Features

### Image Card
Each image shows:
- **Thumbnail preview** - Click to view full size
- **Filename** - Hover for full name
- **File size** - In megabytes
- **Generation date** - When it was created
- **Download button** - Save to your device
- **Delete button** - Remove from gallery

### Controls
- **🔍 Search** - Filter images by filename
- **🔄 Refresh** - Manually reload gallery
- **⬇️ Download All** - Batch download all images

### Modal View
Click any image to view it full-screen:
- Click anywhere to close
- Press `Esc` to close
- High-resolution display

## Integration with Image Generation

### Text-to-Image + Web Server

```bash
go run main.go \
  -config prompts.json \
  --output ./generated_images \
  --web \
  --web-port 8080
```

### Image-to-Image + Web Server

```bash
go run main.go \
  -config prompts.json \
  --img2img \
  --output ./img2img_output \
  --web \
  --web-port 8080
```

## API Endpoints

The web server exposes these endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Main gallery page |
| `/api/images` | GET | JSON list of all images |
| `/api/download/{filename}` | GET | Download specific image |
| `/api/delete/{filename}` | DELETE | Delete specific image |
| `/images/{filename}` | GET | Serve image file |

### Example API Usage

#### Get all images (JSON)

```bash
curl http://localhost:8080/api/images
```

Response:
```json
[
  {
    "filename": "mystical_forest_1.png",
    "path": "/images/mystical_forest_1.png",
    "size": 2457600,
    "mod_time": "2025-12-23T10:30:00Z",
    "thumbnail": "/images/mystical_forest_1.png",
    "download_url": "/api/download/mystical_forest_1.png"
  }
]
```

#### Download an image

```bash
curl -O http://localhost:8080/api/download/mystical_forest_1.png
```

#### Delete an image

```bash
curl -X DELETE http://localhost:8080/api/delete/mystical_forest_1.png
```

## Security Notes

⚠️ **Important**: This web server is intended for local use only.

- No authentication is implemented
- Anyone with access to the port can view/delete images
- Do not expose to the public internet without adding authentication
- Use a firewall to restrict access if needed

For production use, consider:
- Adding authentication (JWT, OAuth, etc.)
- Using HTTPS/TLS
- Rate limiting
- Access control lists

## Troubleshooting

### Port Already in Use

```
Error: listen tcp :8080: bind: address already in use
```

**Solution**: Use a different port
```bash
go run main.go -config config.json --web --web-port 9000
```

### No Images Showing

1. Check that images exist in the output directory
2. Verify the `--output` flag points to the correct directory
3. Check browser console for errors
4. Try refreshing the page

### Cannot Access from Another Device

By default, the server binds to all network interfaces. To access from another device:

1. Find your machine's IP address:
   ```bash
   ip addr show  # Linux
   ifconfig      # macOS
   ipconfig      # Windows
   ```

2. Access from another device:
   ```
   http://<YOUR_IP>:8080
   ```

3. Make sure your firewall allows connections on the port

## Customization

The web UI template is embedded in `pkg/webserver/webserver.go`. You can customize:

- Colors and styling (CSS in `<style>` section)
- Layout and grid size
- Auto-refresh interval (currently 5 seconds)
- Image card information displayed

## Performance

- **Auto-refresh**: Gallery updates every 5 seconds
- **Lazy loading**: Images load as you scroll
- **Efficient**: Only loads metadata, not full images
- **Responsive**: Works on mobile, tablet, and desktop

## Example Complete Workflow

```bash
# 1. Generate images with text-to-image
go run main.go -config my_prompts.json --output ./gallery

# 2. Start web server to view results
go run main.go -config empty.json --web --output ./gallery

# 3. Open browser to http://localhost:8080

# 4. Browse, download, or delete images from the web UI
```

## Tips

💡 **Keep the web server running** while generating images in another terminal
💡 **Use descriptive filenames** in your prompts for better searchability
💡 **Organize by output directory** for different projects
💡 **Bookmark the URL** for quick access

---

Enjoy your FLUX image gallery! 🎨✨
