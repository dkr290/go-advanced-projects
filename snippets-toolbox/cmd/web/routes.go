package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *appconfig) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(app.recoverPanic)
	router.Use(app.logRequest)
	router.Use(secureHeaders)
	//create a file server which serves files out of "./ui/static direct all "
	//path is relative to the project directory root
	fs := http.FileServer(http.Dir("./ui/static/"))

	//use the handler function to register the fileserver as handler

	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.Get("/", app.home)
	router.Get("/snippet/view/{id}", app.snippetView)
	router.Get("/snippet/create", app.snippetCreate)
	router.Post("/snippet/create", app.snippetCreatePost)

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	}))

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.

	return router
}
