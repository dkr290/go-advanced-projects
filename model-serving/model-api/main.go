package main

import (
	"log"
	"os"

	"github.com/dkr290/go-advanced-projects/model-api/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

const (
	modelsDir      = "models"
	llamaCPPPath   = "llama.cpp/build/bin/llama-cli" // Update this path
	maxConcurrency = 4
)

var sem = make(chan struct{}, maxConcurrency)

func main() {
	app := fiber.New()
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		panic(err)
	}
	h := handlers.NewHandlers(modelsDir, sem, llamaCPPPath)
	app.Post("api/pull", h.PullModel)
	app.Post("api/chat", h.GenerateRequest)
	app.Get("api/models", h.ListModels)
	app.Delete("api/models/:name", h.DeleteModel)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
