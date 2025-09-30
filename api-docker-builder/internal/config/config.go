// Package config is container the configuration for the builder and also port env var
package config

import (
	"os"
)

type Config struct {
	Port          string
	DockerHost    string
	BuildTimeout  int
	MaxConcurrent int
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", ":8080"),
		DockerHost:    getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
		BuildTimeout:  60, // seconds
		MaxConcurrent: 5,  // max concurrent builds
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
