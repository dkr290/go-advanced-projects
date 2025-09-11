package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(api huma.API, dockerService *services.DockerService) {
	// API v1 routes
	handlers := NewHandlers(dockerService)

	huma.Post(api, "/api/v1/build", handlers.BuildImage)

	// Build routes
	v1.Get("/build/:buildId/status", handlers.GetBuildStatus)
	v1.Get("/builds", handlers.ListBuilds)

	// Documentation route
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":        "Docker Image Builder API",
			"version":     "1.0.0",
			"description": "API for building Docker images with different model versions",
			"endpoints": fiber.Map{
				"POST /api/v1/build":                "Build a new Docker image",
				"GET /api/v1/build/:buildId/status": "Get build status",
				"GET /api/v1/builds":                "List all builds",
				"GET /health":                       "Health check",
			},
			"supported_model_versions": []string{
				"python-flask",
				"python-fastapi",
				"nodejs",
			},
		})
	})
}

