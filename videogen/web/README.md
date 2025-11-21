# Wan2.1 Video Generator - Web Interface

Beautiful web UI for the Wan2.1 Video Generation API built with Go, HTMX, and Bootstrap 5.

## 🎨 Features

- ✅ **Modern UI** - Clean, responsive design with Bootstrap 5
- ✅ **Real-time Updates** - HTMX for seamless interactions
- ✅ **Text-to-Video** - Generate videos from text descriptions
- ✅ **Image-to-Video** - Animate static images
- ✅ **Video-to-Video** - Transform existing videos
- ✅ **Video Gallery** - Browse all generated videos
- ✅ **Settings Management** - Configure models and parameters
- ✅ **Dark Mode** - Eye-friendly dark theme
- ✅ **Progress Tracking** - Real-time job status updates

## 📋 Prerequisites

- Go 1.21+
- Running Wan2.1 API server (from the main project)

## 🚀 Quick Start

### 1. Setup

```bash
cd videogen/web

# Install dependencies
go mod download

# Copy environment config
cp .env.example .env

# Edit config (optional)
nano .env
```

### 2. Start the API Server

First, make sure the main API server is running:

```bash
# In the main project directory
cd ../..  # Back to wan2-video-server

# Terminal 1: Python backend
cd python_backend
source venv/bin/activate
python server.py

# Terminal 2: Go API server
./wan2-video-server
```

### 3. Start the Web Interface

```bash
# In videogen/web directory
go run main.go
```

The web interface will be available at: **http://localhost:3000**

## 📁 Project Structure

```
videogen/web/
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── .env.example           # Configuration template
│
├── handlers/              # HTTP handlers
│   ├── pages.go          # Page handlers
│   └── api.go            # API proxy handlers
│
├── middleware/            # HTTP middleware
│   └── cors.go           # CORS middleware
│
├── templates/             # HTML templates
│   ├── layouts/
│   │   └── base.html     # Base layout
│   ├── pages/
│   │   ├── index.html    # Home page
│   │   ├── text-to-video.html
│   │   ├── image-to-video.html
│   │   ├── video-to-video.html
│   │   ├── gallery.html
│   │   └── settings.html
│   └── components/        # Reusable components
│       ├── job-status.html
│       ├── video-result.html
│       └── error.html
│
└── static/                # Static assets
    ├── css/
    │   └── style.css     # Custom styles
    └── js/
        └── app.js        # JavaScript
```

## ⚙️ Configuration

Edit `.env` file:

```env
# Web UI Port
PORT=3000

# API Backend URL (the main Wan2.1 server)
API_BASE_URL=http://localhost:8080

# Upload settings
MAX_UPLOAD_SIZE=100MB
ALLOWED_IMAGE_TYPES=jpg,jpeg,png,gif,webp
ALLOWED_VIDEO_TYPES=mp4,avi,mov,webm

# UI Settings
VIDEOS_PER_PAGE=12
ENABLE_GALLERY=true
```

## 🎯 Usage

### Text-to-Video

1. Navigate to **Text to Video** page
2. Enter your prompt (e.g., "A cat playing with a ball")
3. Optionally adjust advanced settings
4. Click **Generate Video**
5. Wait for processing (real-time status updates)
6. Download your video!

### Image-to-Video

1. Navigate to **Image to Video** page
2. Upload an image
3. Enter optional prompt for guidance
4. Adjust settings
5. Generate and download

### Video-to-Video

1. Navigate to **Video to Video** page
2. Upload a video
3. Enter transformation prompt
4. Adjust strength and other settings
5. Generate transformed video

## 🎨 UI Features

### Real-time Updates with HTMX

The UI uses HTMX for seamless updates without page reloads:

- ✅ Form submissions
- ✅ Job status polling
- ✅ Dynamic content loading
- ✅ Error handling

### Responsive Design

Works perfectly on:
- 💻 Desktop
- 📱 Mobile
- 📱 Tablet

### Dark Mode

Beautiful dark theme optimized for long sessions.

## 🔧 Development

### Run in Development Mode

```bash
# With auto-reload (using air)
go install github.com/cosmtrek/air@latest
air

# Or standard run
go run main.go
```

### Build for Production

```bash
# Build binary
go build -o videogen-web main.go

# Run
./videogen-web
```

## 🐳 Docker

```bash
# Build
docker build -t videogen-web .

# Run
docker run -p 3000:3000 \
  -e API_BASE_URL=http://api-server:8080 \
  videogen-web
```

## 📊 API Endpoints

The web UI proxies these endpoints to the main API:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/generate/text-to-video` | POST | Generate from text |
| `/api/generate/image-to-video` | POST | Generate from image |
| `/api/generate/video-to-video` | POST | Transform video |
| `/api/job/:id` | GET | Check job status |
| `/api/models` | GET | List models |
| `/api/switch-model` | POST | Switch model |

## 🎭 Screenshots

### Home Page
- Beautiful hero section
- Feature cards
- Quick start guide

### Generation Pages
- Clean form interface
- Advanced settings accordion
- Real-time progress
- Video preview and download

### Gallery
- Grid layout
- Video thumbnails
- Quick preview

## 🛠️ Technologies Used

- **Backend**: Go + Gin
- **Frontend**: HTMX + Bootstrap 5
- **Icons**: Bootstrap Icons
- **Styling**: Custom CSS + Bootstrap

## 🔗 Integration

This web UI is designed as a separate microservice that communicates with the main Wan2.1 API server.

**Architecture:**

```
Browser → Web UI (Port 3000) → API Server (Port 8080) → Python Backend (Port 5000) → GPU
```

## 📝 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 3000 | Web UI port |
| `GIN_MODE` | release | Gin mode (debug/release) |
| `API_BASE_URL` | http://localhost:8080 | API server URL |
| `MAX_UPLOAD_SIZE` | 100MB | Max file upload size |

## 🐛 Troubleshooting

### "Failed to connect to API server"

- Ensure the main API server is running on port 8080
- Check `API_BASE_URL` in `.env`
- Verify firewall settings

### "HTMX not working"

- Check browser console for errors
- Ensure CDN is accessible
- Check CORS settings

### Videos not loading

- Verify API server is serving files correctly
- Check network tab in browser dev tools
- Ensure correct video path

## 📚 Documentation

- [Main API Documentation](../../README.md)
- [HTMX Documentation](https://htmx.org/docs/)
- [Bootstrap Documentation](https://getbootstrap.com/docs/5.3/)

## 🎉 Features Coming Soon

- [ ] User authentication
- [ ] Video editing
- [ ] Batch processing
- [ ] Custom presets
- [ ] Video sharing
- [ ] Advanced filters

## 💡 Tips

1. **Faster Generation**: Use lower resolution and fewer frames for testing
2. **Better Quality**: Increase inference steps and guidance scale
3. **Reproducibility**: Use same seed value
4. **Keyboard Shortcut**: Ctrl/Cmd + Enter to submit forms

## 📄 License

MIT License - See main project LICENSE file

---

**Built with ❤️ using Go, HTMX, and Bootstrap**
