package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/services"
)

type RootInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Endpoints   map[string]string `json:"endpoints"`
	Supported   []string          `json:"supported_model_versions"`
}

func RegisterHumaRoutes(api huma.API, dockerService *services.DockerService) {
	// API v1 routes
	handlers := NewHandlers(dockerService)

	huma.Post(api, "/api/v1/build", handlers.BuildImage)

	// Build status and list
	huma.Get(api, "/api/v1/build/{buildId}/status", handlers.GetBuildStatus)
	huma.Get(api, "/api/v1/builds", handlers.ListBuilds)
	// Documentation route
	// convert the root info route to Huma as well
	huma.Get(api, "/api/v1", func(ctx context.Context, _ *struct{}) (*RootInfo, error) {
		return &RootInfo{
			Name:        "Docker Image Builder API",
			Version:     "1.0.0",
			Description: "API for building Docker images with different model versions",
			Endpoints: map[string]string{
				"POST /api/v1/build":                 "Build a new Docker image",
				"GET /api/v1/build/{buildId}/status": "Get build status",
				"GET /api/v1/builds":                 "List all builds",
				"GET /health":                        "Health check",
			},
			Supported: []string{"python-flask", "python-fastapi", "nodejs"},
		}, nil
	})
}
