package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	mainpb "grpc_streams_client/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	client := mainpb.NewCalculatorServiceClient(conn)

	ctx := context.Background()

	req := mainpb.GenerateFibunacciRequest{
		N: 10,
	}

	stream, err := client.GenerateFibunacci(ctx, &req)
	if err != nil {
		log.Fatalln("Errro calling generate fibunacci x", err)
	}
	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			log.Println("End of stream")
			break
		}
		if err != nil {
			log.Fatalln("Errro receiving data", err)
		}
		log.Println("Fibunacci number is:", resp.GetNumber())
	}

	stream1, err := client.SendNumbers(ctx)
	if err != nil {
		log.Fatalln("Errro error creating stream", err)
	}
	for n := range 9 {
		err := stream1.Send(&mainpb.SendNumbersRequest{
			Number: int32(n),
		})
		if err != nil {
			log.Fatalln("Error sending number:", err)
		}
		time.Sleep(time.Second)
	}

	resp, err := stream1.CloseAndRecv()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp.Sum)
}
