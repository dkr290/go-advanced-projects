# Web Server with Image Upload Feature

## New Feature: Image Upload via Web Interface

The web server now includes a powerful image upload feature that makes it easy to prepare images for image-to-image processing.

## 🚀 Quick Start

### 1. Start the Web Server
```bash
./gfluxgo --config config.json --web --web-port 8080
```

### 2. Open Browser
Navigate to: `http://localhost:8080`

### 3. Upload Images
- Use the upload section at the top of the page
- Drag & drop images or click to browse
- Choose destination:
  - **Upload to ./images/ directory**: For img2img input images
  - **Upload to gallery**: To add images to the main gallery

## 📁 Upload to ./images/ Directory

This is specifically designed for Qwen/FLUX image-to-image processing:

### Why Use This?
- Upload images that will be used as input for img2img generation
- Images are saved to `./images/` directory
- Automatically creates directory if it doesn't exist
- Preserves original filenames (adds timestamp for duplicates)

### How It Works:
1. Upload images via web interface to `./images/`
2. Configure your prompts in `config.json`
3. Run img2img generation:
   ```bash
   ./gfluxgo --config config.json --hf-model "Qwen/Qwen-Image-Edit" --use-qwen --img2img
   ```
4. Generated images appear in the web gallery

## 🖼️ Upload to Gallery

Use this to:
- Add existing images to the gallery
- Share images with others via the web interface
- Create a collection of reference images

## 🔧 API Endpoints

### Upload to Images Directory
```bash
curl -X POST -F "image=@your_image.jpg" http://localhost:8080/api/upload-images
```

### Upload to Gallery
```bash
curl -X POST -F "image=@your_image.png" http://localhost:8080/api/upload
```

## 📋 Supported Formats
- **Images**: PNG, JPG, JPEG, WebP, GIF, BMP
- **Max Size**: 10MB per file
- **Multiple Files**: Select and upload multiple files at once

## 🎯 Use Cases

### 1. Batch Image Preparation
```bash
# Step 1: Start web server
./gfluxgo --web --web-port 8080 --output ./batch_output

# Step 2: Upload multiple images via web interface
# Step 3: Create config with matching prompts
# Step 4: Run batch img2img processing
./gfluxgo --config batch_config.json --img2img
```

### 2. Collaborative Workflow
1. Team members upload images via web interface
2. Configure different style prompts
3. Generate variations for each image
4. Share results via the same web interface

### 3. Quick Testing
```bash
# Test with minimal setup
mkdir -p test_output
./gfluxgo --web --web-port 8080 --output ./test_output

# Upload test images via browser
# No need to manually copy files to images directory
```

## 🛡️ Security Features

- **File Type Validation**: Only image files allowed
- **Size Limits**: Prevents server overload
- **Path Security**: Blocks directory traversal attempts
- **Duplicate Handling**: Prevents file overwrites

## ⚡ Performance Tips

1. **Compress Images**: Upload compressed images for faster transfer
2. **Batch Upload**: Upload multiple files at once
3. **Use WebP**: Smaller file sizes with good quality
4. **Monitor Progress**: Progress bar shows upload status

## 🔍 Troubleshooting

### Common Issues:

1. **"Invalid file type"**
   - Solution: Use supported formats (PNG, JPG, WebP, etc.)

2. **"File too large"**
   - Solution: Compress image or reduce dimensions

3. **Upload fails silently**
   - Solution: Check browser console for errors
   - Solution: Verify server has write permissions

4. **Images not appearing in gallery**
   - Solution: Click "Refresh" button
   - Solution: Check server output directory

### Debug Mode:
```bash
# Enable debug logging
./gfluxgo --web --debug --output ./debug_output
```

## 📊 Monitoring

The web interface provides:
- Upload progress bar
- Success/error messages
- File count and size information
- Auto-refresh of gallery after upload

## 🎨 Customization

### Modify Upload Limits:
Edit `pkg/webserver/webserver.go`:
```go
// Change max file size (currently 10MB)
err := r.ParseMultipartForm(10 << 20) // 10 MB

// Add/remove allowed file types
allowedExts := []string{".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp"}
```

### Custom Styling:
Modify the `galleryHTML` constant in the same file to change:
- Upload area appearance
- Progress bar colors
- Status message styling

## 🔗 Integration Examples

### With Qwen-Image-Edit:
```bash
# Complete workflow example
./gfluxgo --web --web-port 8080 --output ./qwen_results

# 1. Upload images via web interface
# 2. Create config_qwen.json with prompts
# 3. Run img2img
./gfluxgo --config config_qwen.json --hf-model "Qwen/Qwen-Image-Edit" --use-qwen --img2img
```

### With FLUX:
```bash
# Similar workflow for FLUX
./gfluxgo --config config_flux.json --img2img
```

## 📈 Next Steps

Future enhancements planned:
- [ ] Image cropping before upload
- [ ] Bulk upload via ZIP files
- [ ] Image metadata editing
- [ ] User accounts and permissions
- [ ] Upload history and statistics

## ❓ Need Help?

1. Check the web interface for error messages
2. Run with `--debug` flag for detailed logs
3. Review `WEB_UPLOAD_FEATURE.md` for technical details
4. Test with `test_upload.html` for basic functionality