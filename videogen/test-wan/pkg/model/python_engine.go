package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/types"
)

// PythonEngine implements the Engine interface using a Python backend
type PythonEngine struct {
	config     *config.Config
	log        *logger.Logger
	httpClient *http.Client
	jobs       map[string]*types.JobStatus
	jobsMutex  sync.RWMutex
}

// NewPythonEngine creates a new Python-based engine
func NewPythonEngine(cfg *config.Config) (Engine, error) {
	log := logger.NewLogger()
	log.Info("Initializing Python backend engine")

	engine := &PythonEngine{
		config: cfg,
		log:    log,
		httpClient: &http.Client{
			Timeout: cfg.Python.Timeout,
		},
		jobs: make(map[string]*types.JobStatus),
	}

	// Test connection to Python backend
	if err := engine.testConnection(); err != nil {
		log.Warnf("Python backend not available: %v", err)
		log.Info("Make sure to start the Python backend server")
	}

	return engine, nil
}

func (e *PythonEngine) testConnection() error {
	resp, err := e.httpClient.Get(e.config.Python.URL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned status %d", resp.StatusCode)
	}

	return nil
}

// GenerateTextToVideo generates a video from text prompt
func (e *PythonEngine) GenerateTextToVideo(params *types.GenerationParams) (*types.GenerationResult, error) {
	e.log.Infof("Generating text-to-video: %s", params.Prompt)

	requestBody, err := json.Marshal(map[string]interface{}{
		"prompt":              params.Prompt,
		"negative_prompt":     params.NegativePrompt,
		"num_frames":          params.NumFrames,
		"fps":                 params.FPS,
		"width":               params.Width,
		"height":              params.Height,
		"seed":                params.Seed,
		"guidance_scale":      params.GuidanceScale,
		"num_inference_steps": params.NumInferenceSteps,
	})
	if err != nil {
		return nil, err
	}

	resp, err := e.httpClient.Post(
		e.config.Python.URL+"/api/generate/text-to-video",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend error: %s", string(body))
	}

	var result types.GenerationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateImageToVideo generates a video from an image and prompt
func (e *PythonEngine) GenerateImageToVideo(params *types.GenerationParams) (*types.GenerationResult, error) {
	e.log.Infof("Generating image-to-video: image=%s, prompt=%s", params.ImagePath, params.Prompt)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file
	file, err := os.Open(params.ImagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("image", filepath.Base(params.ImagePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	// Add parameters
	writer.WriteField("prompt", params.Prompt)
	writer.WriteField("negative_prompt", params.NegativePrompt)
	writer.WriteField("num_frames", fmt.Sprintf("%d", params.NumFrames))
	writer.WriteField("fps", fmt.Sprintf("%d", params.FPS))
	writer.WriteField("width", fmt.Sprintf("%d", params.Width))
	writer.WriteField("height", fmt.Sprintf("%d", params.Height))
	writer.WriteField("seed", fmt.Sprintf("%d", params.Seed))
	writer.WriteField("guidance_scale", fmt.Sprintf("%f", params.GuidanceScale))
	writer.WriteField("num_inference_steps", fmt.Sprintf("%d", params.NumInferenceSteps))

	writer.Close()

	req, err := http.NewRequest("POST", e.config.Python.URL+"/api/generate/image-to-video", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend error: %s", string(body))
	}

	var result types.GenerationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateVideoToVideo generates a video from another video and prompt
func (e *PythonEngine) GenerateVideoToVideo(params *types.GenerationParams) (*types.GenerationResult, error) {
	e.log.Infof("Generating video-to-video: video=%s, prompt=%s", params.VideoPath, params.Prompt)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add video file
	file, err := os.Open(params.VideoPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("video", filepath.Base(params.VideoPath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	// Add parameters
	writer.WriteField("prompt", params.Prompt)
	writer.WriteField("negative_prompt", params.NegativePrompt)
	writer.WriteField("fps", fmt.Sprintf("%d", params.FPS))
	writer.WriteField("strength", fmt.Sprintf("%f", params.Strength))
	writer.WriteField("seed", fmt.Sprintf("%d", params.Seed))
	writer.WriteField("guidance_scale", fmt.Sprintf("%f", params.GuidanceScale))
	writer.WriteField("num_inference_steps", fmt.Sprintf("%d", params.NumInferenceSteps))

	writer.Close()

	req, err := http.NewRequest("POST", e.config.Python.URL+"/api/generate/video-to-video", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backend error: %s", string(body))
	}

	var result types.GenerationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetModelInfo returns information about the model
func (e *PythonEngine) GetModelInfo() *types.ModelInfo {
	return &types.ModelInfo{
		Name:        e.config.Model.Name,
		Version:     "1.0.0",
		Provider:    e.config.Model.Provider,
		GPUEnabled:  e.config.GPU.Enabled,
		GPUDeviceID: e.config.GPU.DeviceID,
		CacheDir:    e.config.Model.CacheDir,
	}
}

// GetJobStatus returns the status of a job
func (e *PythonEngine) GetJobStatus(jobID string) *types.JobStatus {
	e.jobsMutex.RLock()
	defer e.jobsMutex.RUnlock()

	if job, exists := e.jobs[jobID]; exists {
		return job
	}
	return nil
}

// UpdateJobStatus updates the status of a job
func (e *PythonEngine) UpdateJobStatus(jobID, status, message string) error {
	e.jobsMutex.Lock()
	defer e.jobsMutex.Unlock()

	now := time.Now().Format(time.RFC3339)

	if job, exists := e.jobs[jobID]; exists {
		job.Status = status
		job.Message = message
		job.UpdatedAt = now
		if status == "completed" {
			job.OutputPath = message
		}
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
