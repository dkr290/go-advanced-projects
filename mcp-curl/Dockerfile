FROM golang:1.25.2 AS builder
WORKDIR /app
COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN apt-get update && \
  apt-get install -y ca-certificates && \
  update-ca-certificates
COPY main.go .
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
# Build a static binary for linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mcp-curl -ldflags "-X main.build=main" ./

FROM curlimages/curl:8.6.0
WORKDIR /app
COPY --from=builder /app/mcp-curl .
CMD ["./mcp-curl"]

