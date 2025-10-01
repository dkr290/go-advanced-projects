// Package handlers
package handlers

import (
	"context"

	"crd-api-deploy/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	service *service.APIService
}

func NewHandler() *Handlers {
	return &Handlers{}
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

func (h *Handlers) CreateAPIHandler(
	ctx context.Context, input *CreateAPIInput,
) (*CreateAPIOutput, error) {
	result, err := h.service.CreateSimpleAPI(ctx, &input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create SimpleAPI", err)
	}

	return &CreateAPIOutput{
		Body: *result,
	}, nil
}

func (h *Handlers) GetAPIHandler(ctx context.Context, input *GetAPIInput) (*GetAPIOutput, error) {
	return nil, nil
}

func (h *Handlers) ListHandler(ctx context.Context, input *ListAPIInput) (*ListAPIOutput, error) {
	return nil, nil
}
