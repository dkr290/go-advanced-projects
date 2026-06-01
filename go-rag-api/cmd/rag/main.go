package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkr290/go-advanced-projects/go-rag-api/app"
	"github.com/dkr290/go-advanced-projects/go-rag-api/config"
)

func main() {
	// We neeed to
	// - setup the app
	// - set up config
	// - set up llm client
	// - set up Read-Eval-Print loop (REPL)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()
	if err := app.Run(ctx, config.Load()); err != nil {

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
