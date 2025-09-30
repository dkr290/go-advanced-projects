# Docker Image Builder API

A Golang Fiber-based API for building Docker images dynamically based on different model versions and parameters.

## Features

- RESTful API built with Golang Fiber framework
- Dynamic Docker image building using Docker-in-Docker
- Support for multiple model versions (Python Flask, Python FastAPI, Node.js)
- Real-time build status tracking
- Configurable build parameters and environment variables
- Health checks and monitoring
- Kubernetes-ready deployment

## Project Structure

```
.
├── main.go                           # Application entry point
├── go.mod                           # Go module definition
├── go.sum                           # Go dependencies
├── Dockerfile                       # API service Dockerfile
├── docker-compose.yml               # Local development setup
├── README.md                        # This file
├── internal/                        # Internal packages
│   ├── api/                         # API layer
│   │   ├── handlers.go             # HTTP handlers
│   │   └── routes.go               # Route definitions
│   ├── config/                     # Configuration
│   │   └── config.go               # Configuration management
│   ├── models/                     # Data models
│   │   └── models.go               # Request/Response models
│   └── services/                   # Business logic
│       └── docker.go               # Docker service implementation
└── utils/                          # Utility functions
    └── dockerfile_generator.go     # Dockerfile generation utilities
```

## API Endpoints

### Build Image

- **POST** `/api/v1/build`
- Build a new Docker image with specified parameters

**Request Body:**

```json
{
  "model_version": "python-flask",
  "version": "1.0.0",
  "name": "my-app",
  "tag": "latest",
  "description": "My application description",
  "environment": {
    "ENV_VAR": "value"
  },
  "build_args": {
    "BUILD_ARG": "value"
  }
}
```

**Response:**

```json
{
  "build_id": "uuid",
  "status": "pending",
  "message": "Build started",
  "image_name": "my-app:latest",
  "started_at": "2023-01-01T00:00:00Z"
}
```

### Get Build Status

- **GET** `/api/v1/build/{buildId}/status`
- Get the current status of a build

**Response:**

```json
{
  "build_id": "uuid",
  "status": "building",
  "message": "Building image...",
  "image_name": "my-app:latest",
  "started_at": "2023-01-01T00:00:00Z",
  "completed_at": null,
  "logs": ["Step 1/5 : FROM python:3.11-slim", "..."]
}
```

### List Builds

- **GET** `/api/v1/builds`
- List all builds

**Response:**

```json
{
  "builds": [...],
  "total": 5
}
```

### Health Check

- **GET** `/health`
- API health check

## Supported Model Versions

- `python-flask`: Python Flask application
- `python-fastapi`: Python FastAPI application
- `nodejs`: Node.js application

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DOCKER_HOST`: Docker daemon endpoint (default: unix:///var/run/docker.sock)

## Running Locally

### Using Docker Compose (Recommended)

```bash
# Start the services
docker-compose up -d

# Check logs
docker-compose logs -f api-builder

# Stop services
docker-compose down
```

### Using Go directly

```bash
# Install dependencies
go mod download

# Run the application
go run main.go
```

## Testing the API

### Build a Python Flask image

```bash
curl -X POST http://localhost:8080/api/v1/build \
  -H "Content-Type: application/json" \
  -d '{
    "model_version": "python-flask",
    "version": "1.0.0",
    "name": "my-flask-app",
    "tag": "latest",
    "description": "My Flask application",
    "environment": {
      "FLASK_ENV": "production"
    }
  }'
```

### Check build status

```bash
curl http://localhost:8080/api/v1/build/{BUILD_ID}/status
```

### List all builds

```bash
curl http://localhost:8080/api/v1/builds
```

## Kubernetes Deployment

The application is designed to run in Kubernetes with Docker-in-Docker support. Make sure to:

1. Use a privileged security context for DinD
2. Mount the Docker socket or use Docker-in-Docker sidecar
3. Configure appropriate RBAC permissions
4. Set resource limits and requests

## Build Statuses

- `pending`: Build is queued
- `building`: Build is in progress
- `success`: Build completed successfully
- `failed`: Build failed

## Error Handling

The API returns structured error responses:

```json
{
  "error": "Error description",
  "code": 400,
  "details": "Additional error details"
}
```

## Logging

The application uses structured JSON logging with different log levels. Logs include:

- HTTP request/response logging
- Build status updates
- Error tracking
- Performance metrics

## Security Considerations

When deploying in production:

- Use proper authentication/authorization
- Implement rate limiting
- Secure Docker daemon access
- Use network policies in Kubernetes
- Regularly update base images
- Scan images for vulnerabilities

# docker run

docker run --rm --name=api-docker-build --privileged -p 8080:8080 api-docker-builder:latest

