package types

type contextKey string

const UserContextKey contextKey = "user"

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}
