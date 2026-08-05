FROM golang:1.25-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

# Build a static binary for linux so it runs in the minimal runtime image.
# -mod=vendor avoids any network access during the build (useful behind
# proxies/TLS-inspecting firewalls).
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -ldflags="-s -w" -o /out/duckduckgo-mcp .

FROM debian:bookworm-slim
WORKDIR /app


COPY --from=builder /out/duckduckgo-mcp .

# HTTP address/path served by the MCP streamable handler. Override at
# runtime with `-e MCP_ADDR=:8080 -e MCP_PATH=/mcp` if needed.
ENV MCP_ADDR=:8080
ENV MCP_PATH=/mcp

EXPOSE 8080

USER mcp
ENTRYPOINT ["/app/duckduckgo-mcp"]
