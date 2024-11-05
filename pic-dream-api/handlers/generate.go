package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/generate"
)

func (h *Handlers) HandleGenereateIndex(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(r, w, generate.Index())
}
