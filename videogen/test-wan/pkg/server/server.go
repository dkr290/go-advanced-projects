package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wan2-video-server/pkg/config"
	"github.com/wan2-video-server/pkg/handlers"
	"github.com/wan2-video-server/pkg/logger"
	"github.com/wan2-video-server/pkg/middleware"
	"github.com/wan2-video-server/pkg/model"
)

// Server represents the HTTP server
type Server struct {
	config      *config.Config
	router      *gin.Engine
	httpServer  *http.Server
	modelEngine model.Engine
	log         *logger.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	log := logger.NewLogger()

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Initialize model engine
	var modelEngine model.Engine
	var err error

	if cfg.Python.Enabled {
		modelEngine, err = model.NewPythonEngine(cfg)
	} else {
		modelEngine, err = model.NewLocalEngine(cfg)
	}

	if err != nil {
		log.Fatalf("Failed to initialize model engine: %v", err)
	}

	return &Server{
		config:      cfg,
		router:      router,
		modelEngine: modelEngine,
		log:         log,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Setup middleware
	s.setupMiddleware()

	// Setup routes
	s.setupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler:        s.router,
		ReadTimeout:    s.config.Process.RequestTimeout,
		WriteTimeout:   s.config.Process.RequestTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in a goroutine
	go func() {
		s.log.Infof("Server starting on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	s.waitForShutdown()

	return nil
}

func (s *Server) setupMiddleware() {
	// Logger middleware
	s.router.Use(middleware.Logger(s.log))

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// CORS middleware
	s.router.Use(middleware.CORS())

	// Rate limiting
	s.router.Use(middleware.RateLimiter(s.config.Process.MaxConcurrentRequests))
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", handlers.HealthCheck)
	s.router.GET("/", handlers.Index)

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// Model info
		v1.GET("/model/info", handlers.NewModelInfoHandler(s.modelEngine))

		// Video generation endpoints
		videoHandler := handlers.NewVideoHandler(s.modelEngine, s.config)
		v1.POST("/generate/text-to-video", videoHandler.TextToVideo)
		v1.POST("/generate/image-to-video", videoHandler.ImageToVideo)
		v1.POST("/generate/video-to-video", videoHandler.VideoToVideo)

		// Job status
		v1.GET("/job/:id", handlers.NewJobHandler(s.modelEngine).GetJobStatus)

		// Model management
		modelHandler := handlers.NewModelManagementHandler(s.config)
		v1.GET("/models", modelHandler.ListModels)
		v1.POST("/models/download", modelHandler.DownloadModel)
	}

	// Static files for outputs
	s.router.Static("/outputs", "./outputs")
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Fatalf("Server forced to shutdown: %v", err)
	}

	s.log.Info("Server exited")
}
