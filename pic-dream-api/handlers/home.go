package handlers

import (
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/home"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {

	//return fmt.Errorf("failed to generate picture")

	return home.Index().Render(r.Context(), w)

}
