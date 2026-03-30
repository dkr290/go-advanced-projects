package main

import (
	"log"
	"net"

	"github.com/dkr290/go-advanced-projects/protobuf-test/grpc-payment-app-serverstream/proto"
	"google.golang.org/grpc"
)


type Server struct{
 proto.UnimplementedPamentServiceServer  	
}


var addr string = "0.0.0.0:50051"

func main()  {

	lis ,err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	log.Println("Server is started on port ", addr)

	server := grpc.NewServer()
	proto.RegisterPamentServiceServer(server, &Server{})

	if err := server.Serve(lis);err != nil {
		log.Fatalf("Failed to serve: %c\n", err)
	}


}
