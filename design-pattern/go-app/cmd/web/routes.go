package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	//if somthing crashes to recover from panic
	//middleware sits betwewen router and handler

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Get("/", app.ShowHome)
	return mux

}
