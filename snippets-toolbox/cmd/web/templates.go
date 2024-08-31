package main

import (
	"dkr290/go-advanced-projects/snippets-toolbox/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

//define templateData type to act as holding ctructure for any dynamic data
//to pass to the html templates

type TemplateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

//function to return nicely format date can be only return one value and potentially error
//it can accept multiply parameters

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

//initialize a template.FuncMap object and store in a global variable
//this is keymap which acts as lookup between the names of our template functions
// and functions themselves

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	//new map for the template chache
	cache := make(map[string]*template.Template)

	//filepath.Glob is used to get slice of all filepaths that match pattern ./ui/html/pages/*.html

	pages, err := filepath.Glob("./ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	//loop through the pages filepaths one by one

	for _, page := range pages {
		//extract the filename like view.html

		name := filepath.Base(page)
		//create a slice containing filepaths for our base template
		// The template.FuncMap must be registered with the template set
		//before call the ParseFiles() method. This means we have to use
		// template.New() to create an empty template set, use the Funcs() method to register
		// template.FuncMap, and then parse the file as normal.

		tmpl := template.New(name).Funcs(functions)

		ts, err := tmpl.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}
	return cache, nil
}
