package main

import (
	"fmt"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
)

func (s *Server) GreetManyTimes(in *proto.GreetRequest,	str grpc.ServerStreamingServer[proto.GreetResponse]) error {

	log.Printf("Greet many times function was involked %v\n", in)


	for i := range 10 {
		res := fmt.Sprintf("Hello %s number %d", in.FirstName, i)
		err := str.Send(&proto.GreetResponse{

			Result: res,
		})
		if err != nil {
			log.Fatalf("Erro get teh ersponse %v", err)
		}
	}

	return nil
}
