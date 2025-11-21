package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/types"
)

// LocalEngine implements the Engine interface for local model inference
// This is a placeholder for future direct Go implementation
type LocalEngine struct {
	config    *config.Config
	log       *logger.Logger
	jobs      map[string]*types.JobStatus
	jobsMutex sync.RWMutex
}

// NewLocalEngine creates a new local inference engine
func NewLocalEngine(cfg *config.Config) (Engine, error) {
	log := logger.NewLogger()
	log.Info("Initializing local inference engine")
	log.Warn("Local engine is not fully implemented. Consider using Python backend.")

	return &LocalEngine{
		config: cfg,
		log:    log,
		jobs:   make(map[string]*types.JobStatus),
	}, nil
}

// GenerateTextToVideo generates a video from text prompt
func (e *LocalEngine) GenerateTextToVideo(params *types.GenerationParams) (*types.GenerationResult, error) {
	return nil, fmt.Errorf("local engine not implemented - use Python backend")
}

// GenerateImageToVideo generates a video from an image and prompt
func (e *LocalEngine) GenerateImageToVideo(params *types.GenerationParams) (*types.GenerationResult, error) {
	return nil, fmt.Errorf("local engine not implemented - use Python backend")
}

// GenerateVideoToVideo generates a video from another video and prompt
func (e *LocalEngine) GenerateVideoToVideo(params *types.GenerationParams) (*types.GenerationResult, error) {
	return nil, fmt.Errorf("local engine not implemented - use Python backend")
}

// GetModelInfo returns information about the model
func (e *LocalEngine) GetModelInfo() *types.ModelInfo {
	return &types.ModelInfo{
		Name:        e.config.Model.Name,
		Version:     "1.0.0",
		Provider:    "local",
		GPUEnabled:  e.config.GPU.Enabled,
		GPUDeviceID: e.config.GPU.DeviceID,
		CacheDir:    e.config.Model.CacheDir,
	}
}

// GetJobStatus returns the status of a job
func (e *LocalEngine) GetJobStatus(jobID string) *types.JobStatus {
	e.jobsMutex.RLock()
	defer e.jobsMutex.RUnlock()

	if job, exists := e.jobs[jobID]; exists {
		return job
	}
	return nil
}

// UpdateJobStatus updates the status of a job
func (e *LocalEngine) UpdateJobStatus(jobID, status, message string) error {
	e.jobsMutex.Lock()
	defer e.jobsMutex.Unlock()

	now := time.Now().Format(time.RFC3339)

	if job, exists := e.jobs[jobID]; exists {
		job.Status = status
		job.Message = message
		job.UpdatedAt = now
	} else {
		e.jobs[jobID] = &types.JobStatus{
			JobID:     jobID,
			Status:    status,
			Message:   message,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	return nil
}
