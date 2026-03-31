package main

import (
	"context"
	"log"
	"time"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
)

func doAvgNumbers(c pb.AvgServiceClient) {
	log.Println("do average  was invoked")

	nums := []*pb.AvgRequest{
		{Number: 32},
		{Number: 28},
		{Number: 3},
		{Number: 19},
	}

	stream, err := c.Avg(context.Background())
	if err != nil {
		log.Fatalf("Error while calling avg %v\n", err)
	}
	for _, num := range nums {
		log.Printf("sending number %v\n", num)
		stream.Send(num)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reeiving responce from avg function %v\n", err)
	}

	log.Printf("The avarage result is %f\b", res.Result)
}
