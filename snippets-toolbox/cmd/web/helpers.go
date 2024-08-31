package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

func (a *appconfig) render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper

	ts, ok := a.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		a.serveError(w, err)
		return
	}

	buff := new(bytes.Buffer)

	w.WriteHeader(status) //200 ok or 400 Bad Request

	// execute template
	err := ts.ExecuteTemplate(buff, "base", data)
	if err != nil {
		a.serveError(w, err)
	}

	buff.WriteTo(w)

}

// Create an newTemplateData() helper, which returns a pointer to a templateData
// struct initialized with the current year. Note that we're not using the

func (a *appconfig) newTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
	}
}
