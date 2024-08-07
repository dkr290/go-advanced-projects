package main

import (
	"fmt"
	"html/template"
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

	//initialize slice containing two files. It's importnant
	//the base template should be first one

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	//template.ParseFiles() function to read the files and sto the templates in the templateset. variadic parameter as noted in the function
	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serveError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		a.serveError(w, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
func (a *appconfig) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
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
