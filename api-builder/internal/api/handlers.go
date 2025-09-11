package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/models"
	"github.com/dkr290/go-advanced-projects/api-builder/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	dockerService *services.DockerService
}

func NewHandlers(dockerService *services.DockerService) *Handlers {
	return &Handlers{
		dockerService: dockerService,
	}
}

// BuildImage handles POST /api/v1/build
func (h *Handlers) BuildImage(
	ctx context.Context,
	input *models.BuildRequest,
) (*models.BuildResponse, error) {
	if input.ModelVersion == "" || input.Version == "" || input.Name == "" || input.Tag == "" {
		return nil, huma.Error400BadRequest("model_version, version, name, and tag are required")
	}
	logrus.Infof(
		"Building image: %s:%s with model version: %s",
		input.Name,
		input.Tag,
		input.ModelVersion,
	)

	resp, err := h.dockerService.BuildImage(ctx, &models.BuildRequest{
		ModelVersion: input.ModelVersion,
		Version:      input.Version,
		Name:         input.Name,
		Tag:          input.Tag,
	})
	if err != nil {
		logrus.Errorf("Failed to build image: %v", err)
		return nil, huma.Error500InternalServerError("Failed to initiate build", err)

	}

	return &models.BuildResponse{
		BuildID: resp.BuildID,
		Status:  resp.Status,
	}, nil
}

// GetBuildStatus handles GET /api/v1/build/:buildId/status
func (h *Handlers) GetBuildStatus(c *fiber.Ctx) error {
	buildID := c.Params("buildId")
	if buildID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Build ID is required",
			Code:  fiber.StatusBadRequest,
		})
	}

	status, err := h.dockerService.GetBuildStatus(buildID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "Build not found",
			Code:    fiber.StatusNotFound,
			Details: err.Error(),
		})
	}

	return c.JSON(status)
}

// ListBuilds handles GET /api/v1/builds
func (h *Handlers) ListBuilds(c *fiber.Ctx) error {
	builds := h.dockerService.ListBuilds()
	return c.JSON(fiber.Map{
		"builds": builds,
		"total":  len(builds),
	})
}
