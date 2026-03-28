package main

import (
	"context"
	"flag"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = "localhost:50051"
	x    int
	y    int
)

func main() {
	flag.IntVar(&x, "xvar", 0, "set teh value for x")
	flag.IntVar(&y, "yvar", 0, "set teh value for y")

	flag.Parse()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	c := proto.NewCalculateServiceClient(conn)

	x := int32(*&x)
	y := int32(*&y)

	calcReq := proto.CalculateRequest{
		X: x,
		Y: y,
	}

	resp, err := c.Calculate(context.Background(), &calcReq)
	if err != nil {
		log.Fatalf("Cannot create request %v\n", err)
	}

	log.Printf("The sum is %d", resp.Sum)
}
