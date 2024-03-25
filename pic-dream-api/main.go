package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main(){

	if err := getEnv(); err != nil{
		log.Fatal(err)
	}

	router := chi.NewMux()

	//router.Get("/")

    port := os.Getenv("HTTP_LISTEN_ADDR")
    slog.Info("application is running","port",port)
	log.Fatal(http.ListenAndServe(os.Getenv("HTTP_LISTEN_ADDR"),router))

}


func getEnv() error {
	if err := godotenv.Load(); err != nil {
		return err 
	}
	return nil 

}