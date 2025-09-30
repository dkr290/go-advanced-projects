// Package utils  for various utils and functions
package utils

import (
	"fmt"
	"strings"

	"github.com/dkr290/go-advanced-projects/api-builder/internal/models"
)

// GeneratePythonFlaskDockerfile generates a Dockerfile for Python Flask applications
func GeneratePythonFlaskDockerfile(req *models.BuildImageRequest) string {
	var dockerfile strings.Builder

	dockerfile.WriteString("FROM python:3.11-slim\n\n")
	dockerfile.WriteString("# Set working directory\n")
	dockerfile.WriteString("WORKDIR /app\n\n")

	dockerfile.WriteString("# Install system dependencies\n")
	dockerfile.WriteString("RUN apt-get update && apt-get install -y \\\n")
	dockerfile.WriteString("    gcc \\\n")
	dockerfile.WriteString("    && rm -rf /var/lib/apt/lists/*\n\n")

	dockerfile.WriteString("# Copy requirements and install Python dependencies\n")
	dockerfile.WriteString("COPY requirements.txt .\n")
	dockerfile.WriteString("RUN pip install --no-cache-dir -r requirements.txt\n\n")

	dockerfile.WriteString("# Copy application code\n")
	dockerfile.WriteString("COPY . .\n\n")

	dockerfile.WriteString("# Add labels\n")
	dockerfile.WriteString(fmt.Sprintf("LABEL version=\"%s\"\n", req.Version))
	dockerfile.WriteString(fmt.Sprintf("LABEL description=\"%s\"\n", req.Description))
	dockerfile.WriteString(fmt.Sprintf("LABEL model_version=\"%s\"\n", req.ModelVersion))
	dockerfile.WriteString("\n")

	dockerfile.WriteString("# Expose port\n")
	dockerfile.WriteString("EXPOSE 5000\n\n")

	dockerfile.WriteString("# Health check\n")
	dockerfile.WriteString(
		"HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\\n",
	)
	dockerfile.WriteString("  CMD curl -f http://localhost:5000/health || exit 1\n\n")

	dockerfile.WriteString("# Run the application\n")
	dockerfile.WriteString("CMD [\"python\", \"app.py\"]\n")

	return dockerfile.String()
}

// GeneratePythonFastAPIDockerfile generates a Dockerfile for Python FastAPI applications
func GeneratePythonFastAPIDockerfile(req *models.BuildImageRequest) string {
	var dockerfile strings.Builder

	dockerfile.WriteString("FROM python:3.11-slim\n\n")
	dockerfile.WriteString("# Set working directory\n")
	dockerfile.WriteString("WORKDIR /app\n\n")

	dockerfile.WriteString("# Install system dependencies\n")
	dockerfile.WriteString("RUN apt-get update && apt-get install -y \\\n")
	dockerfile.WriteString("    gcc \\\n")
	dockerfile.WriteString("    && rm -rf /var/lib/apt/lists/*\n\n")

	dockerfile.WriteString("# Copy requirements and install Python dependencies\n")
	dockerfile.WriteString("COPY requirements.txt .\n")
	dockerfile.WriteString("RUN pip install --no-cache-dir -r requirements.txt\n\n")

	dockerfile.WriteString("# Copy application code\n")
	dockerfile.WriteString("COPY . .\n\n")

	dockerfile.WriteString("# Add labels\n")
	dockerfile.WriteString(fmt.Sprintf("LABEL version=\"%s\"\n", req.Version))
	dockerfile.WriteString(fmt.Sprintf("LABEL description=\"%s\"\n", req.Description))
	dockerfile.WriteString(fmt.Sprintf("LABEL model_version=\"%s\"\n", req.ModelVersion))
	dockerfile.WriteString("\n")

	dockerfile.WriteString("# Expose port\n")
	dockerfile.WriteString("EXPOSE 8000\n\n")

	dockerfile.WriteString("# Health check\n")
	dockerfile.WriteString(
		"HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\\n",
	)
	dockerfile.WriteString("  CMD curl -f http://localhost:8000/health || exit 1\n\n")

	dockerfile.WriteString("# Run the application\n")
	dockerfile.WriteString(
		"CMD [\"uvicorn\", \"main:app\", \"--host\", \"0.0.0.0\", \"--port\", \"8000\"]\n",
	)

	return dockerfile.String()
}

// GenerateNodeJSDockerfile generates a Dockerfile for Node.js applications
func GenerateNodeJSDockerfile(req *models.BuildImageRequest) string {
	var dockerfile strings.Builder

	dockerfile.WriteString("FROM node:18-alpine\n\n")
	dockerfile.WriteString("# Set working directory\n")
	dockerfile.WriteString("WORKDIR /app\n\n")

	dockerfile.WriteString("# Copy package files\n")
	dockerfile.WriteString("COPY package*.json ./\n\n")

	dockerfile.WriteString("# Install dependencies\n")
	dockerfile.WriteString("RUN npm ci --only=production\n\n")

	dockerfile.WriteString("# Copy application code\n")
	dockerfile.WriteString("COPY . .\n\n")

	dockerfile.WriteString("# Add labels\n")
	dockerfile.WriteString(fmt.Sprintf("LABEL version=\"%s\"\n", req.Version))
	dockerfile.WriteString(fmt.Sprintf("LABEL description=\"%s\"\n", req.Description))
	dockerfile.WriteString(fmt.Sprintf("LABEL model_version=\"%s\"\n", req.ModelVersion))
	dockerfile.WriteString("\n")

	dockerfile.WriteString("# Create non-root user\n")
	dockerfile.WriteString("RUN addgroup -g 1001 -S nodejs\n")
	dockerfile.WriteString("RUN adduser -S nextjs -u 1001\n")
	dockerfile.WriteString("USER nextjs\n\n")

	dockerfile.WriteString("# Expose port\n")
	dockerfile.WriteString("EXPOSE 3000\n\n")

	dockerfile.WriteString("# Health check\n")
	dockerfile.WriteString(
		"HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\\n",
	)
	dockerfile.WriteString("  CMD curl -f http://localhost:3000/health || exit 1\n\n")

	dockerfile.WriteString("# Run the application\n")
	dockerfile.WriteString("CMD [\"npm\", \"start\"]\n")

	return dockerfile.String()
}
