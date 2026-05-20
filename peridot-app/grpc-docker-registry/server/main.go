// server/main.go
package main

import (
	"fmt"
	"log"
	"math"
	"net"

	"github.com/dkr290/peridot-app/grpc-docker-registry/config"
	"github.com/dkr290/peridot-app/grpc-docker-registry/internal/services"
	"github.com/dkr290/peridot-app/grpc-docker-registry/internal/storage"
	pb "github.com/dkr290/peridot-app/grpc-docker-registry/proto/gen"
	"github.com/dkr290/peridot-app/grpc-docker-registry/utils"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log := utils.New(cfg.DebugFlag)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to listen: %v", err))
		return
	}

	// Initialize content-addressable blob store
	blobStore := storage.NewFileBlobStore(cfg.StoragePath)

	// Initialize manifest store
	manifestStore := storage.NewFileManifestStore(cfg.StoragePath)

	srv := &services.ImageService{
		Log:           log,
		Cfg:           cfg,
		BlobStore:     blobStore,
		ManifestStore: manifestStore,
	}
	maxMsgSize := grpc.MaxRecvMsgSize(math.MaxInt32)

	s := grpc.NewServer(maxMsgSize)
	pb.RegisterImageRegistryServiceServer(s, srv)
	log.Info(fmt.Sprintf("Server listening at %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		log.Error(fmt.Sprintf("Failed to serve: %v", err))
		return
	}
}
