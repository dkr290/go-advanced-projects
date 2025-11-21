# API Reference

Complete API documentation for Wan2.1 Video Server.

## Base URL

```
http://localhost:8080
```

## Table of Contents

- [Health & Info](#health--info)
- [Video Generation](#video-generation)
- [Job Management](#job-management)
- [Model Management](#model-management)

---

## Health & Info

### Health Check

Check if the server is running and healthy.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "service": "wan2-video-server",
  "version": "1.0.0"
}
```

### Model Information

Get information about the loaded model.

**Endpoint:** `GET /api/v1/model/info`

**Response:**
```json
{
  "name": "Wan2.1",
  "version": "1.0.0",
  "provider": "huggingface",
  "gpu_enabled": true,
  "gpu_device_id": 0,
  "cache_dir": "./models"
}
```

---

## Video Generation

### Text-to-Video

Generate a video from a text prompt.

**Endpoint:** `POST /api/v1/generate/text-to-video`

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "prompt": "A cat playing with a ball of yarn",
  "negative_prompt": "blurry, low quality, distorted",
  "num_frames": 64,
  "fps": 24,
  "width": 512,
  "height": 512,
  "seed": 42,
  "guidance_scale": 7.5,
  "num_inference_steps": 50
}
```

**Parameters:**

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `prompt` | string | Yes | - | Text description of the video to generate |
| `negative_prompt` | string | No | "" | Things to avoid in generation |
| `num_frames` | int | No | 64 | Number of frames (max: 128) |
| `fps` | int | No | 24 | Frames per second |
| `width` | int | No | 512 | Video width in pixels |
| `height` | int | No | 512 | Video height in pixels |
| `seed` | int | No | -1 | Random seed (-1 for random) |
| `guidance_scale` | float | No | 7.5 | How closely to follow prompt (1-20) |
| `num_inference_steps` | int | No | 50 | Quality vs speed (25-100) |

**Response:**
```json
{
  "job_id": "job_1234567890_abcdef123456",
  "status": "processing",
  "message": "Video generation started"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/generate/text-to-video \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A beautiful sunset over the ocean",
    "num_frames": 64,
    "fps": 24,
    "guidance_scale": 7.5
  }'
```

---

### Image-to-Video

Generate a video from an input image and optional text prompt.

**Endpoint:** `POST /api/v1/generate/image-to-video`

**Content-Type:** `multipart/form-data`

**Form Fields:**

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `image` | file | Yes | - | Input image file (JPG, PNG) |
| `prompt` | string | No | "" | Additional guidance text |
| `negative_prompt` | string | No | "" | Things to avoid |
| `num_frames` | int | No | 64 | Number of frames |
| `fps` | int | No | 24 | Frames per second |
| `width` | int | No | 512 | Video width |
| `height` | int | No | 512 | Video height |
| `seed` | int | No | -1 | Random seed |
| `guidance_scale` | float | No | 7.5 | Prompt adherence |
| `num_inference_steps` | int | No | 50 | Quality steps |

**Response:**
```json
{
  "job_id": "job_1234567890_abcdef123456",
  "status": "processing",
  "message": "Video generation started"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/generate/image-to-video \
  -F "image=@/path/to/image.jpg" \
  -F "prompt=A dog running in a park" \
  -F "num_frames=64" \
  -F "fps=24"
```

---

### Video-to-Video

Transform an input video based on a text prompt.

**Endpoint:** `POST /api/v1/generate/video-to-video`

**Content-Type:** `multipart/form-data`

**Form Fields:**

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `video` | file | Yes | - | Input video file (MP4, AVI) |
| `prompt` | string | No | "" | Transformation guidance |
| `negative_prompt` | string | No | "" | Things to avoid |
| `fps` | int | No | 24 | Output FPS |
| `strength` | float | No | 0.8 | Transformation strength (0-1) |
| `seed` | int | No | -1 | Random seed |
| `guidance_scale` | float | No | 7.5 | Prompt adherence |
| `num_inference_steps` | int | No | 50 | Quality steps |

**Response:**
```json
{
  "job_id": "job_1234567890_abcdef123456",
  "status": "processing",
  "message": "Video generation started"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/generate/video-to-video \
  -F "video=@/path/to/video.mp4" \
  -F "prompt=Transform into anime style" \
  -F "strength=0.8"
```

---

## Job Management

### Get Job Status

Check the status of a generation job.

**Endpoint:** `GET /api/v1/job/:id`

**Parameters:**
- `id` (path) - Job ID returned from generation request

**Response (Processing):**
```json
{
  "job_id": "job_1234567890_abcdef123456",
  "status": "processing",
  "message": "Generating video...",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:15Z"
}
```

**Response (Completed):**
```json
{
  "job_id": "job_1234567890_abcdef123456",
  "status": "completed",
  "message": "Generation complete",
  "output_path": "./outputs/text2video_20240115_103045_abc123.mp4",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:45Z"
}
```

**Response (Failed):**
```json
{
  "job_id": "job_1234567890_abcdef123456",
  "status": "failed",
  "message": "Error: Out of memory",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:20Z"
}
```

**Status Values:**
- `pending` - Job queued
- `processing` - Currently generating
- `completed` - Successfully completed
- `failed` - Generation failed

**Example:**
```bash
curl http://localhost:8080/api/v1/job/job_1234567890_abcdef123456
```

---

## Model Management

### List Models

Get a list of available models.

**Endpoint:** `GET /api/v1/models`

**Response:**
```json
{
  "models": [
    {
      "id": "Lightricks/LTX-Video",
      "name": "Wan2.1",
      "provider": "huggingface",
      "cache_dir": "./models",
      "downloaded": true
    }
  ]
}
```

**Example:**
```bash
curl http://localhost:8080/api/v1/models
```

---

### Download Model

Download a model from Hugging Face.

**Endpoint:** `POST /api/v1/models/download`

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "model_id": "Lightricks/LTX-Video"
}
```

**Response:**
```json
{
  "message": "Model download started",
  "model_id": "Lightricks/LTX-Video"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/models/download \
  -H "Content-Type: application/json" \
  -d '{"model_id": "Lightricks/LTX-Video"}'
```

---

## Error Responses

All endpoints may return error responses:

**400 Bad Request:**
```json
{
  "error": "Invalid request: prompt is required"
}
```

**429 Too Many Requests:**
```json
{
  "error": "Too many concurrent requests. Please try again later."
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error: model not loaded"
}
```

**503 Service Unavailable:**
```json
{
  "error": "Model not loaded"
}
```

---

## Rate Limiting

- Maximum concurrent requests: 2 (configurable)
- Requests exceeding limit will receive 429 status
- Wait for current jobs to complete before submitting new ones

---

## File Access

Generated videos are accessible via the static file server:

```
http://localhost:8080/outputs/<filename>
```

Example:
```
http://localhost:8080/outputs/text2video_20240115_103045_abc123.mp4
```

---

## Best Practices

1. **Use appropriate parameters:**
   - Lower `num_frames` for faster generation
   - Higher `num_inference_steps` for better quality
   - Use consistent `seed` for reproducible results

2. **Monitor job status:**
   - Poll `/api/v1/job/:id` periodically
   - Don't submit duplicate requests

3. **Handle errors gracefully:**
   - Check for rate limiting (429)
   - Retry failed jobs with adjusted parameters
   - Reduce resolution if out of memory

4. **Optimize prompts:**
   - Be specific and descriptive
   - Use negative prompts to avoid unwanted elements
   - Adjust `guidance_scale` to balance creativity and adherence
