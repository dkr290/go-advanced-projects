package main

import (
	"fmt"
	"log"
	"net"

	mainpb "grpc_server_stream/proto/gen"

	"google.golang.org/grpc"
)

type server struct {
	mainpb.UnimplementedChatServiceServer
}

func (s *server) Chat(stream mainpb.ChatService_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				resp := &mainpb.ChatResponse{
					Message: "Server received all messagess",
				}
				return stream.SendAndClose(resp)
			}
			return err
		}
		// processed received messagess
		fmt.Printf("Received message from client %s\n", req.GetMessage())

	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	mainpb.RegisterChatServiceServer(grpcServer, &server{})
	log.Println("The server started")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
