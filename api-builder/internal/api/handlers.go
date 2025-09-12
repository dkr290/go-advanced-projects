package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/models"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/services"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	dockerService *services.DockerService
	clog          *logrus.Logger
}

type GetBuildStatus struct {
	BuildID string `path:"buildId" doc:"The Build ID to query"`
}
type BuildImageOutput struct {
	Body models.BuildImageResponse `json:"body"`
}
type GetBuildStatusOutput struct {
	Body models.BuildStatus `json:"body"`
}

type BuildImageInput struct {
	Body models.BuildImageRequest
}

func NewHandlers(dockerService *services.DockerService, clog *logrus.Logger) *Handlers {
	return &Handlers{
		dockerService: dockerService,
		clog:          clog,
	}
}

// BuildImage handles POST /api/v1/build
func (h *Handlers) BuildImage(
	ctx context.Context,
	input *BuildImageInput,
) (*BuildImageOutput, error) {
	req := input.Body
	if req.ModelVersion == "" || req.Version == "" || req.Name == "" || req.Tag == "" {
		return nil, huma.Error400BadRequest("model_version, version, name, and tag are required")
	}
	h.clog.Infof(
		"Building image: %s:%s with model version: %s",
		req.Name,
		req.Tag,
		req.ModelVersion,
	)

	resp, err := h.dockerService.BuildImage(ctx, &models.BuildImageRequest{
		ModelVersion: req.ModelVersion,
		Version:      req.Version,
		Name:         req.Name,
		Tag:          req.Tag,
		Description:  req.Description,
	})
	if err != nil {
		h.clog.Errorf("Failed to build image: %v", err)
		return nil, huma.Error500InternalServerError("Failed to initiate build", err)

	}

	return &BuildImageOutput{Body: *resp}, nil
}

// GetBuildStatus handles GET /api/v1/build/:buildId/status
func (h *Handlers) GetBuildStatus(
	ctx context.Context,
	input *GetBuildStatus,
) (*GetBuildStatusOutput, error) {
	if input.BuildID == "" {
		return nil, huma.Error400BadRequest("buildId is required")
	}
	status, err := h.dockerService.GetBuildStatus(input.BuildID)
	if err != nil {
		return nil, huma.Error404NotFound("Build not found", err)
	}

	return &GetBuildStatusOutput{Body: *status}, nil
}
