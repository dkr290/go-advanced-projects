package orderhandlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/db"
	"github.com/go-chi/chi/v5"
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

}
