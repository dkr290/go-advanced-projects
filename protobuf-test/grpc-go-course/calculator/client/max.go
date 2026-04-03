package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
)

func doMax(client pb.MaxServiceClient) {
	fmt.Println("DoMax was invoked")

	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("Error while creating the stream: %v\n", err)
	}

	reqs := []*pb.MaxRequest{
		{Number: 1},
		{Number: 5},
		{Number: 3},
		{Number: 6},
		{Number: 2},
		{Number: 30},
	}

	waitC := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Println("Send request ", req)
			err := stream.Send(req)
			if err != nil {
				fmt.Println("error sending number", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	var recap []int32
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error receiving from server", err)
			}

			log.Println("Received", resp.Result)
			recap = append(recap, resp.Result)
			
		}

		close(waitC)
	}()

	<-waitC

	log.Println("Numbers received ", recap)
}
