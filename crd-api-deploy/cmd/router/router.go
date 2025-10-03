// Package router
package router

import (
	"net/http"

	"model-image-deployer/internal/handlers"
	"model-image-deployer/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func RegisterRoutes(mux *http.ServeMux, config huma.Config, apiService service.APIServiceInterface) *http.ServeMux {
	handler := handlers.NewHandler(apiService)
	api := humago.New(mux, config)
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Root path",
		Description: "Root Path",
		Tags:        []string{"Root-check"},
	}, handler.RootHandler)

	huma.Register(api, huma.Operation{
		OperationID:   "create-crd",
		Method:        http.MethodPost,
		Path:          "/crd-create",
		Summary:       "Create or Apply model CRD",
		Description:   "Creates or applies a new Custom Resource Definition for a model in the Kubernetes cluster",
		Tags:          []string{"Apply-ModelCrd"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateCRDHandler)

	huma.Register(api, huma.Operation{
		OperationID: "list-crds",
		Method:      http.MethodGet,
		Path:        "/crd-list",
		Summary:     "List all models CRD",
		Description: "Lists all crd resources for the models",
		Tags:        []string{"List-ModelCrds"},
	}, handler.ListHandler)
	huma.Register(api, huma.Operation{
		OperationID: "delete-crds",
		Method:      http.MethodPost,
		Path:        "/crd-delete",
		Summary:     "Delete model CRD",
		Description: "Delete crd resources for the models",
		Tags:        []string{"Delete-ModelCrd"},
	}, handler.DeleteHandler)

	return mux
}
