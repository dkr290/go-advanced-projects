package handlers

import "model-image-deployer/internal/models"

type CreateCrdInput struct {
	Body models.CreateCrdRequest `json:"body"`
}
type CreateCrdOutput struct {
	Body models.CreateCrdResponse `json:"body"`
}

// ListAPIOutput represents the output for listing SimpleAPIs
type ListAPIOutput struct {
	Body models.ListAPIResponse `json:"body"`
}
type DeleteAPIInput struct {
	Body models.DeleteCrdInput `json:"body"`
}

type DeleteAPIOutput struct {
	Body models.DeleteCrdResponse `json:"body"`
}
