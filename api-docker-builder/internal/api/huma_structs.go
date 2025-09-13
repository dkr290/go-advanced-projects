package api

import (
	"time"
)

type GetBuildStatus struct {
	BuildID string `path:"buildId" doc:"The Build ID to query"`
}
type BuildImageOutput struct {
	Body struct {
		BuildID   string    `json:"build_id"`
		Status    string    `json:"status"`
		Message   string    `json:"message"`
		ImageName string    `json:"image_name"`
		StartedAt time.Time `json:"started_at"`
	}
}
type GetBuildStatusOutput struct {
	Body struct {
		BuildID     string     `json:"build_id"`
		Status      string     `json:"status"` // pending, building, success, failed
		Message     string     `json:"message"`
		ImageName   string     `json:"image_name"`
		StartedAt   time.Time  `json:"started_at"`
		CompletedAt *time.Time `json:"completed_at,omitempty"`
		Logs        []string   `json:"logs,omitempty"`
	}
}

type BuildImageInput struct {
	Body struct {
		ModelVersion string `json:"model_version" example:"python-flask"             enum:"python-flask,python-fastapi,nodejs" description:"Base template to use"`
		Version      string `json:"version"       example:"1.0.0"                                                              description:"Application version label"`
		Name         string `json:"name"          example:"myapp"                                                              description:"Image name"`
		Tag          string `json:"tag"           example:"latest"                                                             description:"Image tag"`
		Description  string `json:"description"   example:"Initial build for my app"`
	}
}
