package main

import (
	"flag"
	"internal-pypi/internal/config"
	"log"
	"net/http"
	"os"
)

type Config struct {
	*config.AppConfig
}

func main() {

	// Initialize the AppConfig
	aconf := config.New("admin", "password", ":4000")
	app := Config{AppConfig: aconf}
	//define a new command-line flag with the name addr and default value
	flag.StringVar(&app.Username, "username", app.Username, "he username for authentication")
	flag.StringVar(&app.Password, "password", app.Password, "password for the authentication")
	flag.StringVar(&app.Port, "port", app.Port, "The port for application to run")

	flag.Parse()
	log.Println(app.Port)
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
		log.Fatalf("could not create upload directory: %v", err)
	}

	srv := http.Server{
		Addr:     app.Port,
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(),
	}

	app.InfoLog.Printf("Starting the application on port %s", app.Port)
	if err := srv.ListenAndServe(); err != nil {

		app.ErrorLog.Fatal(err)
	}
}
