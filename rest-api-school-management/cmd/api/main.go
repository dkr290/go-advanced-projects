package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

var port = ":8080"

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Root route!")
	})

	teachers := app.Group("/teachers", func(c *fiber.Ctx) error {
		return c.Next()
	})

	teachers.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello, teachers route!")
	})
	teachers.Patch("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Patch Method on Teachers struct")
	})
	teachers.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Post  Method on Teachers struct ")
	})
	teachers.Delete("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Delete  Method on Teachers struct")
	})

	teachers.Put("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Put  Method on Teachers struct")
	})

	students := app.Group("/students", func(c *fiber.Ctx) error {
		return c.Next()
	})
	students.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello, students route!")
	})
	students.Patch("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Patch Method on students struct")
	})
	students.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Post  Method on students struct ")
	})
	students.Delete("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Delete  Method on students struct")
	})

	students.Put("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Put  Method on students struct")
	})
	execs := app.Group("/execs", func(c *fiber.Ctx) error {
		return c.Next()
	})
	execs.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello, execs route!")
	})
	execs.Patch("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Patch Method on execs struct")
	})
	execs.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Post  Method on execs struct ")
	})
	execs.Delete("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Delete  Method on execs struct")
	})

	execs.Put("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Method())
		return c.SendString("Hello Put  Method on execs struct")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(port))
}
