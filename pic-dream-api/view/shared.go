package view

import (
	"context"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
