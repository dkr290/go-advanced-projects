package main

import (
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
)

func(s *Server) Prime(req *proto.PrimeRequest, stream grpc.ServerStreamingServer[proto.PrimeResponse]) error{

	log.Println("Primes function was invoked with", req)

	number := req.Number
	divisor := int64(2)

	for number >1 {
		if number % divisor == 0 {
			err := stream.Send(&proto.PrimeResponse{
				Result: divisor,
			})
			if err != nil {
				log.Fatalln(err)
			}
			number/= divisor
		} else {
			divisor++
		}
	}

	return nil
}

