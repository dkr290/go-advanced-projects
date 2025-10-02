package handlers

import "github.com/dkr290/go-advanced-projects/crd-api-deploy/internal/models"

type CreateAPIInput struct {
	Body models.CreateAPIRequest `json:"body"`
}
type CreateAPIOutput struct {
	Body models.CreateAPIResponse `json:"body"`
}

type GetAPIInput struct {
	Body models.GetSigleCrdInput `json:"body"`
}

type GetAPIOutput struct {
	Body models.GetAPIResponse `json:"body"`
}

// ListAPIInput represents the input for listing SimpleAPIs
type ListAPIInput struct {
	Body models.ListAPIInput `json:"body"`
}

// ListAPIOutput represents the output for listing SimpleAPIs
type ListAPIOutput struct {
	Body models.ListAPIResponse `json:"body"`
}
