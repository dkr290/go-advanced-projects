package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dkr290/go-events-booking-api/pic-dream-api/models"
)

func getAuthenticatedUser(r *http.Request) models.AuthenticatedUser {

	user, ok := r.Context().Value(userKey).(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}
	}

	return user
}

func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
		}
	}
}
