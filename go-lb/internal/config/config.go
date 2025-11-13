// Package config - doing all the configuration or fetching environment variables if they are present
package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	FrontIP        string
	BackendURLs    []string
	Weights        []int
	Percentages    []int
	Algorithm      string
	HealthPath     string
	HealthInterval string
	Mode           string
	DebugLog       string
}

func Load() *Config {
	var (
		debugLog = flag.String("debuglog", "false", "The debug logging")
		frontIP  = flag.String("front", "0.0.0.0:8080", "Front listen IP:port")
		backends = flag.String(
			"backends",
			"https://search.brave.com/,https://duckduckgo.com/",
			"Comma-separated backend URLs or IP:port (http(s)://host[:port]/path or host:port)",
		)

		weights     = flag.String("weights", "1,1", "Comma-separated weights for backends")
		percentages = flag.String(
			"percentages",
			"50,50",
			"Comma-separated percentages for backends",
		)
		algo = flag.String("algo", "roundrobin", "Algorithm: roundrobin|weighted|percentage")
	)
	healthPath := flag.String("health-path", "/", "Health check path (e.g. /health)")
	healthInterval := flag.String("health-interval", "5s", "Health check interval (e.g. 5s, 10s)")
	mode := flag.String("mode", "http", "Proxy mode: http or tcp")
	flag.Parse()

	backendRaw := strings.Split(getEnvOrDefault("LB_BACKENDS", *backends), ",")
	backendURLs := make([]string, 0, len(backendRaw))
	for _, b := range backendRaw {
		b = strings.TrimSpace(b)
		if b != "" {
			backendURLs = append(backendURLs, b)
		}
	}
	cfg := &Config{
		FrontIP:        getEnvOrDefault("LB_FRONT", *frontIP),
		BackendURLs:    backendURLs,
		Algorithm:      getEnvOrDefault("LB_ALGO", *algo),
		HealthPath:     getEnvOrDefault("LB_HEALTH_PATH", *healthPath),
		HealthInterval: getEnvOrDefault("LB_HEALTH_INTERVAL", *healthInterval),
		Mode:           getEnvOrDefault("LB_MODE", *mode),
		DebugLog:       getEnvOrDefault("DEBUG_LOG", *debugLog),
	}
	cfg.Weights = parseIntSlice(getEnvOrDefault("LB_WEIGHTS", *weights))
	cfg.Percentages = parseIntSlice(getEnvOrDefault("LB_PERCENTAGES", *percentages))
	if *algo != "roundrobin" && *algo != "weighted" && *algo != "percentage" {
		log.Fatal("wrong algorithm")
	}
	return cfg
}

func getEnvOrDefault(env, def string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return def
}

func parseIntSlice(s string) []int {
	parts := strings.Split(s, ",")
	res := make([]int, 0, len(parts))
	for _, p := range parts {
		var v int
		fmt.Sscanf(strings.TrimSpace(p), "%d", &v)
		res = append(res, v)
	}
	return res
}
