package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

var port = ":8080"

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

func getTeachers(c *fiber.Ctx) error {
	return c.SendString("Hello, teachers route!")
}

func postTeachers(c *fiber.Ctx) error {
	return c.SendString("Hello Post  Method on Teachers struct ")
}

func deleteTeachers(c *fiber.Ctx) error {
	return c.SendString("Hello Delete  Method on Teachers struct")
}

func patchTeachers(c *fiber.Ctx) error {
	return c.SendString("Hello Patch Method on Teachers struct")
}

func putTeachers(c *fiber.Ctx) error {
	return c.SendString("Hello Put  Method on Teachers struct")
}

func getStudents(c *fiber.Ctx) error {
	return c.SendString("Hello, students route!")
}

func patchStudents(c *fiber.Ctx) error {
	return c.SendString("Hello Patch Method on students struct")
}

func postStudents(c *fiber.Ctx) error {
	return c.SendString("Hello Post  Method on students struct ")
}

func deleteStudents(c *fiber.Ctx) error {
	return c.SendString("Hello Delete  Method on students struct")
}

func putStudents(c *fiber.Ctx) error {
	return c.SendString("Hello Put  Method on students struct")
}

func getExecs(c *fiber.Ctx) error {
	return c.SendString("Hello, execs route!")
}

func patchExecs(c *fiber.Ctx) error {
	return c.SendString("Hello Patch Method on execs struct")
}

func postExecs(c *fiber.Ctx) error {
	return c.SendString("Hello Post  Method on execs struct ")
}

func deleteExecs(c *fiber.Ctx) error {
	return c.SendString("Hello Delete  Method on execs struct")
}

func putExecs(c *fiber.Ctx) error {
	return c.SendString("Hello Put  Method on execs struct")
}

func teacherHandler(r fiber.Router) {
	r.Get("/", getTeachers)
	r.Post("/", postTeachers)
	r.Delete("/", deleteTeachers)
	r.Put("/", putTeachers)
	r.Patch("/", patchTeachers)
}

func studentHandler(r fiber.Router) {
	r.Get("/", getStudents)
	r.Post("/", postStudents)
	r.Delete("/", deleteStudents)
	r.Put("/", putStudents)
	r.Patch("/", patchStudents)
}

func execHandler(r fiber.Router) {
	r.Get("/", getExecs)
	r.Post("/", postExecs)
	r.Delete("/", deleteExecs)
	r.Put("/", putExecs)
	r.Patch("/", patchExecs)
}

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Root route!")
	})

	teachers := app.Group("/teachers")
	teacherHandler(teachers)

	students := app.Group("/students")
	studentHandler(students)

	execs := app.Group("/execs")
	execHandler(execs)
	// Start the server on port
	log.Fatal(app.Listen(port))
}
