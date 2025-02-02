package handlers

import "github.com/dkr290/go-advanced-projects/cars-htmx/internal/pkg/db"

type Handler struct {
	store *db.Storage
}

func New(store *db.Storage) *Handler {
	return &Handler{
		store: store,
	}
}
