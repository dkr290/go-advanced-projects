package handlers

import "github.com/dkr290/go-advanced-projects/crd-api-deploy/internal/models"

type CreateAPIInput struct {
	Body models.CreateAPIRequest `json:"body"`
}
type CreateAPIOutput struct {
	Body models.CreateAPIResponse `json:"body"`
}

type GetAPIInput struct {
	Body models.GetAPIInput `json:"body"`
}

type GetAPIOutput struct {
	Body models.GetAPIResponse `json:"body"`
}

// ListAPIInput represents the input for listing SimpleAPIs
type ListAPIInput struct {
	Namespace string `query:"namespace" default:"default" doc:"Namespace to list SimpleAPI resources from"`
}

// ListAPIOutput represents the output for listing SimpleAPIs
type ListAPIOutput struct {
	Body models.ListSimpleAPIResponse `json:"body"`
}
