package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/dkr290/go-advanced-projects/crd-api-deploy/cmd/router"
	"github.com/dkr290/go-advanced-projects/crd-api-deploy/internal/middleware"
)

type Options struct {
	Port  int    `help:"Port to listen on"      short:"p" default:"8080"`
	Debug bool   `help:"Debug flag for logging" short:"d" default:"false"     doc:"Enable Debug Logging"`
	Host  string `help:"Host to run on"         short:"s" default:"localhost" doc:"Hostname to listen to"`
}

func main() {
	mux := http.NewServeMux()

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		config := huma.DefaultConfig("SimpleAPI CRD Deployer", "1.0.0")

		config.Info.Description = "A REST API for deploying SimpleAPI Custom Resource Definitions to Kubernetes clusters"

		config.Info.Contact = &huma.Contact{
			Name: "API Support",
		}

		config.Servers = []*huma.Server{
			{URL: fmt.Sprintf("http://localhost:%d", options.Port)},
		}
		mux := router.RegisterRoutes(mux, config)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: middleware.SecurityHeaders(middleware.ResponseTimeMiddleware(mux)),
		}

		hooks.OnStart(func() {
			log.Printf("Starting server on port %d", options.Port)
			log.Printf(
				"API documentation available at http://%s:%d/docs",
				options.Host,
				options.Port,
			)

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed to start: %v", err)
			}
		})

		hooks.OnStop(func() {
			log.Println("Shutting down server...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
			} else {
				log.Println("Server exited gracefully")
			}
		})

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Received shutdown signal")
			hooks.OnStop(func() {})
		}()
	})

	cli.Run()
}
