// Package models is for build parameters and responce of the api
package models

import (
	"time"
)

// BuildRequest represents the request to build a Docker image
type BuildRequest struct {
	ModelVersion string             `json:"model_version"         validate:"required"`
	Version      string             `json:"version"               validate:"required"`
	Name         string             `json:"name"                  validate:"required"`
	Tag          string             `json:"tag"                   validate:"required"`
	Description  string             `json:"description"`
	Environment  map[string]string  `json:"environment,omitempty"`
	BuildArgs    map[string]*string `json:"build_args,omitempty"`
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
