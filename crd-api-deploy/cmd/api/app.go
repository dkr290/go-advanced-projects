// Package api
package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"model-image-deployer/cmd/router"
	"model-image-deployer/config"
	"model-image-deployer/internal/k8s"
	"model-image-deployer/internal/logger"
	"model-image-deployer/internal/middleware"
	"model-image-deployer/internal/service"
	"model-image-deployer/internal/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/rs/zerolog/log"
)

// Options struct for command-line flags.
type Options struct {
	Port  int    `help:"Port to listen on"      short:"p" default:"8080"`
	Debug bool   `help:"Debug flag for logging" short:"d" default:"false"     doc:"Enable Debug Logging"`
	Host  string `help:"Host to run on"         short:"s" default:"localhost" doc:"Hostname to listen to"`
}

// App struct to hold the CLI application.
type App struct {
	humacli.CLI
}

// NewApp initializes and returns a new App.
func NewApp() *App {
	mux := http.NewServeMux()

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		cfg := config.Load()
		logger.Init(options.Debug)

		humaConfig := huma.DefaultConfig("SimpleAPI CRD Deployer", "1.0.0")
		humaConfig.Info.Description = "A REST API for deploying SimpleAPI Custom Resource Definitions to Kubernetes clusters"
		humaConfig.Info.Contact = &huma.Contact{
			Name: "API Support",
		}
		humaConfig.Servers = []*huma.Server{
			{URL: fmt.Sprintf("http://localhost:%d", options.Port)},
		}

		k8sClient, err := k8s.NewClient()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create k8s client")
		}
		templateEngine, err := template.NewEngine()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create template engine")
		}
		apiService, err := service.NewAPIService(cfg, k8sClient, templateEngine)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create API service")
		}

		mux := router.RegisterRoutes(mux, humaConfig, apiService)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: middleware.SecurityHeaders(mux),
		}

		hooks.OnStart(func() {
			log.Info().Msgf("Starting server on port %d", options.Port)
			log.Info().Msgf(
				"API documentation available at http://%s:%d/docs",
				options.Host,
				options.Port,
			)

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("Server failed to start")
			}
		})

		hooks.OnStop(func() {
			log.Info().Msg("Shutting down server...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Error().Err(err).Msg("Server forced to shutdown")
			} else {
				log.Info().Msg("Server exited gracefully")
			}
		})

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Info().Msg("Received shutdown signal")
			hooks.OnStop(func() {})
		}()
	})

	return &App{cli}
}
