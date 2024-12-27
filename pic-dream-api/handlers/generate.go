package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/view/generate"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) HandleGenereateIndex(w http.ResponseWriter, r *http.Request) error {
	// images := make([]types.Image, 20)
	data := generate.ViewData{
		Images: []types.Image{},
	}
	//	images[0].ImageStatus = types.ImageStatusPending
	return helpers.Render(r, w, generate.Index(data))
}

func (h *Handlers) HandleGenereateCreate(w http.ResponseWriter, r *http.Request) error {
	return helpers.Render(
		r,
		w,
		generate.GalerryImage(types.Image{ImageStatus: types.ImageStatusPending}),
	)
}

func (h *Handlers) HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	// fetch from db
	image := types.Image{
		ImageStatus: types.ImageStatusPending,
	}
	slog.Info("checkiung image status", "id", id)
	return helpers.Render(r, w, generate.GalerryImage(image))
}
