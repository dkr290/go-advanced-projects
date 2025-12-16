# Multi-stage build for Stable Diffusion Generator
# This Dockerfile builds a containerized version with CUDA support

# ============================================
# Stage 1: Builder - Compile everything
# ============================================
FROM  nvidia/cuda:12.4.1-devel-ubuntu22.04 AS builder

# Prevent interactive prompts during build
ENV DEBIAN_FRONTEND=noninteractive

# Install build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    wget \
    ca-certificates \
    libgomp1 \
    pkg-config \
    libcurl4-openssl-dev \
    ocl-icd-opencl-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Go
ARG GO_VERSION=1.25.5
RUN rm -rf /usr/local/go /usr/local/go1.* /go /root/go && \
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOROOT="/usr/local/go"

WORKDIR /build


# Clone and build stable-diffusion.cpp with CUDA support
RUN git clone --recursive https://github.com/Binozo/GoStableDiffusion.git && \
    cd GoStableDiffusion && \
    git submodule update --init --recursive && \
    cd stable-diffusion.cpp && \
    mkdir build && cd build && \
    cmake ..  -DSD_CUDA=ON  -DSD_BUILD_SHARED_LIBS=ON && \
    cmake --build . --config Release -j$(nproc) && \
    cp bin/libstable-diffusion.so /usr/local/lib/ && \
    cp bin/sd /usr/local/lib/ && \
    ldconfig -v


# Copy Go application source
WORKDIR /app
COPY . .
RUN go mod download


# Build Go binary with CUDA stubs linked so ld can resolve cu* at build time
ENV CGO_ENABLED=1
ENV CUDA_HOME=/usr/local/cuda
ENV LD_LIBRARY_PATH="${CUDA_HOME}/lib64:${LD_LIBRARY_PATH}"
# Use CUDA stubs for link; at runtime the NVIDIA driver (libcuda) is injected by the host (Runpod)
ENV CGO_LDFLAGS="-L/usr/local/cuda/lib64/stubs -lcuda -Wl,-rpath,/usr/local/cuda/lib64 -Wl,-rpath-link,/usr/local/cuda/lib64"
RUN go build -o sd-generator main.go


# ============================================ 
# Stage 2: Runtime - Minimal production image
# ============================================

FROM nvidia/cuda:12.4.1-devel-ubuntu22.04 

ENV DEBIAN_FRONTEND=noninteractive
# Install only runtime dependencies
RUN apt-get update && apt-get install -y  \
    ca-certificates \
    libgomp1 \
    libcurl4 \
    ocl-icd-libopencl1 \ 
    openssh-server \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary and C++ library from builder stage
COPY --from=builder /app/sd-generator /app/
COPY --from=builder /usr/local/lib/libstable-diffusion.so /usr/local/lib/
COPY --from=builder /usr/local/lib/sd /usr/local/lib/
RUN ldconfig

# Create directories for runtime data
# Models and output will be mounted as volumes
RUN mkdir -p /app/models/lora /app/output

# Ensure libs are found (driver libs are injected by NVIDIA runtime on Runpod)
ENV LD_LIBRARY_PATH="/usr/local/lib:/usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}"

# SSH for Runpod
RUN mkdir -p /run/sshd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config
# Expose SSH port
EXPOSE 22

# Start SSH server in foreground for Runpod
CMD ["/usr/sbin/sshd", "-D"]

# Labels for documentation
LABEL maintainer="your-email@example.com"
LABEL description="FLUX Stable Diffusion image generator with CUDA support"
LABEL version="1.0"

# Expected volume mounts:
# - /app/character_config.json (config file, read-only recommended)
# - /app/models (model cache, persistent storage)
# - /app/output (generated images)
