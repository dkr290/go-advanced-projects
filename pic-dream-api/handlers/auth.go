package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/userauth"
	"github.com/nedpals/supabase-go"
)

func (s *Handlers) HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(r, w, userauth.LogIn())
}

func (s *Handlers) HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(r, w, userauth.SignUp())
}

func (s *Handlers) HandleSignupCretate(w http.ResponseWriter, r *http.Request) error {
	params := userauth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	if !helpers.ValidateEmail(params.Email) {
		slog.Error("email is not valid")
		return helpers.Render(r, w, userauth.SignUpForm(params, userauth.SignupErrors{
			Email: "The email is invalid",
		}))

	}
	if err := helpers.ValidatePassword(params.Password); err != nil {
		slog.Error("password is not valid")
		return helpers.Render(r, w, userauth.SignUpForm(params, userauth.SignupErrors{
			Password: "The password is invalid: " + err.Error(),
		}))

	}
	if params.Password != params.ConfirmPassword {
		slog.Error("The passwords are not the same ")
		return helpers.Render(r, w, userauth.SignUpForm(params, userauth.SignupErrors{
			ConfirmPassword: "Password mismatch",
		}))
	}

	return nil
}

func (s *Handlers) HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	// this should be valid only if we create user and password
	// if err := helpers.ValidatePassword(credentials.Password); err != nil {
	// 	slog.Error("password is not valid")
	// 	return helpers.Render(r, w, userauth.LoginForm(credentials, userauth.LoginErrors{
	// 		Password: "The password is invalid: " + err.Error(),
	// 	}))
	//
	// }
	//
	// calling the supabase
	resp, err := s.sb.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("user login error", slog.String("error", err.Error()))
		return helpers.Render(r, w, userauth.LoginForm(credentials, userauth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))

	}
	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
