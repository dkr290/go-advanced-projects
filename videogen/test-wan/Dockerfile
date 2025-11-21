# Docker support (multi-stage build)
FROM nvidia/cuda:11.8.0-cudnn8-runtime-ubuntu22.04 AS python-base

# Install Python and dependencies
RUN apt-get update && apt-get install -y \
    python3.10 \
    python3-pip \
    python3-venv \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /app

# Copy Python backend
COPY python_backend/ /app/python_backend/

# Setup Python environment
RUN cd python_backend && \
    python3 -m venv venv && \
    . venv/bin/activate && \
    pip install --upgrade pip && \
    pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118 && \
    pip install -r requirements.txt

# Go builder stage
FROM golang:1.21-alpine AS go-builder

WORKDIR /build

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o wan2-video-server main.go

# Final stage
FROM nvidia/cuda:11.8.0-cudnn8-runtime-ubuntu22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    python3.10 \
    python3-pip \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Python environment from python-base
COPY --from=python-base /app/python_backend /app/python_backend

# Copy Go binary
COPY --from=go-builder /build/wan2-video-server /app/

# Copy configuration
COPY .env.example /app/.env

# Create necessary directories
RUN mkdir -p /app/uploads /app/outputs /app/models

# Expose ports
EXPOSE 8080 5000

# Create startup script
RUN echo '#!/bin/bash\n\
cd /app/python_backend && . venv/bin/activate && python server.py &\n\
sleep 5\n\
cd /app && ./wan2-video-server\n\
' > /app/start.sh && chmod +x /app/start.sh

CMD ["/app/start.sh"]
