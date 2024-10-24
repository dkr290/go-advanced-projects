package main

import "net/http"

func (app *appconfig) routes() http.Handler {
	mux := http.NewServeMux()

	//create a file server which serves files out of "./ui/static direct all "
	//path is relative to the project directory root
	fs := http.FileServer(http.Dir("./ui/static/"))

	//use the handler function to register the fileserver as handler

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
