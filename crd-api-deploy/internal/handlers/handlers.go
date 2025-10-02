// Package handlers
package handlers

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/crd-api-deploy/internal/service"
)

type Handlers struct {
	service *service.APIService
}

func NewHandler() *Handlers {
	svc, err := service.NewAPIService()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize APIService: %v", err))
	}
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

func (h *Handlers) CreateAPIHandler(
	ctx context.Context, input *CreateAPIInput,
) (*CreateAPIOutput, error) {
	result, err := h.service.CreateAPP(ctx, &input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create SimpleAPI", err)
	}

	return &CreateAPIOutput{
		Body: *result,
	}, nil
}

func (h *Handlers) GetAPIHandler(ctx context.Context, input *GetAPIInput) (*GetAPIOutput, error) {
	result, err := h.service.GetAPPResouce(ctx, &input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create SimpleAPI", err)
	}
	return &GetAPIOutput{
		Body: *result,
	}, nil
}

func (h *Handlers) ListHandler(ctx context.Context, input *ListAPIInput) (*ListAPIOutput, error) {
	result, err := h.service.ListAPPs(ctx, &input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create SimpleAPI", err)
	}
	return &ListAPIOutput{
		Body: *result,
	}, nil
}
