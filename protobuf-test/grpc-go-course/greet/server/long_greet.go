package main

import (
	"fmt"
	"io"
	"log"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
)

func (s *Server) LogGreet(stream grpc.ClientStreamingServer[pb.GreetRequest, pb.GreetResponse]) error {
	log.Println("LongGreet function is")

	res :=""

for {
	req ,err := stream.Recv()
	if err == io.EOF{
		return stream.SendAndClose(&pb.GreetResponse{
			Result: res,
		})
	}
	if err != nil {
		log.Fatalf("Error while reading client stream %v\n", err)
	}
	log.Println("Receiving", req)
	res += fmt.Sprintf("Hello %s!\n", req.GetFirstName())
}

}
