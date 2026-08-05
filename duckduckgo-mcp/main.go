// Command duckduckgo-mcp runs an MCP server exposing a DuckDuckGo web
// search tool over the streamable HTTP transport, so it can be reached at
// http://<host>:<port><path> (e.g. http://0.0.0.0:8080/mcp) from a
// LocalAI-managed model or any other MCP-compatible client.
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"duckduckgo-mcp/config"
	"duckduckgo-mcp/internal/handlers"
	"duckduckgo-mcp/pkg/utils"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))
	slog.SetDefault(logger)

	ddgClient := utils.NewDuckDuckGoClient(&http.Client{Timeout: cfg.RequestTimeout}, cfg.UserAgent)
	searchHandler := handlers.NewSearchHandler(ddgClient, cfg.MaxResults)

	getServer := func(_ *http.Request) *mcp.Server {
		server := mcp.NewServer(
			&mcp.Implementation{Name: "duckduckgo-mcp", Version: "1.0.0"},
			&mcp.ServerOptions{Logger: logger},
		)

		mcp.AddTool(server, &mcp.Tool{
			Name:        "duckduckgo_search",
			Description: "Search the internet for any topic using DuckDuckGo and return a list of relevant results (title, url, snippet).",
		}, searchHandler.Handle)

		return server
	}

	mcpHandler := mcp.NewStreamableHTTPHandler(getServer, &mcp.StreamableHTTPOptions{
		Stateless: true,
	})

	mux := http.NewServeMux()
	mux.Handle(cfg.MCPPath, mcpHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("starting duckduckgo-mcp server", "addr", cfg.Addr, "path", cfg.MCPPath)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down duckduckgo-mcp server")
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
