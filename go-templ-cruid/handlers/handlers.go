package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid/view/home"
)

func (h *Handlers) HandleHome(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
