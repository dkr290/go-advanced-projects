package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {

	return home.Index().Render(r.Context(), w)
}
