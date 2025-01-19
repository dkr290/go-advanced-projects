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
	st := store.NewKeyValuesStore()
	var m sync.Mutex
	h := handlers.NewHandlers(st, &m)

	app := fiber.New()
	api := app.Group("api/v1")
	api.Post("/set", h.HandlerSet)
	api.Get("/get/:database/:key", h.HandlerGet)
	api.Delete("/del/:database/:key", h.HandleDelete)
	api.Get("/all/:database", h.HandlerGetAllRecords)

	v2store := store.NewV2KeyValuesStore()
	h2 := handlers.NewV2Handlers(v2store)
	api1 := app.Group("api/v2")
	api1.Post("/set", h2.V2HandlerSet)
	api1.Get("/get/:database/:key", h2.V2HandlerGet)
	api1.Delete("/del/:database/:key", h2.V2HandleDelete)
	api1.Get("/all/:database", h2.V2HandlerGetAllRecords)

	return app.Listen(port)
}

func getEnvs() {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = ":3000"
	}
}
