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
		ImageName    string `json:"image_name"          example:"myapp"                                                              doc:"Image name"`
		Tag          string `json:"tag"           example:"latest"                                                             doc:"Image tag"`
		Description  string `json:"description"   example:"description" doc:"description"`
		RepoURL      string `json:"repourl" example:"https://github.com/repouser1/goproject1" doc:"The Repository URL to clone from"`
		RepoUsername string `json:"repousername"  example:"user2" doc:"The username"`
		RepoPassword string `json:"repopassword"  example:"password123" doc:"Password for the user"`
		UserAuth     bool   `json:"userauth" example:"true" doc:"true or false depends if the github needs authentication"`
	}
}
