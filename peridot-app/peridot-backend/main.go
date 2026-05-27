package main
import (
	"log"
	"net/http"

	"github.com/humaio/huma/v4"
	"github.com/humaio/huma/v4/adapters/humago"
	"github.com/dkr290/peridot-app/peridot-backend/api"
)

func main() {
	// Create HUMA API with default configuration
	config := huma.DefaultConfig("Peridot API", "1.0.0")
	
	// Register custom operation processor if needed
	config.OperationProcessors = append(config.OperationProcessors, func(operation *huma.Operation) error {
		return nil
	})

	// Create router and adapter
	router := http.NewServeMux()
	api := humago.New(router, config)

	// Register routes
	api.Register(api.NewImageHandler())

	// Start server
	log.Println("🚀 Starting Peridot API on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

