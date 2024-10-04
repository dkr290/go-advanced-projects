package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
)

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {

	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}
