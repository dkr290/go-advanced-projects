package producthandlers

import (
	"fmt"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/dkr290/go-advanced-projects/ecom/helpers"
	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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
	//create empty payload
	var product types.CreateProductPayload
	if err := helpers.ParseJson(r, &product); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := helpers.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return

	}
	err := h.DB.CreateProduct(product)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = helpers.WriteJSON(w, http.StatusCreated, product)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("write json error : %v", err))
		return
	}
	// TODO validation
}
