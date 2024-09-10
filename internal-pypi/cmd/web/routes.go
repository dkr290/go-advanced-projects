package main

import (
	"net/http"
)

func (app *Config) routes() *http.ServeMux {
	mux := http.NewServeMux()
	//create a file server which serves files out of "./ui/static direct all "
	//path is relative to the project directory root
	fs := http.FileServer(http.Dir("./static/"))

	//use the handler function to register the fileserver as handler
	mux.Handle("/static", http.StripPrefix("/static", fs))
	mux.HandleFunc("/favicon.ico", app.favIconHandler)
	mux.HandleFunc("/", basicAuth(app.indexHandler, app.Username, app.Password))
	mux.HandleFunc("/simple/", basicAuth(app.simpleHandler, app.Username, app.Password))
	mux.HandleFunc("/packages/", basicAuth(app.packageHandler, app.Username, app.Password))
	mux.HandleFunc("/upload", basicAuth(app.uploadHandler, app.Username, app.Password))
	mux.HandleFunc("/about", basicAuth(app.aboutHandler, app.Username, app.Password))
	mux.HandleFunc("/contact", basicAuth(app.contactHandler, app.Username, app.Password))
	mux.HandleFunc("/test", app.testHandler)

	return mux
}
