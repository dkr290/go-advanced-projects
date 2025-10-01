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
	huma.Get(api, "/", handler.RootHandler)
	huma.Register(api, huma.Operation{
		OperationID:   "create-simpleapi",
		Method:        http.MethodPost,
		Path:          "/crd",
		Summary:       "Create a new SimpleAPI CRD",
		Description:   "Creates a new SimpleAPI Custom Resource Definition in the Kubernetes cluster",
		Tags:          []string{"SimpleAPI"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "get-simpleapi",
		Method:      http.MethodGet,
		Path:        "/crd/{name}",
		Summary:     "Get a SimpleAPI resource",
		Description: "Retrieves a SimpleAPI resource by name and namespace",
		Tags:        []string{"SimpleAPI"},
	}, handler.GetAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "list-simpleapis",
		Method:      http.MethodGet,
		Path:        "/crd",
		Summary:     "List SimpleAPI resources",
		Description: "Lists all SimpleAPI resources in the specified namespace",
		Tags:        []string{"SimpleAPI"},
	}, handler.ListHandler)

	return mux
}
