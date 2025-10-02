// Package router
package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/dkr290/go-advanced-projects/crd-api-deploy/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, config huma.Config) *http.ServeMux {
	handler := handlers.NewHandler()
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
		Summary:       "Create or Apply a CRD",
		Description:   "Creates or applies a new Custom Resource Definition in the Kubernetes cluster",
		Tags:          []string{"Apply-CRD"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "get-single-crd",
		Method:      http.MethodPost,
		Path:        "/crd-get",
		Summary:     "Get an application",
		Description: "Retrieves an application resource by name and namespace and crd group, kind and version",
		Tags:        []string{"GET-Crd"},
	}, handler.GetAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "list-crds",
		Method:      http.MethodPost,
		Path:        "/crd-list",
		Summary:     "List all resources",
		Description: "Lists all crd resources in the specified namespace",
		Tags:        []string{"List-Crds"},
	}, handler.ListHandler)

	return mux
}
