package user

import (
	"fmt"
	"net/http"

	"github.com/dkr290/go-advanced-projects/ecom/helpers"
	"github.com/dkr290/go-advanced-projects/ecom/types"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/login", h.handleLogin)
	router.Post("/register", h.handleRegister)
	router.Get("/test", h.testHandler)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//get json payload
	//check if the user exists

	//if it does not we create new user
	var payload types.RegisterUserPayload
	if err := helpers.ParseJson(r, payload); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
	}
	// check if the user exists

}

func (h *Handler) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Test")
}
