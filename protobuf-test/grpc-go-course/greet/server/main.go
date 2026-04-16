package main

import (
	"log"
	"net"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	addr string = "0.0.0.0:50051"
	tls         = true
)

type Server struct {
	pb.UnimplementedGreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %v", err)
	}
	opts := []grpc.ServerOption{}
	if tls {
		certfile := "ssl/server.crt"
		key := "ssl/server.key"
		creds, err := credentials.NewServerTLSFromFile(certfile, key)
		if err != nil {
			log.Fatalf("Failed loading certificates %v", err)
		}
		opts = []grpc.ServerOption{
			grpc.Creds(creds),
		}
	}

	log.Printf("Listening on %s\n", addr)

	server := grpc.NewServer(opts...)
	pb.RegisterGreetServiceServer(server, &Server{})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %c\n", err)
	}
}
