package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dkr290/go-advanced-projects/api-service/foundation/logger"
	"go.uber.org/zap"
)

var build = "develop"

func main() {

	log, err := logger.New("SALES-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}

}

func run(log *zap.SugaredLogger) error {
	// GOMAXPROX
	log.Infow("startup", "GOMAXPROX", runtime.GOMAXPROCS(0), "BUILD", build)
	//--------------------------------------------------------------
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown

	log.Infow("shutdown", "status", "shutdown started", "signal", sig)
	defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}
