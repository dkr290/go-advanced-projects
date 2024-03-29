package handlers

import (
	"net/http"

	"github.com/dkr290/go-events-booking-api/pic-dream-api/view/userauth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {

	return userauth.LogIn().Render(r.Context(), w)
}
