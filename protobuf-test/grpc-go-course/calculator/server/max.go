package main

import (
	"io"
	"log"
	"math"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
)

func (s *Server) Max(steam grpc.BidiStreamingServer[proto.MaxRequest, proto.MaxResponse]) error {
	log.Println("Start the new Max function")
	max := int32(math.MinInt32)
	for {
		req, err := steam.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving reuest %v\n", err)
		}

		if req.Number > max {
			max = int32(req.Number)
			err = steam.Send(&proto.MaxResponse{
				Result: int32(max),
			})
			if err != nil {
				log.Fatalf("Error while sending data to the client %v\n", err)
			}
		}

	}
}
