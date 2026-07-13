package main

import (
	"log"

	"fincalc/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/", handlers.Index)
	app.Get("/test", handlers.Test)

	log.Fatal(app.Listen(":3000"))
}
