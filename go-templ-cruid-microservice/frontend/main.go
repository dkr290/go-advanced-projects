package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/frontend/handlers"
	"github.com/dkr290/go-advanced-projects/go-templ-cruid-microservice/frontend/helpers"
	"github.com/go-chi/chi"
)

var (
	backendService string
	portAddress    string
)

func main() {
	getEnvs()
	log.Fatal(Run())
}

func Run() error {
	h := handlers.NewHandlers(backendService)

	r := chi.NewRouter()
	r.Get("/", helpers.MakeHandler(h.HandleHome))
	// get all tasks
	r.Get("/tasks", helpers.MakeHandler(h.HandleFetchTasks))
	// Fetch add task form
	r.Get("/getnewtaskform", helpers.MakeHandler(h.HandleGetTaskForm))
	// add task post handler
	r.Post("/tasks", helpers.MakeHandler(h.HandleAddTask))

	// fetch update form for update by ID
	r.Get("/gettaskupdateform/{id}", helpers.MakeHandler(h.HandleGetTaskUpdateForm))

	// update the task
	r.Put("/task/{id}", helpers.MakeHandler(h.HandleUpdateTask))
	r.Post("/task/{id}", helpers.MakeHandler(h.HandleUpdateTask))
	r.Get("/favicon.ico", h.FavIconHandler)
	r.Get("/test", h.TestHandler)

	// delete the task
	r.Delete("/task/{id}", helpers.MakeHandler(h.HandleDeleteTask))

	slog.Info("application is running", "port", portAddress)
	return http.ListenAndServe(portAddress, r)
}

func getEnvs() {
	backendService = os.Getenv("BACKEND_SERVICE")
	if len(backendService) == 0 {
		backendService = "backend:3000"
	}
	portAddress = os.Getenv("HTTP_LISTEN_ADDR")
	if len(portAddress) == 0 {
		portAddress = ":8090"
	}
}
