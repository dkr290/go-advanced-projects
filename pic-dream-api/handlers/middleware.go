package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
)

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//to not go into the middleware, the public not need to be authenticated
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			Email:    "dani@abv.bg",
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
