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

	"crd-api-deploy/cmd/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8090"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		mux := http.NewServeMux()

		config := huma.DefaultConfig("SimpleAPI CRD Deployer", "1.0.0")
		config.Info.Description = "A REST API for deploying SimpleAPI Custom Resource Definitions to Kubernetes clusters"
		config.Info.Contact = &huma.Contact{
			Name: "API Support",
		}
		config.Servers = []*huma.Server{
			{URL: fmt.Sprintf("http://localhost:%d", options.Port)},
		}

		api := humago.New(mux, config)

		err := router.RegisterRoutes(api)
		if err != nil {
			log.Fatalf("Failed to register routes: %v", err)
		}

		api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
			ctx.SetHeader("Access-Control-Allow-Origin", "*")
			ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if ctx.Method() == "OPTIONS" {
				ctx.SetStatus(200)
				return
			}

			next(ctx)
		})

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: mux,
		}

		hooks.OnStart(func() {
			log.Printf("Starting server on port %d", options.Port)
			log.Printf("API documentation available at http://localhost:%d/docs", options.Port)

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

