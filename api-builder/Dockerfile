# Build stage
FROM golang:1.25 AS builder

COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN apt-get update && \
  apt-get install -y ca-certificates && \
  update-ca-certificates
# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM docker:28.4.0-dind


# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080


# Run the application
CMD ["./main"]
