// Package config centralizes all environment variable driven configuration
// for the duckduckgo-mcp server.
package config

import (
	"os"
	"strconv"
	"time"
)

// AppConfig holds all runtime configuration loaded from the environment.
type AppConfig struct {
	// Addr is the host:port the HTTP server listens on, e.g. ":8080".
	Addr string
	// MCPPath is the HTTP path the MCP streamable handler is mounted on,
	// e.g. "/mcp".
	MCPPath string
	// RequestTimeout bounds how long a single DuckDuckGo search request
	// is allowed to take.
	RequestTimeout time.Duration
	// UserAgent is sent on outgoing search requests to DuckDuckGo.
	UserAgent string
	// MaxResults caps the number of search results returned per query.
	MaxResults int
	// LogLevel controls slog verbosity: debug, info, warn, error.
	LogLevel string
}

// Load reads configuration from environment variables, applying sane
// defaults for anything that is not set.
func Load() AppConfig {
	return AppConfig{
		Addr:           getEnv("MCP_ADDR", ":8080"),
		MCPPath:        getEnv("MCP_PATH", "/mcp"),
		RequestTimeout: getEnvDuration("MCP_REQUEST_TIMEOUT", 15*time.Second),
		UserAgent:      getEnv("MCP_USER_AGENT", "Mozilla/5.0 (compatible; duckduckgo-mcp/1.0; +https://github.com/)"),
		MaxResults:     getEnvInt("MCP_MAX_RESULTS", 10),
		LogLevel:       getEnv("MCP_LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
