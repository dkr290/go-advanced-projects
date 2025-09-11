// Package models is for build parameters and responce of the api
package models

import (
	"time"
)

// BuildRequest represents the request to build a Docker image
type BuildRequest struct {
	ModelVersion string `json:"model_version" example:"python-flask"             enum:"python-flask,python-fastapi,nodejs" description:"Base template to use"`
	Version      string `json:"version"       example:"1.0.0"                                                              description:"Application version label"`
	Name         string `json:"name"          example:"myapp"                                                              description:"Image name"`
	Tag          string `json:"tag"           example:"latest"                                                             description:"Image tag"`
	Description  string `json:"description"   example:"Initial build for my app"`
}

// BuildResponse represents the response after initiating a build
type BuildResponse struct {
	BuildID   string    `json:"build_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	ImageName string    `json:"image_name"`
	StartedAt time.Time `json:"started_at"`
}

// BuildStatus represents the current status of a build
type BuildStatus struct {
	BuildID     string     `json:"build_id"`
	Status      string     `json:"status"` // pending, building, success, failed
	Message     string     `json:"message"`
	ImageName   string     `json:"image_name"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Logs        []string   `json:"logs,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}
