package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "simple_grpc_server/proto/gen"
	fwpb "simple_grpc_server/proto/gen/firewell"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedCalculateServiceServer
	pb.UnimplementedGreeterServiceServer
	fwpb.UnimplementedAufwiedersehenServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Println("Received request")
	return &pb.AddResponse{
		Sim: req.A + req.B,
	}, nil
}

func main() {
	port := ":50051"
	cert := "cert.pem"
	key := "key.pem"
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen ", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatal("Failed to load the certificate")
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// skipping somehting here
	pb.RegisterCalculateServiceServer(grpcServer, &server{})
	pb.RegisterGreeterServiceServer(grpcServer, &server{})
	fwpb.RegisterAufwiedersehenServer(grpcServer, &server{})
	fmt.Println("Starting the server on port", port)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("Failed to serve", err)
	}
}
