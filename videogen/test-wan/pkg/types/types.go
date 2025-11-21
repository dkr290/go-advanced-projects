package types

// TextToVideoRequest represents a text-to-video generation request
type TextToVideoRequest struct {
	Prompt            string  `json:"prompt" binding:"required"`
	NegativePrompt    string  `json:"negative_prompt"`
	NumFrames         int     `json:"num_frames"`
	FPS               int     `json:"fps"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	Seed              int64   `json:"seed"`
	GuidanceScale     float64 `json:"guidance_scale"`
	NumInferenceSteps int     `json:"num_inference_steps"`
}

// GenerationParams represents parameters for video generation
type GenerationParams struct {
	Prompt            string
	NegativePrompt    string
	ImagePath         string
	VideoPath         string
	NumFrames         int
	FPS               int
	Width             int
	Height            int
	Seed              int64
	GuidanceScale     float64
	NumInferenceSteps int
	Strength          float64
}

// GenerationResult represents the result of video generation
type GenerationResult struct {
	OutputPath string
	Duration   float64
	Frames     int
	FPS        int
}

// JobStatus represents the status of a generation job
type JobStatus struct {
	JobID      string `json:"job_id"`
	Status     string `json:"status"` // pending, processing, completed, failed
	Message    string `json:"message"`
	OutputPath string `json:"output_path,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ModelInfo represents information about the loaded model
type ModelInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Provider    string `json:"provider"`
	GPUEnabled  bool   `json:"gpu_enabled"`
	GPUDeviceID int    `json:"gpu_device_id,omitempty"`
	CacheDir    string `json:"cache_dir"`
}
