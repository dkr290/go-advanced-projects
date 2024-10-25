package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/userauth"
	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "accessToken"
)

func (s *Handlers) HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(r, w, userauth.AccountSetup())
}

func (s *Handlers) HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	params := userauth.AccountSetupFormParams{
		Username: r.FormValue("username"),
	}
	fmt.Println("test")
	if err := helpers.ValidateUser(params.Username); err != nil {
		slog.Error("username is not valid")
		return helpers.Render(
			r,
			w,
			userauth.AccountSetupForm(params, userauth.AccountSetupFormErrors{
				Username: "Username is invalid",
			}),
		)
	}

	return helpers.HxRedirect(w, r, "/")
}

func (s *Handlers) HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(r, w, userauth.LogIn())
}

func (s *Handlers) HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(r, w, userauth.SignUp())
}

func (s *Handlers) HandleLoginGithub(w http.ResponseWriter, r *http.Request) error {
	resp, err := s.sb.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "github",
		RedirectTo: s.github_redirect_url,
	})
	if err != nil {
		return err
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
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
			Email: "the email is invalid",
		}))

	}
	if err := helpers.ValidatePassword(params.Password); err != nil {
		slog.Error("password is not valid")
		return helpers.Render(r, w, userauth.SignUpForm(params, userauth.SignupErrors{
			Password: "the password is invalid: " + err.Error(),
		}))

	}
	if params.Password != params.ConfirmPassword {
		slog.Error("The passwords are not the same ")
		return helpers.Render(r, w, userauth.SignUpForm(params, userauth.SignupErrors{
			ConfirmPassword: "password do not match",
		}))
	}

	sbUser, err := s.sb.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return helpers.Render(r, w, userauth.SignupSuccess(sbUser.Email))
}

func (s *Handlers) HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = ""
	if err := session.Save(r, w); err != nil {
		return err
	}

	// Redirect the user to the home page or login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (s *Handlers) HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	// calling the supabase
	resp, err := s.sb.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("user login error", slog.String("error", err.Error()))
		return helpers.Render(r, w, userauth.LoginForm(credentials, userauth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))

	}

	if err := setAuthSession(w, r, resp.AccessToken); err != nil {
		return err
	}
	return helpers.HxRedirect(w, r, "/")
}

func (s *Handlers) HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return helpers.Render(r, w, userauth.CallbackScript())
	}
	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = accessToken
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}
