package main

import (
	"log"
	"net/http"

	"github.com/dkr290/go-advanced-projects/cars-htmx/handlers"
	"github.com/dkr290/go-advanced-projects/cars-htmx/internal/pkg/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := db.InitConfig()
	d, err := db.InitSqlLiteDb(conf)
	if err != nil {
		log.Fatal(err)
	}
	database := db.Storage{
		Db: d,
	}
	//  just to use database
	http.Handle(
		"/views/css/",
		http.StripPrefix("/views/css/", http.FileServer(http.Dir("./views/css"))),
	)
	h := handlers.New(&database)
	app := fiber.New()
	app.Static("/views/css/", "./views/css")
	app.Get("/", h.HandleHome)
	log.Println("Database type:", database.Db.Name())
	log.Fatal(app.Listen(":3000"))
}
