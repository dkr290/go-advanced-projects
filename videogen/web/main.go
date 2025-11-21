package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"videogen/web/handlers"
	"videogen/web/middleware"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Set Gin mode
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Setup middleware
	r.Use(middleware.CORS())

	// Load HTML templates
	r.LoadHTMLGlob("templates/**/*")

	// Static files
	r.Static("/static", "./static")
	r.Static("/assets", "./assets")

	// Routes
	setupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Web UI starting on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// Web pages
	r.GET("/", handlers.IndexHandler)
	r.GET("/text-to-video", handlers.TextToVideoHandler)
	r.GET("/image-to-video", handlers.ImageToVideoHandler)
	r.GET("/video-to-video", handlers.VideoToVideoHandler)
	r.GET("/gallery", handlers.GalleryHandler)
	r.GET("/settings", handlers.SettingsHandler)

	// API endpoints for HTMX
	api := r.Group("/api")
	{
		api.POST("/generate/text-to-video", handlers.GenerateTextToVideoHandler)
		api.POST("/generate/image-to-video", handlers.GenerateImageToVideoHandler)
		api.POST("/generate/video-to-video", handlers.GenerateVideoToVideoHandler)
		api.GET("/job/:id", handlers.GetJobStatusHandler)
		api.GET("/models", handlers.GetModelsHandler)
		api.POST("/switch-model", handlers.SwitchModelHandler)
		api.GET("/gallery/list", handlers.GetGalleryListHandler)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
}
