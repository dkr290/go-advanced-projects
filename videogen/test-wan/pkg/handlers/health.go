package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "wan2-video-server",
		"version": "1.0.0",
	})
}

// Index handles the root endpoint
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Wan2.1 Video Generation Server",
		"version": "1.0.0",
		"endpoints": gin.H{
			"health":            "/health",
			"model_info":        "/api/v1/model/info",
			"text_to_video":     "/api/v1/generate/text-to-video",
			"image_to_video":    "/api/v1/generate/image-to-video",
			"video_to_video":    "/api/v1/generate/video-to-video",
			"job_status":        "/api/v1/job/:id",
			"list_models":       "/api/v1/models",
			"download_model":    "/api/v1/models/download",
		},
	})
}
