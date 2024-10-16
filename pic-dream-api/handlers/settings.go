package handlers

import (
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/settings"
)

func (h *Handlers) HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return helpers.Render(r, w, settings.Index(user))
}
