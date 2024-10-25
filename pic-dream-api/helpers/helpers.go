package helpers

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
)

func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
		}
	}
}

func Render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func ValidateEmail(email string) bool {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return pattern.MatchString(email)
}

func ValidatePassword(password string) error {
	if len(password) < 4 {
		return errors.New("password must be at least 4 characters long")
	}
	if len(password) > 20 {
		return errors.New("password must be no more than 20 characters long")
	}
	if matched, _ := regexp.MatchString("[A-Z]", password); !matched {
		return errors.New("password must contain at least one uppercase letter")
	}
	if matched, _ := regexp.MatchString("[a-z]", password); !matched {
		return errors.New("password must contain at least one lowercase letter")
	}
	if matched, _ := regexp.MatchString("[0-9]", password); !matched {
		return errors.New("password must contain at least one digit")
	}
	if matched, _ := regexp.MatchString("[!@#$%^&*(),.?\":{}|<>]", password); !matched {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// set header and make cliebnt side redirect
func HxRedirect(w http.ResponseWriter, r *http.Request, to string) error {
	if len(r.Header.Get("HX-Request")) > 0 {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}

func ValidateUser(username string) error {
	if len(username) < 2 {
		return errors.New("username must be at least 2 characters long")
	}
	if len(username) > 20 {
		return errors.New("username must be no more than 20 characters long")
	}

	return nil
}
