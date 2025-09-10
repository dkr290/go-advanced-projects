package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

var port = ":8080"

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Root route!")
	})

	app.Get("/teachers", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello, teachers route!")
	})
	app.Get("/students", func(c *fiber.Ctx) error {
		return c.SendString("Hello, students route!")
	})
	app.Get("/execs", func(c *fiber.Ctx) error {
		return c.SendString("Hello, execs route!")
	})
	// Start the server on port 3000
	log.Fatal(app.Listen(port))
}
