package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/dkr290/go-events-booking-api/pic-dream-api/models"
)

func WithUser(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := models.AuthenticatedUser{
			Email:    "dani@gmail.com",
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), models.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
