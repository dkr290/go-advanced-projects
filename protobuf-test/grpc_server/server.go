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
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedCalculateServiceServer
	pb.UnimplementedGreeterServiceServer
	//	fwpb.UnimplementedAufwiedersehenServer
	pb.UnimplementedBidFirewellServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Println("Received request")
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("no metadata received")
	}
	log.Println("The meta is:", meta)

	val, ok := meta["authorization"]
	if !ok {
		log.Println("no value with auth key in metadata")
	}
	log.Println("Authorization", val)

	return &pb.AddResponse{
		Sim: req.A + req.B,
	}, nil
}

func (s *server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Println("Received greet request")
	return &pb.GreetResponse{
		Message: fmt.Sprintf("Hello %s.", req.Name),
	}, nil
}

func (s *server) BidGoodBye(
	ctx context.Context,
	req *fwpb.BidGoodByeRequest,
) (*fwpb.BidGoodByeResponse, error) {
	log.Println("Received bidgoodbye request")
	return &fwpb.BidGoodByeResponse{
		Message: fmt.Sprintf("%s.", req.Name),
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
	// fwpb.RegisterAufwiedersehenServer(grpcServer, &server{})
	pb.RegisterBidFirewellServiceServer(grpcServer, &server{})
	fmt.Println("Starting the server on port", port)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal("Failed to serve", err)
	}
}
