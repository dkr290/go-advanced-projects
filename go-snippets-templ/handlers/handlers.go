package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/go-snippets-templ/view/home"
)

func (h *Handlers) HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
