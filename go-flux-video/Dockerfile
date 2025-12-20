# Multi-stage build for Stable Diffusion Generator
# This Dockerfile builds a containerized version with CUDA support

# ============================================
# Stage 1: Builder - Compile everything
# ============================================
FROM golang:1.25-bookworm AS builder

# Prevent interactive prompts during build
ENV DEBIAN_FRONTEND=noninteractive


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o sd-generator main.go





# ============================================
# Stage 2: Runtime - Minimal production image
# ============================================

#FROM nvidia/cuda:12.4.1-devel-ubuntu22.04
FROM nvidia/cuda:12.4.1-runtime-ubuntu22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1
# Install Python and runtime dependencies
RUN apt-get update && apt-get install -y  \
    python3 \
    python3-pip \
    python3-venv \
    ca-certificates \
    openssh-server \
    vim \
    nano \
    htop \
    && rm -rf /var/lib/apt/lists/* 




WORKDIR /app

COPY --from=builder /app/sd-generator /app/

# Copy Python script
COPY scripts/python_generate.py /app/scripts/

# Install Python dependencies
RUN pip3 install --no-cache-dir \
    torch \
    diffusers \
    transformers \
    accelerate \
    safetensors \
    sentencepiece \
    protobuf \
    pillow \
    xformers \
    hf_transfer \
    gguf>=0.10.0 \
    peft  \
    && rm -rf ~/.cache/pip

# Create directories for runtime data
RUN mkdir -p /app/models/lora /app/output

# Set library path for CUDA
ENV LD_LIBRARY_PATH="/usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}"


# SSH for Runpod
RUN mkdir -p /run/sshd && \
    echo 'root:root' | chpasswd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
# Expose SSH port
EXPOSE 22

# Start SSH server in foreground for Runpod
CMD ["/usr/sbin/sshd", "-D"]

ENV HF_HOME=/app/models/.cache
# Labels for documentation
LABEL maintainer="your-email@example.com"
LABEL description="FLUX Stable Diffusion image generator with CUDA support"
LABEL version="1.0"

# Expected volume mounts:
# - /app/character_config.json (config file, read-only recommended)
# - /app/models (model cache, persistent storage)
# - /app/output (generated images)
