package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"duckduckgo-mcp/internal/handlers"
	"duckduckgo-mcp/pkg/utils"
)

// TestDuckDuckGoSearchToolOverHTTP starts the MCP streamable HTTP handler in
// an httptest server, connects a real MCP client to it, and invokes the
// duckduckgo_search tool end-to-end against the live DuckDuckGo HTML
// endpoint. It is skipped in short mode / offline CI environments.
func TestDuckDuckGoSearchToolOverHTTP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent integration test in short mode")
	}

	ddgClient := utils.NewDuckDuckGoClient(&http.Client{Timeout: 15 * time.Second}, "duckduckgo-mcp-test/1.0")
	searchHandler := handlers.NewSearchHandler(ddgClient, 5)

	getServer := func(_ *http.Request) *mcp.Server {
		server := mcp.NewServer(&mcp.Implementation{Name: "duckduckgo-mcp-test", Version: "test"}, nil)
		mcp.AddTool(server, &mcp.Tool{
			Name:        "duckduckgo_search",
			Description: "Search the internet using DuckDuckGo",
		}, searchHandler.Handle)
		return server
	}

	mcpHandler := mcp.NewStreamableHTTPHandler(getServer, &mcp.StreamableHTTPOptions{Stateless: true})

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: ts.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer session.Close()

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "duckduckgo_search",
		Arguments: map[string]any{
			"query":       "golang programming language",
			"max_results": 3,
		},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}
	if result.IsError {
		t.Fatalf("tool returned error result: %+v", result.Content)
	}
	if len(result.Content) == 0 {
		t.Fatalf("expected non-empty content in tool result")
	}
}
