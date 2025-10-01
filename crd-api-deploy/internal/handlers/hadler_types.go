package handlers

import "crd-api-deploy/internal/models"

type CreateAPIInput struct {
	Body models.CreateAPIRequest `json:"body"`
}
type CreateAPIOutput struct {
	Body models.CreateAPIResponse `json:"body"`
}

type GetAPIInput struct {
	Name      string `path:"name" doc:"Name of the SimpleAPI resource"`
	Namespace string `            doc:"Namespace of the SimpleAPI resource" query:"namespace" default:"default"`
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
