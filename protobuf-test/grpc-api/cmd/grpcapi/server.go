package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/internals/api/handlers"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/pkg/config"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/pkg/utils"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/repositories/mongodb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log := utils.New(cfg.DebugFlag)

	_, err = mongodb.CreateMongoClient(log)
	if err != nil {
		log.Error(fmt.Sprintf("Error mongodb connection %v\n", err))
		return
	}
	srv := &handlers.Server{
		Log: log,
		Cfg: cfg,
	}

	s := grpc.NewServer()
	pb.RegisterStudentsServiceServer(s, srv)
	pb.RegisterExecsServiceServer(s, srv)
	pb.RegisterTeachersServiceServer(s, srv)

	reflection.Register(s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to Listen %v", err))
		os.Exit(1)
		log.Error(fmt.Sprintf("failed to listen: %v", err))
		return
	}
	log.Info(fmt.Sprintf("server listening on port %d", cfg.ServerPort))
	if err := s.Serve(lis); err != nil {
		log.Error(fmt.Sprintf("failed to serve: %v", err))
	}
}
