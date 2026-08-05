# duckduckgo-mcp

A Go MCP (Model Context Protocol) server exposing a `duckduckgo_search` tool
that lets an LLM search the internet via DuckDuckGo. It speaks MCP over the
**streamable HTTP** transport, so it can run as a standalone Docker
container and be reached at a plain URL such as `http://<host>:8080/mcp` —
suitable for wiring up to a model served by [LocalAI](https://localai.io/).

## Project layout

- `main.go` - HTTP server bootstrap: wires config, the DuckDuckGo client,
  and the MCP tool, then serves the streamable HTTP handler.
- `internal/handlers/` - Thin MCP tool handler (`duckduckgo_search`).
- `pkg/utils/` - DuckDuckGo HTML-endpoint search client (reusable domain
  logic).
- `config/` - Environment-variable driven configuration.

## Configuration (environment variables)

| Variable              | Default        | Description                                   |
|-----------------------|----------------|------------------------------------------------|
| `MCP_ADDR`             | `:8080`        | HTTP listen address                            |
| `MCP_PATH`             | `/mcp`         | Path the MCP handler is mounted on             |
| `MCP_REQUEST_TIMEOUT`  | `15s`          | Timeout for outbound DuckDuckGo requests       |
| `MCP_USER_AGENT`       | (generic UA)   | User-Agent sent to DuckDuckGo                  |
| `MCP_MAX_RESULTS`      | `10`           | Default number of results per search           |
| `MCP_LOG_LEVEL`        | `info`         | `debug`, `info`, `warn`, or `error`             |

## Tool: `duckduckgo_search`

Input:

```json
{ "query": "golang generics", "max_results": 5 }
```

Output:

```json
{
  "query": "golang generics",
  "count": 3,
  "results": [
    { "title": "...", "url": "https://...", "snippet": "..." }
  ]
}
```

## Run locally

```bash
go build -o duckduckgo-mcp .
./duckduckgo-mcp
# -> serving MCP at http://localhost:8080/mcp, health check at /healthz
```

## Build & run with Docker

```bash
# vendor/ must be present and up to date before building (offline build):
go mod vendor

docker build -t duckduckgo-mcp:latest .
docker run -d --name duckduckgo-mcp -p 8080:8080 duckduckgo-mcp:latest

curl http://localhost:8080/healthz
```

> The Dockerfile bakes in an optional corporate root CA
> (`zscaler-root-ca.crt`) so outbound HTTPS calls to DuckDuckGo succeed
> behind TLS-inspecting proxies. Remove the corresponding `COPY`/`RUN`
> lines in the Dockerfile if that's not needed in your environment.

## Wiring up to LocalAI

Point your LocalAI MCP-enabled model configuration at the container's
`/mcp` endpoint, e.g.:

```yaml
mcp:
  servers:
    duckduckgo:
      url: http://duckduckgo-mcp:8080/mcp
```

(Adjust the host/port to match how the container is reached from your
LocalAI deployment — e.g. a Docker Compose service name, or the host's
external IP if running as a separate container: `http://x.x.x.x:8080/mcp`.)

Once connected, the model can call the `duckduckgo_search` tool to look up
any topic on the internet.

## Testing

```bash
go test ./...
```

`main_test.go` contains an integration test that spins up the streamable
HTTP handler in-process, connects a real MCP client, and calls
`duckduckgo_search` end-to-end against the live DuckDuckGo endpoint (skipped
with `go test -short`).
