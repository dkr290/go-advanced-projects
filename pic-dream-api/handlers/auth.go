package handlers

import (
	"fmt"
	"net/http"

	"github.com/dkr290/go-events-booking-api/pic-dream-api/view/userauth"
	"github.com/nedpals/supabase-go"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {

	return userauth.LogIn().Render(r.Context(), w)
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {

	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	fmt.Println(credentials)
	return nil
}
