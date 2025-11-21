package model

import (
	"github.com/wan2-video-server/pkg/types"
)

// Engine defines the interface for video generation models
type Engine interface {
	// GenerateTextToVideo generates a video from text prompt
	GenerateTextToVideo(params *types.GenerationParams) (*types.GenerationResult, error)

	// GenerateImageToVideo generates a video from an image and prompt
	GenerateImageToVideo(params *types.GenerationParams) (*types.GenerationResult, error)

	// GenerateVideoToVideo generates a video from another video and prompt
	GenerateVideoToVideo(params *types.GenerationParams) (*types.GenerationResult, error)

	// GetModelInfo returns information about the model
	GetModelInfo() *types.ModelInfo

	// GetJobStatus returns the status of a job
	GetJobStatus(jobID string) *types.JobStatus

	// UpdateJobStatus updates the status of a job
	UpdateJobStatus(jobID, status, message string) error
}
