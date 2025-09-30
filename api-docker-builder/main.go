package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/api-docker-builder/internal/api"
	"github.com/dkr290/go-advanced-projects/api-docker-builder/internal/config"
	"github.com/dkr290/go-advanced-projects/api-docker-builder/internal/logging"
	"github.com/dkr290/go-advanced-projects/api-docker-builder/internal/services"
)

func main() {
	// Initialize configuration
	cfg := config.Load()
	clog := logging.Init(false)

	// Initialize Docker service
	dockerService, err := services.NewDockerService(clog)
	if err != nil {
		clog.Fatalf("Failed to initialize Docker service: %v", err)
	}
	defer dockerService.Close()

	router := http.NewServeMux()
	// initialize Huma API for openapi
	humaAPI := humago.New(router, huma.DefaultConfig("Docker Builder API", "1.0.0"))

	// Register Huma-documented routes
	api.RegisterHumaRoutes(humaAPI, dockerService, clog)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	clog.Infof("Starting server on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Error starting the server", err)
	}
}
