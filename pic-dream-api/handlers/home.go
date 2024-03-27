package handlers

import (
	"net/http"

	"github.com/dkr290/go-events-booking-api/pic-dream-api/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	//return fmt.Errorf("failed to generate picture")
	return home.Index().Render(r.Context(), w)
}
