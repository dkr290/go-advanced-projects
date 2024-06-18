package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const PORT = ":4000"

type application struct {
	config appConfig
}

type appConfig struct {
	useCache bool
	verbose  bool
}

func main() {
	app := application{}

	flag.BoolVar(&app.config.useCache, "cache", false, "Use template cache")
	flag.BoolVar(&app.config.verbose, "verbose", false, "Use verbose logging")
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	srv := &http.Server{
		Addr:              PORT,
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	fmt.Println("Starting the web application on port", PORT)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
