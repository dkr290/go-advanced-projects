package main

import (
	"context"
	"log"
	"time"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
)

func doLongGreet(c pb.GreetServiceClient) {
	log.Println("do LongGreet was invoked")


	greetReqs := []*pb.GreetRequest{
		{FirstName: "Clement"},
		{FirstName: "Marrie"},
		{FirstName: "Test"},
	}

	stream ,err := c.LogGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling long greet %v\n",err)
	}
	for _, req := range greetReqs{
		log.Printf("sending Req %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res ,err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reeiving responce frpom long greet %v\n", err)
	}

	log.Printf("Long greet %s\b", res.Result)

}

