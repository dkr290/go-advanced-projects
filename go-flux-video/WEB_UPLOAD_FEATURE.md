# Web Server Upload Feature

## Overview

The web server now includes image upload functionality that allows users to:
1. Upload images to the `./images/` directory for image-to-image processing
2. Upload images directly to the gallery/output directory
3. Use drag-and-drop or file browser for uploading
4. Track upload progress with visual feedback

## New API Endpoints

### 1. `/api/upload` - Upload to Gallery
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Field**: `image` (file)
- **Response**: JSON with upload status and file info
- **Destination**: Output directory (same as generated images)

### 2. `/api/upload-images` - Upload to ./images/ directory
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Field**: `image` (file)
- **Response**: JSON with upload status and file info
- **Destination**: `./images/` directory (for img2img processing)

## Web Interface Features

### Upload Section
The web interface now includes an upload section with:
- Drag-and-drop area for easy file selection
- File browser integration
- Progress bar for upload tracking
- Success/error status messages
- Two upload buttons:
  - **Upload to ./images/ directory**: For img2img input images
  - **Upload to gallery**: For adding images to the main gallery

### Supported File Formats
- PNG, JPG, JPEG, WebP, GIF, BMP
- Maximum file size: 10MB per file
- Multiple file selection supported

## Usage Examples

### 1. Using the Web Interface
1. Start the web server:
   ```bash
   ./gfluxgo --config config.json --web --web-port 8080
   ```

2. Open browser: `http://localhost:8080`

3. Use the upload section to:
   - Drag and drop images
   - Or click to browse files
   - Select destination (images directory or gallery)
   - Monitor upload progress

### 2. Using cURL (API)
```bash
# Upload to images directory
curl -X POST -F "image=@/path/to/your/image.jpg" http://localhost:8080/api/upload-images

# Upload to gallery
curl -X POST -F "image=@/path/to/your/image.png" http://localhost:8080/api/upload
```

### 3. Using JavaScript/Fetch
```javascript
// Upload to images directory
async function uploadImage(file) {
    const formData = new FormData();
    formData.append('image', file);
    
    const response = await fetch('/api/upload-images', {
        method: 'POST',
        body: formData
    });
    
    return await response.json();
}
```

## Integration with Image-to-Image Processing

### Workflow:
1. **Upload images** to `./images/` directory via web interface
2. **Configure prompts** in your JSON config file
3. **Run img2img** with Qwen/FLUX model:
   ```bash
   ./gfluxgo --config config.json --hf-model "Qwen/Qwen-Image-Edit" --use-qwen --img2img
   ```
4. **View results** in the web gallery

### Automatic Directory Creation
- The `./images/` directory is automatically created if it doesn't exist
- Files are saved with original names (timestamp added if duplicate exists)

## Security Features

1. **File Type Validation**: Only allowed image formats are accepted
2. **Size Limits**: 10MB maximum per file
3. **Path Security**: Prevents directory traversal attacks
4. **Duplicate Handling**: Automatic timestamp addition for duplicate filenames

## Error Handling

### Common Errors and Solutions:
1. **"Invalid file type"**: Use supported formats (PNG, JPG, WebP, etc.)
2. **"File too large"**: Reduce image size or compress before uploading
3. **"Failed to create directory"**: Check filesystem permissions
4. **"No image file provided"**: Ensure file field is named `image`

### Response Format:
```json
{
    "status": "success",
    "filename": "upload_20241223_143045_image.jpg",
    "path": "/images/upload_20241223_143045_image.jpg",
    "relative_path": "images/upload_20241223_143045_image.jpg",
    "size": 1024576,
    "message": "Image uploaded to ./images/ directory for img2img processing"
}
```

## Testing

### Quick Test:
1. Run the web server:
   ```bash
   ./gfluxgo --web --web-port 8080 --output ./test_output
   ```

2. Open `test_upload.html` in browser
3. Test upload functionality without generating images

### Integration Test:
1. Upload images to `./images/` via web interface
2. Create config with matching number of prompts
3. Run img2img generation
4. Verify results appear in web gallery

## Browser Compatibility

- Modern browsers with HTML5 File API support
- Drag-and-drop works in Chrome, Firefox, Safari, Edge
- Fallback to file browser for older browsers

## Performance Considerations

1. **Memory Usage**: 10MB limit per file prevents memory issues
2. **Concurrent Uploads**: Files are uploaded sequentially
3. **Progress Updates**: Real-time progress bar for better UX
4. **Auto-refresh**: Gallery refreshes automatically after upload

## Next Steps

Potential enhancements:
1. Bulk upload with zip file support
2. Image preview before upload
3. Image cropping/resizing before upload
4. Upload history/log
5. User authentication for uploads