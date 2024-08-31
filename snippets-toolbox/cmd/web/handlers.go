package main

import (
	"dkr290/go-advanced-projects/snippets-toolbox/internal/models"
	"errors"
	"fmt"

	"net/http"
	"strconv"
)

func (a *appconfig) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}
	// Use the template.ParseFiles() function to read the template file into
	// template set. If there's an error, log detailed error message
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.

	//initialize slice conta sxining two files. It's importnant
	//the base template should be first one

	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serveError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := a.newTemplateData(r)
	data.Snippets = snippets

	a.render(w, http.StatusOK, "home.html", data)

}
func (a *appconfig) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serveError(w, err)
		}
	}
	data := a.newTemplateData(r)
	data.Snippet = snippet
	a.render(w, http.StatusOK, "view.html", data)

}
func (a *appconfig) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	//some dummy data to be removed after

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n–Kobayashi Issa"
	expires := 7

	// pass the data to insert method from models receiving id of the new record back
	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serveError(w, err)
		return
	}
	//redirect user to relevant page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
