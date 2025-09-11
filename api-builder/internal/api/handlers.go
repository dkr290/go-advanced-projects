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
}

type GetBuildStatus struct {
	BuildID string `path:"buildId" doc:"The Build ID to query"`
}
type ListBuildResponse struct {
	Builds []*models.BuildStatus `json:"builds" doc:"All known Builds"`
	Total  int                   `json:"total"  doc:"Total number of Builds"`
}
type BuildImageOutput struct {
	Body models.BuildResponse `json:"body"`
}
type GetBuildStatusOutput struct {
	Body models.BuildStatus `json:"body"`
}
type BuildImageRequest struct {
	ModelVersion string `json:"model_version" example:"python-flask"             enum:"python-flask,python-fastapi,nodejs" description:"Base template to use"`
	Version      string `json:"version"       example:"1.0.0"                                                              description:"Application version label"`
	Name         string `json:"name"          example:"myapp"                                                              description:"Image name"`
	Tag          string `json:"tag"           example:"latest"                                                             description:"Image tag"`
	Description  string `json:"description"   example:"Initial build for my app"`
}
type BuildImageInput struct {
	Body BuildImageRequest
}

func NewHandlers(dockerService *services.DockerService) *Handlers {
	return &Handlers{
		dockerService: dockerService,
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
	logrus.Infof(
		"Building image: %s:%s with model version: %s",
		req.Name,
		req.Tag,
		req.ModelVersion,
	)

	resp, err := h.dockerService.BuildImage(ctx, &models.BuildRequest{
		ModelVersion: req.ModelVersion,
		Version:      req.Version,
		Name:         req.Name,
		Tag:          req.Tag,
		Description:  req.Description,
	})
	if err != nil {
		logrus.Errorf("Failed to build image: %v", err)
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

// ListBuilds handles GET /api/v1/builds
func (h *Handlers) ListBuilds(ctx context.Context, _ *struct{}) (*ListBuildResponse, error) {
	builds := h.dockerService.ListBuilds()
	return &ListBuildResponse{
		Builds: builds,
		Total:  len(builds),
	}, nil
}
