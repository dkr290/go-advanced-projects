package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dkr290/go-advanced-projects/model-serving/model-api/pkg/config"
	"github.com/dkr290/go-advanced-projects/model-serving/model-api/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

var (
	modelsDir      = flag.String("mpath", "models", "Path to the models directory")
	maxConcurrency = flag.Int("maxc", 4, "Maximum number of concurrent processes")
	contextSize    = flag.Int("context-size", 2048, "Model context size")
	gpuLayers      = flag.Int("gpu-layers", 0, "Number of layers to offload to GPU")
	numa           = flag.Bool("numa", false, "Enable NUMA optimization")
	threads        = flag.Int("threads", 4, "Number of threads to use")
	batchSize      = flag.Int("batch-size", 512, "Processing batch size")
	verbose        = flag.Bool("verbose", false, "Enable verbose logging")
)

var sem = make(chan struct{}, *maxConcurrency)

func main() {
	CmdLineparams()
	app := fiber.New()
	if err := os.MkdirAll(*modelsDir, 0755); err != nil {
		panic(err)
	}
	// Create configuration struct for llama parameters
	llamaConfig := &config.LlamaConfig{
		ContextSize: *contextSize,
		GPULayers:   *gpuLayers,
		NUMA:        *numa,
		Threads:     *threads,
		BatchSize:   *batchSize,
		Verbose:     *verbose,
	}

	h := handlers.NewHandlers(*modelsDir, sem, llamaConfig)
	app.Post("api/pull", h.PullModelgguf)
	app.Post("api/chat", h.GenerateRequest)
	app.Get("api/models", h.ListModels)
	app.Delete("api/models/:name", h.DeleteModel)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func CmdLineparams() {
	// Define command-line flags with default values
	help := flag.Bool("help", false, "Show usage information")
	// Parse command-line flags
	flag.Parse()
	// Display help if the flag is set
	if *help {
		showUsage()
	}

	// Ensure required parameters are provided
	if *modelsDir == "models" {
		fmt.Printf("-mpath not supplied using defaults %s.\n", *modelsDir)
	}

	if *maxConcurrency == 4 {
		fmt.Printf("-maxc not supplied using defaults %d.\n", *maxConcurrency)
	}
	if *contextSize == 2048 {
		fmt.Printf("-context-size not supplied using defaults %s.\n", *contextSize)
	}
	if *gpuLayers == 0 {
		fmt.Printf("-gpu-layers using defaults or custom  %s.\n", *contextSize)
	}
	if *numa == false {
		fmt.Printf("-numa not true using defaults %s.\n", *numa)
	}
	if *threads == 4 {
		fmt.Printf("-threads not supplied using defaults %s.\n", *threads)
	}
	if *batchSize == 512 {
		fmt.Printf("-batch-size not supplied using defaults %s.\n", *batchSize)
	}
	if *verbose == false {
		fmt.Printf("-verbose using defaults %s.\n", *verbose)
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
