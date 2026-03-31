package main

import (
	"io"
	"log"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
)

func (s *Server) Avg(stream grpc.ClientStreamingServer[pb.AvgRequest, pb.AvgResponse]) error {
	log.Println("Avg function")

	var sum int32
	var count int32

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			var avg float32
			if count > 0 {
				avg = float32(sum) / float32(count)
			}
			return stream.SendAndClose(&pb.AvgResponse{
				Avg: avg,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v\n", err)
		}
		log.Println("Receiving", request)
		for _, num := range request.Numbers {
			sum += num
			count++
		}
	}
} 
