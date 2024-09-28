package producthandlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	router.Get("/all", h.handleGetProduct)
	router.Post("/create", h.handleCreateProduct)
	router.Put("/update/{id}", h.handleUpdateProduct)
	router.Get("/{id}", h.handleGetProductById)
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
}

func (h *ProducsHandlers) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid product id: %v", err))
		return
	}

	if err := helpers.ParseJson(r, &product); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	err = h.DB.UpdateProduct(product, id)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = helpers.WriteJSON(w, http.StatusCreated, product)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("write json error : %v", err))
		return
	}

}
func (h *ProducsHandlers) handleGetProductById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid product id: %v", err))
		return
	}
	product, err := h.DB.GetProductById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.WriteError(w, http.StatusNotFound, err)

		} else {
			helpers.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to retreive product %v", err))
		}
		return
	}
	err = helpers.WriteJSON(w, http.StatusCreated, product)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("write json error : %v", err))
		return
	}

}
