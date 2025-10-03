// Package handlers
package handlers

import (
	"context"
	"errors"
	"model-image-deployer/internal/apierror"
	"model-image-deployer/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

type Handlers struct {
	service service.APIServiceInterface
}

func NewHandler(svc service.APIServiceInterface) *Handlers {
	return &Handlers{service: svc}
}

func (h *Handlers) RootHandler(ctx context.Context, _ *struct{}) (*struct {
	Body struct {
		Message string `json:"message"`
	}
}, error,
) {
	output := &struct {
		Body struct {
			Message string `json:"message"`
		}
	}{
		Body: struct {
			Message string `json:"message"`
		}{
			Message: "Root path Healthy",
		},
	}

	return output, nil
}

func (h *Handlers) CreateCRDHandler(
	ctx context.Context, input *CreateCrdInput,
) (*CreateCrdOutput, error) {
	log.Info().Interface("request", input.Body).Msg("CreateCRD request received")
	result, err := h.service.CreateAPP(ctx, &input.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create SimpleAPI")
		switch {
		case errors.Is(err, apierror.ErrInvalidInput):
			return nil, huma.Error400BadRequest(err.Error())
		case errors.Is(err, apierror.ErrK8sConflict):
			return nil, huma.Error409Conflict(err.Error())
		default:
			return nil, huma.Error500InternalServerError("Failed to create SimpleAPI")
		}
	}

	log.Info().Interface("result", result).Msg("CreateCRD request successful")
	return &CreateCrdOutput{
		Body: *result,
	}, nil
}

func (h *Handlers) ListHandler(ctx context.Context, _ *struct{}) (*ListAPIOutput, error) {
	log.Info().Msg("ListHandler request received")
	result, err := h.service.ListAPPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list resource")
		return nil, huma.Error400BadRequest("Failed to list resource", err)
	}
	log.Info().Int("count", len(result.Items)).Msg("ListHandler request successful")
	return &ListAPIOutput{
		Body: *result,
	}, nil
}

func (h *Handlers) DeleteHandler(
	ctx context.Context,
	input *DeleteAPIInput,
) (*DeleteAPIOutput, error) {
	log.Info().Interface("request", input.Body).Msg("DeleteHandler request received")
	result, err := h.service.DeleteAPP(ctx, &input.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete resource")
		switch {
		case errors.Is(err, apierror.ErrNotFound):
			return nil, huma.Error404NotFound(err.Error())
		default:
			return nil, huma.Error500InternalServerError("Failed to delete resource")
		}
	}
	log.Info().Interface("result", result).Msg("DeleteHandler request successful")
	return &DeleteAPIOutput{
		Body: *result,
	}, nil
}
