package main

import (
	"log"
	"os"
	"sync"

	"github.com/dkr290/go-advanced-projects/kv-store/pkg/handlers"
	"github.com/dkr290/go-advanced-projects/kv-store/pkg/store"
	"github.com/gofiber/fiber/v2"
)

var port string

func main() {
	getEnvs()
	log.Fatal(Run())
}

func Run() error {
	store := store.NewKeyValuesStore()
	var m sync.Mutex
	h := handlers.NewHandlers(store, &m)

	app := fiber.New()
	api := app.Group("api/v1")
	api.Post("/set", h.HandlerSet)
	api.Get("/get/:database/:key", h.HandlerGet)
	api.Delete("/del", h.HandleDelete)
	api.Get("/all", h.HandlerGetAllRecords)

	return app.Listen(port)
}

func getEnvs() {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = ":3000"
	}
}
