package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/handlers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// embed public
var FS embed.FS

func main() {

	if err := getEnv(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handlers.WithUser)

	//router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/", helpers.MakeHandler(handlers.HandleHomeIndex))
	router.Get("/login", helpers.MakeHandler(handlers.HandleLogin))
	router.Post("/login", helpers.MakeHandler(handlers.HandleLoginCreate))

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
