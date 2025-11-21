package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/model"
	"github.com/wan2-video-server/pkg/types"
	"github.com/wan2-video-server/pkg/utils"
)

// VideoHandler handles video generation requests
type VideoHandler struct {
	engine model.Engine
	config *config.Config
	log    *logger.Logger
}

// NewVideoHandler creates a new video handler
func NewVideoHandler(engine model.Engine, cfg *config.Config) *VideoHandler {
	return &VideoHandler{
		engine: engine,
		config: cfg,
		log:    logger.NewLogger(),
	}
}

// TextToVideo handles text-to-video generation
func (h *VideoHandler) TextToVideo(c *gin.Context) {
	var req types.TextToVideoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	// Validate request
	if req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	// Set defaults
	if req.NumFrames == 0 {
		req.NumFrames = 64
	}
	if req.FPS == 0 {
		req.FPS = h.config.Model.DefaultFPS
	}
	if req.Width == 0 {
		req.Width = h.config.Model.DefaultWidth
	}
	if req.Height == 0 {
		req.Height = h.config.Model.DefaultHeight
	}

	h.log.Infof("Text-to-video request: prompt=%s, frames=%d, fps=%d", req.Prompt, req.NumFrames, req.FPS)

	// Generate job ID
	jobID := utils.GenerateJobID()

	// Start generation asynchronously
	go func() {
		params := &types.GenerationParams{
			Prompt:      req.Prompt,
			NegativePrompt: req.NegativePrompt,
			NumFrames:   req.NumFrames,
			FPS:         req.FPS,
			Width:       req.Width,
			Height:      req.Height,
			Seed:        req.Seed,
			GuidanceScale: req.GuidanceScale,
			NumInferenceSteps: req.NumInferenceSteps,
		}

		result, err := h.engine.GenerateTextToVideo(params)
		if err != nil {
			h.log.Errorf("Generation failed for job %s: %v", jobID, err)
			h.engine.UpdateJobStatus(jobID, "failed", err.Error())
			return
		}

		h.engine.UpdateJobStatus(jobID, "completed", result.OutputPath)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":  jobID,
		"status":  "processing",
		"message": "Video generation started",
	})
}

// ImageToVideo handles image-to-video generation
func (h *VideoHandler) ImageToVideo(c *gin.Context) {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get image file
	files := form.File["image"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	imageFile := files[0]

	// Validate file size
	if imageFile.Size > h.config.Process.UploadMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds maximum allowed"})
		return
	}

	// Save uploaded file
	uploadDir := "./uploads"
	utils.EnsureDir(uploadDir)
	
	timestamp := time.Now().Unix()
	imagePath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", timestamp, imageFile.Filename))
	
	if err := c.SaveUploadedFile(imageFile, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Get parameters
	prompt := c.PostForm("prompt")
	negativePrompt := c.PostForm("negative_prompt")
	
	numFrames := utils.ParseIntOrDefault(c.PostForm("num_frames"), 64)
	fps := utils.ParseIntOrDefault(c.PostForm("fps"), h.config.Model.DefaultFPS)
	width := utils.ParseIntOrDefault(c.PostForm("width"), h.config.Model.DefaultWidth)
	height := utils.ParseIntOrDefault(c.PostForm("height"), h.config.Model.DefaultHeight)
	seed := utils.ParseInt64OrDefault(c.PostForm("seed"), -1)
	guidanceScale := utils.ParseFloat64OrDefault(c.PostForm("guidance_scale"), 7.5)
	numInferenceSteps := utils.ParseIntOrDefault(c.PostForm("num_inference_steps"), 50)

	h.log.Infof("Image-to-video request: image=%s, prompt=%s", imagePath, prompt)

	// Generate job ID
	jobID := utils.GenerateJobID()

	// Start generation asynchronously
	go func() {
		params := &types.GenerationParams{
			Prompt:         prompt,
			NegativePrompt: negativePrompt,
			ImagePath:      imagePath,
			NumFrames:      numFrames,
			FPS:            fps,
			Width:          width,
			Height:         height,
			Seed:           seed,
			GuidanceScale:  guidanceScale,
			NumInferenceSteps: numInferenceSteps,
		}

		result, err := h.engine.GenerateImageToVideo(params)
		if err != nil {
			h.log.Errorf("Generation failed for job %s: %v", jobID, err)
			h.engine.UpdateJobStatus(jobID, "failed", err.Error())
			return
		}

		h.engine.UpdateJobStatus(jobID, "completed", result.OutputPath)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":  jobID,
		"status":  "processing",
		"message": "Video generation started",
	})
}

// VideoToVideo handles video-to-video generation
func (h *VideoHandler) VideoToVideo(c *gin.Context) {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get video file
	files := form.File["video"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file is required"})
		return
	}

	videoFile := files[0]

	// Validate file size
	if videoFile.Size > h.config.Process.UploadMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds maximum allowed"})
		return
	}

	// Save uploaded file
	uploadDir := "./uploads"
	utils.EnsureDir(uploadDir)
	
	timestamp := time.Now().Unix()
	videoPath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", timestamp, videoFile.Filename))
	
	if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
		return
	}

	// Get parameters
	prompt := c.PostForm("prompt")
	negativePrompt := c.PostForm("negative_prompt")
	
	fps := utils.ParseIntOrDefault(c.PostForm("fps"), h.config.Model.DefaultFPS)
	strength := utils.ParseFloat64OrDefault(c.PostForm("strength"), 0.8)
	seed := utils.ParseInt64OrDefault(c.PostForm("seed"), -1)
	guidanceScale := utils.ParseFloat64OrDefault(c.PostForm("guidance_scale"), 7.5)
	numInferenceSteps := utils.ParseIntOrDefault(c.PostForm("num_inference_steps"), 50)

	h.log.Infof("Video-to-video request: video=%s, prompt=%s", videoPath, prompt)

	// Generate job ID
	jobID := utils.GenerateJobID()

	// Start generation asynchronously
	go func() {
		params := &types.GenerationParams{
			Prompt:         prompt,
			NegativePrompt: negativePrompt,
			VideoPath:      videoPath,
			FPS:            fps,
			Strength:       strength,
			Seed:           seed,
			GuidanceScale:  guidanceScale,
			NumInferenceSteps: numInferenceSteps,
		}

		result, err := h.engine.GenerateVideoToVideo(params)
		if err != nil {
			h.log.Errorf("Generation failed for job %s: %v", jobID, err)
			h.engine.UpdateJobStatus(jobID, "failed", err.Error())
			return
		}

		h.engine.UpdateJobStatus(jobID, "completed", result.OutputPath)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"job_id":  jobID,
		"status":  "processing",
		"message": "Video generation started",
	})
}
