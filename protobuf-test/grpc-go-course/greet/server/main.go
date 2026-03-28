package main

import (
	"log"
	"net"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	pb.UnimplementedGreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %v", err)
	}

	log.Printf("Listening on %s\n", addr)

	server := grpc.NewServer()
	pb.RegisterGreetServiceServer(server, &Server{})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %c\n", err)
	}
}
