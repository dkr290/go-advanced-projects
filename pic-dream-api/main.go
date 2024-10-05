package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/handlers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/helpers"
	"github.com/dkr290/go-advanced-projects/pic-dream-api/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// embed public
//var FS embed.FS

func main() {

	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
	sbHost := os.Getenv("SUPABASE_URL")
	if len(sbHost) == 0 {
		log.Fatal("Neet supabase URL")
	}
	sbSecret := os.Getenv("SUPABASE_SECRET")
	if len(sbSecret) == 0 {
		log.Fatal("Need supabase token")
	}
	sbClient := db.InitDB(sbHost, sbSecret)
	//to change this when we pass the variable
	fmt.Printf("sbClient: %v\n", sbClient)

	router := chi.NewMux()
	router.Use(handlers.IsLoggedIn)

	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	//router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/", helpers.MakeHandler(handlers.HandleHomeIndex))
	router.Get("/login", helpers.MakeHandler(handlers.HandleLoginIndex))
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
