package main

import (
	"io"
	"log"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
)

func (s *Server)GreetEveryone(stream grpc.BidiStreamingServer[pb.GreetRequest, pb.GreetResponse]) error{
	log.Println("Receiving request from GreetEveryone")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving request %v", err)
		}

		res := "Hello from the server " + req.FirstName

		err = stream.SendMsg(&pb.GreetResponse{
			Result: res,
		})
    	
		if err != nil {
			log.Fatalf("Error while sending data to the client %v", err)
		}

	}
	
}
