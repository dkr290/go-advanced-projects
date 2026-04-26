package main

import (
	"fmt"
	"os"
	"strconv"
)


type Config struct {
	ServerPort int
}

// LoadConfig reads environment variables with sensible defaults.
// Designed for Docker/K8s — no .env file required.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerPort: 50051, // default gRPC port
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT %q: %w", v, err)
		}
		cfg.ServerPort = port
	}

	return cfg, nil
}

