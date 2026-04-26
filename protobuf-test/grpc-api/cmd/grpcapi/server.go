package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/internals/api/handlers"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStudentsServiceServer(s, &handlers.Server{})
	pb.RegisterExecsServiceServer(s, &handlers.Server{})
	pb.RegisterTeachersServiceServer(s, &handlers.Server{})

	reflection.Register(s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening on port %d", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
