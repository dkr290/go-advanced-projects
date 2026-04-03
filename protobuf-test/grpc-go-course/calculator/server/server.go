package main

import (
	"context"
	"log"
	"net"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	proto.UnimplementedCalculateServiceServer
	proto.UnimplementedAvgServiceServer
	proto.UnsafeMaxServiceServer
}

func (s *Server) Calculate(
	ctx context.Context,
	req *proto.CalculateRequest,
) (*proto.CalculateResponse, error) {
	log.Printf(
		"Calculate function on the server is invoked with %d and %d\n",
		req.GetX(),
		req.GetY(),
	)

	sum := req.GetX() + req.GetY()

	return &proto.CalculateResponse{
		Sum: sum,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %v", err)
	}

	log.Printf("Listening on %s\n", addr)

	grpcSerer := grpc.NewServer()

	proto.RegisterCalculateServiceServer(grpcSerer, &Server{})
	proto.RegisterAvgServiceServer(grpcSerer, &Server{})
	proto.RegisterMaxServiceServer(grpcSerer, &Server{})
	if err := grpcSerer.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %c\n", err)
	}
}
