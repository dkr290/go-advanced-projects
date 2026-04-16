package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	addr = "localhost:50051"
	tls  = true
)

func main() {
	opts := []grpc.DialOption{}

	if tls {
		certfile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certfile, "localhost")
		if err != nil {
			log.Fatalf("Failed loading certificates %v", err)
		}
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to the server %v", err)
	}
	defer conn.Close()

	client := proto.NewGreetServiceClient(conn)

	greetRequest := proto.GreetRequest{
		FirstName: "Tom",
	}

	doGreet, err := client.Greet(context.Background(), &greetRequest)
	if err != nil {
		log.Fatalf("Error calling greet %v", err)
	}

	fmt.Printf("Greeting: %s\n", doGreet.GetResult())

	// doGreetManyTimes(client)
	//
	// doLongGreet(client)

	doGreetEveryone(client)
}
