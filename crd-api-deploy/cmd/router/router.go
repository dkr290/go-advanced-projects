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
		Path:          "/crd-create",
		Summary:       "Create or Apply a CRD",
		Description:   "Creates or applies a new Custom Resource Definition in the Kubernetes cluster",
		Tags:          []string{"Apply-CRD"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "get-simpleapi",
		Method:      http.MethodPost,
		Path:        "/crd-get",
		Summary:     "Get an application",
		Description: "Retrieves an application resource by name and namespace and crd group, kind and version",
		Tags:        []string{"GET-Deployment"},
	}, handler.GetAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "list-simpleapis",
		Method:      http.MethodGet,
		Path:        "/crd",
		Summary:     "List all resources",
		Description: "Lists all deployment resources in the specified namespace",
		Tags:        []string{"List-Deployments"},
	}, handler.ListHandler)

	return mux
}
