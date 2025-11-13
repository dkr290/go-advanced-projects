package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"go-lb/internal/config"
	"go-lb/internal/server"
	"go-lb/pkg/logger"
)

func main() {
	cfg := config.Load()

	parseBoolCfg, err := strconv.ParseBool(cfg.DebugLog)
	if err != nil {
		log.Fatal("Something went wrong parsig the bool flag check envs")
	}

	log := logger.New(parseBoolCfg)
	if err := server.Run(cfg, log); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
