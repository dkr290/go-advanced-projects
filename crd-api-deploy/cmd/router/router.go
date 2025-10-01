package router

import (
	"net/http"

	"crd-api-deploy/internal/handlers"

	"github.com/danielgtaylor/huma/v2"
)

// // SimpleAPIHandler handles SimpleAPI operations
// type SimpleAPIHandler struct {
// 	service *service.SimpleAPIService
// }
//
// // NewSimpleAPIHandler creates a new handler
// func NewSimpleAPIHandler() (*SimpleAPIHandler, error) {
// 	svc, err := service.NewSimpleAPIService()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &SimpleAPIHandler{
// 		service: svc,
// 	}, nil
// }
//

// // CreateSimpleAPI creates a new SimpleAPI CRD
// func (h *SimpleAPIHandler) CreateSimpleAPI(
// 	ctx context.Context,
// 	input *CreateSimpleAPIInput,
// ) (*CreateSimpleAPIOutput, error) {
// 	result, err := h.service.CreateSimpleAPI(ctx, &input.Body)
// 	if err != nil {
// 		return nil, huma.Error400BadRequest("Failed to create SimpleAPI", err)
// 	}
//
// 	return &CreateSimpleAPIOutput{
// 		Body: *result,
// 	}, nil
// }

// // GetSimpleAPI retrieves a SimpleAPI resource
// func (h *SimpleAPIHandler) GetSimpleAPI(
// 	ctx context.Context,
// 	input *GetSimpleAPIInput,
// ) (*GetSimpleAPIOutput, error) {
// 	result, err := h.service.GetSimpleAPI(ctx, input.Name, input.Namespace)
// 	if err != nil {
// 		return nil, huma.Error404NotFound("SimpleAPI not found", err)
// 	}
//
// 	return &GetSimpleAPIOutput{
// 		Body: *result,
// 	}, nil
// }

// // ListSimpleAPIs lists SimpleAPI resources
// func (h *SimpleAPIHandler) ListSimpleAPIs(
// 	ctx context.Context,
// 	input *ListSimpleAPIInput,
// ) (*ListSimpleAPIOutput, error) {
// 	result, err := h.service.ListSimpleAPIs(ctx, input.Namespace)
// 	if err != nil {
// 		return nil, huma.Error500InternalServerError("Failed to list SimpleAPIs", err)
// 	}
//
// 	return &ListSimpleAPIOutput{
// 		Body: *result,
// 	}, nil
// }

// RegisterRoutes registers all the routes for the SimpleAPI handler
func RegisterRoutes(api huma.API) error {
	handler := handlers.NewHandler()

	huma.Register(api, huma.Operation{
		OperationID:   "create-simpleapi",
		Method:        http.MethodPost,
		Path:          "/api/v1/simpleapis",
		Summary:       "Create a new SimpleAPI CRD",
		Description:   "Creates a new SimpleAPI Custom Resource Definition in the Kubernetes cluster",
		Tags:          []string{"SimpleAPI"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "get-simpleapi",
		Method:      http.MethodGet,
		Path:        "/api/v1/simpleapis/{name}",
		Summary:     "Get a SimpleAPI resource",
		Description: "Retrieves a SimpleAPI resource by name and namespace",
		Tags:        []string{"SimpleAPI"},
	}, handler.GetAPIHandler)

	huma.Register(api, huma.Operation{
		OperationID: "list-simpleapis",
		Method:      http.MethodGet,
		Path:        "/api/v1/simpleapis",
		Summary:     "List SimpleAPI resources",
		Description: "Lists all SimpleAPI resources in the specified namespace",
		Tags:        []string{"SimpleAPI"},
	}, handler.ListHandler)

	return nil
}
