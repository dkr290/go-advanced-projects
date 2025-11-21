package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/model"
)

// ModelManagementHandler handles model management operations
type ModelManagementHandler struct {
	config *config.Config
	log    *logger.Logger
}

// NewModelManagementHandler creates a new model management handler
func NewModelManagementHandler(cfg *config.Config) *ModelManagementHandler {
	return &ModelManagementHandler{
		config: cfg,
		log:    logger.NewLogger(),
	}
}

// ListModels returns a list of available models
func (h *ModelManagementHandler) ListModels(c *gin.Context) {
	models := []gin.H{
		{
			"id":          h.config.Model.HuggingFaceModelID,
			"name":        h.config.Model.Name,
			"provider":    h.config.Model.Provider,
			"cache_dir":   h.config.Model.CacheDir,
			"downloaded":  h.isModelDownloaded(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// DownloadModel downloads a model from Hugging Face
func (h *ModelManagementHandler) DownloadModel(c *gin.Context) {
	var req struct {
		ModelID string `json:"model_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.ModelID == "" {
		req.ModelID = h.config.Model.HuggingFaceModelID
	}

	h.log.Infof("Starting model download: %s", req.ModelID)

	// Download in background
	go func() {
		downloader := model.NewHuggingFaceDownloader(h.config)
		if err := downloader.Download(); err != nil {
			h.log.Errorf("Failed to download model: %v", err)
		} else {
			h.log.Infof("Model downloaded successfully")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Model download started",
		"model_id": req.ModelID,
	})
}

func (h *ModelManagementHandler) isModelDownloaded() bool {
	modelPath := filepath.Join(h.config.Model.CacheDir, h.config.Model.HuggingFaceModelID)
	if _, err := os.Stat(modelPath); err == nil {
		return true
	}
	return false
}
