package main

import (
	"context"
	"io"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
)

func doGreetManyTimes(c proto.GreetServiceClient) {
	log.Println("do GreetManyTimes was invoked")

	req := &proto.GreetRequest{
		FirstName: "Clement",
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling greet many times", err)
	}
	for {
		msg,err := stream.Recv()

		if err == io.EOF{
			break
		}


		if err != nil {
			log.Fatalln("Error reading the stream",err)
		}

		log.Printf("Greet many times %s\n", msg.GetResult())
	}
}
