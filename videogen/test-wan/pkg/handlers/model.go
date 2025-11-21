package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wan2-video-server/pkg/model"
)

// ModelInfoHandler handles model information requests
type ModelInfoHandler struct {
	engine model.Engine
}

// NewModelInfoHandler creates a new model info handler
func NewModelInfoHandler(engine model.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		info := engine.GetModelInfo()
		c.JSON(http.StatusOK, info)
	}
}

// JobHandler handles job status requests
type JobHandler struct {
	engine model.Engine
}

// NewJobHandler creates a new job handler
func NewJobHandler(engine model.Engine) *JobHandler {
	return &JobHandler{engine: engine}
}

// GetJobStatus returns the status of a job
func (h *JobHandler) GetJobStatus(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	status := h.engine.GetJobStatus(jobID)
	if status == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}
