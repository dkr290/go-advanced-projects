package main

import (
	"log"

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
	h := handlers.New(&database)
	app := fiber.New()
	app.Static("/public/", "./public/")
	app.Get("/", h.HandleHome)
	app.Get("/cars", h.HandleListCars)
	app.Post("/cars", h.HandleAddCar)
	log.Println("Database type:", database.Db.Name())
	log.Fatal(app.Listen(":3000"))
}
