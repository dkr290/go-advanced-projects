package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dkr2909/go-advanced-projects/design-pattern/go-app/pets"
	"github.com/dkr2909/go-advanced-projects/design-pattern/go-app/toolbox"
	"github.com/go-chi/chi/v5"
)

func (app *application) ShowHome(w http.ResponseWriter, r *http.Request) {
	app.render(w, "home.html", nil)
}

func (app *application) ShowPage(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")
	app.render(w, fmt.Sprintf("%s.html", page), nil)
}

func (app *application) CreateDogFromFactory(w http.ResponseWriter, r *http.Request) {
	dog := pets.NewPet("dog")

	data, err := json.Marshal(dog)
	if err != nil {
		log.Fatal(err)
	}

	// Set the content type and send response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)

}
func (app *application) CreateCatFromFactory(w http.ResponseWriter, r *http.Request) {
	dog := pets.NewPet("cat")

	data, err := json.Marshal(dog)
	if err != nil {
		log.Fatal(err)
	}

	// Set the content type and send response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)

}

func (app *application) TestPatterns(w http.ResponseWriter, r *http.Request) {
	app.render(w, "test.html", nil)
}

func (app *application) CreateDogFromAbstractFactory(w http.ResponseWriter, r *http.Request) {
	var t toolbox.Tools

	dog, err := pets.NewPetFromAbstractFactory("dog")
	if err != nil {
		_ = t.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(dog)
	if err != nil {
		log.Fatal(err)
	}
	// Set the content type and send response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)

}

func (app *application) CreateCatFromAbstractFactory(w http.ResponseWriter, r *http.Request) {
	var t toolbox.Tools

	cat, err := pets.NewPetFromAbstractFactory("cat")
	if err != nil {
		_ = t.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(cat)
	if err != nil {
		log.Fatal(err)
	}
	// Set the content type and send response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)

}
