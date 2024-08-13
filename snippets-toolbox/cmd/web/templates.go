package main

import (
	"dkr290/go-advanced-projects/snippets-toolbox/internal/models"
	"html/template"
	"path/filepath"
)

//define templateData type to act as holding ctructure for any dynamic data
//to pass to the html templates

type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}
	return cache, nil
}
