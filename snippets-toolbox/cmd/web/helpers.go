package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (a *appconfig) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errotLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *appconfig) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *appconfig) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}
