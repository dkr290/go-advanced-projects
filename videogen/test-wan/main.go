package main

import (
	"os"

	"github.com/wan2-video-server/cmd"
	"github.com/wan2-video-server/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Execute root command
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
		os.Exit(1)
	}
}
