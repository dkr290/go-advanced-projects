package main

import (
	"log"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/api"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/config"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Initialize Docker service
	dockerService, err := services.NewDockerService()
	if err != nil {
		logrus.Fatalf("Failed to initialize Docker service: %v", err)
	}
	defer dockerService.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logrus.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// initialize Huma API for openapi
	humaAPI := humafiber.New(app, huma.DefaultConfig("Docker Builder API", "1.0.0"))

	// Register Huma-documented routes
	api.RegisterHumaRoutes(humaAPI, dockerService)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "docker-image-builder",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	logrus.Infof("Starting server on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
