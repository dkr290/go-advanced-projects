package handlers

import (
	"html/template"
	"log"

	"github.com/gofiber/fiber/v3"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
}

func Index(c fiber.Ctx) error {
	data := map[string]string{
		"Title":   "Fincalc",
		"Message": "Welcome to Fincalc",
	}
	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "index.html", data)
}

func Test(c fiber.Ctx) error {
	data := map[string]string{
		"Title": "Test Page",
	}
	c.Set("Content-Type", "text/html")
	return templates.ExecuteTemplate(c.Response().BodyWriter(), "test.html", data)
}
