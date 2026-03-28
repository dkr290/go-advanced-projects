package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
)

func (s *Server) Greet(
	ctx context.Context,
	in *proto.GreetRequest,
) (*proto.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v\n", in)
	return &proto.GreetResponse{
		Result: fmt.Sprintf("Hello %s", in.GetFirstName()),
	}, nil
}
