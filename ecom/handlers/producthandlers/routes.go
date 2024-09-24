package producthandlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/dkr290/go-advanced-projects/ecom/helpers"
	"github.com/go-chi/chi/v5"
)

type ProducsHandlers struct {
	DB db.ProductDatabaseInt
}

func NewProductHandler(db db.ProductDatabaseInt) *ProducsHandlers {
	return &ProducsHandlers{
		DB: db,
	}
}

func (h *ProducsHandlers) RegisterRoutes(router chi.Router) {
	router.Get("/products", h.handleGetProduct)
	router.Post("/products", h.handleCreateProduct)
}

func (h *ProducsHandlers) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	prods, err := h.DB.GetProducts()
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.CustomErrorMessage("error handle create product", err))
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, prods)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.CustomErrorMessage("error write Json ", err))
	}
}

func (h *ProducsHandlers) handleCreateProduct(w http.ResponseWriter, r *http.Request) {

}
