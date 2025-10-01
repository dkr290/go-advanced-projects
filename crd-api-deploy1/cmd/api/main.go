package main

import (
	"fmt"
	"log"
	"net/http"

	"crd-api-deploy/cmd/router"
	"crd-api-deploy/internal/middleware"
)

type Options struct {
	Port  int    `help:"Port to listen on"      short:"p" default:"8080"`
	Debug bool   `help:"Debug flag for logging" short:"d" default:"false"     doc:"Enable Debug Logging"`
	Host  string `help:"Host to run on"         short:"s" default:"localhost" doc:"Hostname to listen to"`
}

var port = ":8080"

func main() {
	mux := router.RegisterRoutes()
	server := &http.Server{
		Addr:    port,
		Handler: middleware.ResponseTimeMiddleware(mux),
	}

	fmt.Println("The server is starting on port", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Error Starting the server", err)
	}

	// cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
	// 	server := &http.Server{
	// 		Addr:    fmt.Sprintf(":%d", options.Port),
	// 		Handler: middleware.SecurityHeaders(mux),
	// 	}
	//
	// 	hooks.OnStart(func() {
	// 		log.Printf("Starting server on port %d", options.Port)
	// 		log.Printf(
	// 			"API documentation available at http://%s:%d/docs",
	// 			options.Host,
	// 			options.Port,
	// 		)
	//
	// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 			log.Fatalf("Server failed to start: %v", err)
	// 		}
	// 	})
	//
	// 	hooks.OnStop(func() {
	// 		log.Println("Shutting down server...")
	//
	// 		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// 		defer cancel()
	//
	// 		if err := server.Shutdown(ctx); err != nil {
	// 			log.Printf("Server forced to shutdown: %v", err)
	// 		} else {
	// 			log.Println("Server exited gracefully")
	// 		}
	// 	})
	//
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// 	go func() {
	// 		<-c
	// 		log.Println("Received shutdown signal")
	// 		hooks.OnStop(func() {})
	// 	}()
	// })
	//
	// cli.Run()
}
