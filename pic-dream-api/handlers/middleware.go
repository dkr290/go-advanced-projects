package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
)

func (h *Handlers) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// to not go into the middleware, the public not need to be authenticated
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticatedUser(r)

		if !user.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// to not go into the middleware, the public not need to be authenticated
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := h.sb.Auth.User(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			Email:    resp.Email,
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	// do we have an user key
	// and is that key authenticated user
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}
