package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dkr290/go-advanced-projects/model-serving/model-api/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

const (
	modelsDir      = "models"
	llamaCPPPath   = "llama.cpp/build/bin/llama-cli" // Update this path
	maxConcurrency = 4
)

var sem = make(chan struct{}, maxConcurrency)

func main() {
	app := fiber.New()
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		panic(err)
	}
	h := handlers.NewHandlers(modelsDir, sem, llamaCPPPath)
	app.Post("api/pull", h.PullModelgguf)
	app.Post("api/chat", h.GenerateRequest)
	app.Get("api/models", h.ListModels)
	app.Delete("api/models/:name", h.DeleteModel)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func CmdLineparams() (modelsDir *string, llamaCPPPath *string, maxConcurrency *int) {
	// Define command-line flags with default values
	modelsDir = flag.String("mpath", "models", "Path to the models directory")
	llamaCPPPath = flag.String(
		"llmbin",
		"llama.cpp/build/bin/llama-cli",
		"Path to the LlamaCPP binary",
	)
	maxConcurrency = flag.Int("maxc", 4, "Maximum number of concurrent processes")
	help := flag.Bool("help", false, "Show usage information")
	// Parse command-line flags
	flag.Parse()
	// Display help if the flag is set
	if *help {
		showUsage()
	}

	// Ensure required parameters are provided
	if *modelsDir == "" {
		fmt.Printf("-mpath not supplied using defaults %s.\n", *modelsDir)
		showUsage()
	}
	if *llamaCPPPath == "" {
		fmt.Printf("-llmbin not supplied using defaults %s.\n", *llamaCPPPath)
		showUsage()
	}
	if *maxConcurrency == 0 {
		fmt.Printf("-maxc not supplied using defaults %d.\n", *maxConcurrency)
		showUsage()
	}
	return
}

// showUsage prints usage information and exits
func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  -modelspath string        Path to the models directory (default \"models\")")
	fmt.Println(
		"  -llamabinary string     Path to the LlamaCPP binary (default \"llama.cpp/build/bin/llama-cli\")",
	)
	fmt.Println("  -maxconcurrency int      Maximum number of concurrent processes (default 4)")
	fmt.Println("  -help                    Show this help message")
	os.Exit(0)
}
