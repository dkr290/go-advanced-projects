package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/dkr2909/go-advanced-projects/design-pattern/go-app/templatecache"
)

type templateData struct {
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, t string, td *templateData) {
	var err error
	var tmpl *template.Template
	if app.config.useCache {
		tmpl, err = templatecache.Get(t, app.config.verbose, app.templatePaths(t)...)
		if err != nil {
			log.Println("Error getting the template:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
	}
	if tmpl == nil {
		newTemplate, err := template.ParseFiles(app.templatePaths(t)...)
		if err != nil {
			log.Println("Error building template:", err)
			return
		}
		if app.config.verbose {
			log.Println("Not using cache option")
		}
		tmpl = newTemplate
	}

	if td == nil {
		td = &templateData{}
	}

	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (app *application) templatePaths(t string) []string {
	return []string{
		"./templates/base.layout.html",
		"./templates/partials/header.html",
		"./templates/partials/footer.html",
		fmt.Sprintf("./templates/%s", t),
	}
}

//old render function
// type templateData struct {
// 	Data map[string]any
// }
//
// func (app *application) render(w http.ResponseWriter, t string, td *templateData) {
//
// 	var tmpl *template.Template
//
// 	// template cache  if we are using template cache, try to get the template from the map
//
// 	if app.config.useCache {
// 		if temmplateFromMap, ok := app.templateMap[t]; ok {
// 			tmpl = temmplateFromMap
// 		}
// 	}
//
// 	if tmpl == nil {
// 		newTemplate, err := app.buildTemplateFromDisk(t)
// 		if err != nil {
// 			log.Println("Error building template:", err)
// 			return
// 		}
// 		log.Println("building template from disk")
// 		tmpl = newTemplate
// 	}
//
// 	if td == nil {
// 		td = &templateData{}
// 	}
//
// 	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
// 		log.Println("Error executing template")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// func (app *application) buildTemplateFromDisk(t string) (*template.Template, error) {
//
// 	templateSlice := []string{
// 		"./templates/base.layout.html",
// 		"./templates/partials/header.html",
// 		"./templates/partials/footer.html",
// 		fmt.Sprintf("./templates/%s", t),
// 	}
//
// 	tmpl, err := template.ParseFiles(templateSlice...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	app.templateMap[t] = tmpl
//
// 	return tmpl, nil
//
// }
