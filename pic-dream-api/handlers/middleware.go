package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/dkr290/go-events-booking-api/pic-dream-api/models"
)

const userKey = "user"

func WithUser(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := models.AuthenticatedUser{}
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
