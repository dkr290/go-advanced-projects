package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/go-snippets-templ/handlers"
	"github.com/dkr290/go-advanced-projects/go-snippets-templ/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
	router := chi.NewMux()
	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	// router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	h := handlers.NewHandlers()
	router.Get("/", helpers.MakeHandler(h.HandleHomeIndex))
	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application is running", "port", port)
	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"), router))
}

func getEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
