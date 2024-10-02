package orderhandlers

import (
	"fmt"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/dkr290/go-advanced-projects/ecom/helpers"
	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type OrderHandlers struct {
	DB         db.OrderDatabaseInt
	productsDB db.ProductDatabaseInt
}

func NewOrderHandler(db db.OrderDatabaseInt, pdb db.ProductDatabaseInt) *OrderHandlers {
	return &OrderHandlers{
		DB:         db,
		productsDB: pdb,
	}
}

func (h *OrderHandlers) RegisterRoutes(router chi.Router) {
	router.Post("cart/checkout", h.handleCheckout)

}

func (h *OrderHandlers) handleCheckout(w http.ResponseWriter, r *http.Request) {

	var cart types.CardCheckoutPayload
	err := helpers.ParseJson(r, &cart)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := helpers.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return

	}
	//get products
	var items []types.CartItem
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}

	products, err := h.productsDB.GetProductByIds(productIDs)

}
