package main

import (
	"context"
	"io"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
)


func doPrimes(c proto.CalculateServiceClient) {
	log.Println("do Primes was invoked")

	req := &proto.PrimeRequest{
		Number: 833254235405675,
	}

	stream, err := c.Prime(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while prime", err)
	}
	for {
		msg,err := stream.Recv()

		if err == io.EOF{
			break
		}


		if err != nil {
			log.Fatalln("Error reading the stream",err)
		}

		log.Printf("Prime %d\n", msg.GetResult())
	}
}
