package handlers

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/settings"
)

func (h *Handlers) HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return helpers.Render(r, w, settings.Index(user))
}

func (h *Handlers) HandleSettingsUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}
	if err := helpers.ValidateUser(params.Username); err != nil {
		slog.Error("username is not valid")
		return helpers.Render(
			r,
			w,
			settings.ProfileForm(params, settings.ProfileErrors{
				Username: "Username should be min 3 chars and max 20 chars long",
			}),
		)
	}
	user := getAuthenticatedUser(r)
	user.Account.Username = params.Username

	if err := h.Bun.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	log.Println("setting updated user ")
	return helpers.Render(r, w, settings.ProfileForm(params, settings.ProfileErrors{}))
}

func (h *Handlers) HandleSettingsPasswordReset(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return h.sb.Auth.ResetPasswordForEmail(r.Context(), user.Email)
}
