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
	"google.golang.org/grpc/encoding/gzip"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
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

	chatClient, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Error creasting chat stream %v\n", err)
	}
	waitc := make(chan struct{})

	// send message

	go func() {
		message := []string{"Hello", "How are you ?", "GoodBye"}
		for _, mes := range message {
			log.Println("Sending message", message)
			err := chatClient.Send(&mainpb.ChatRequest{Message: mes})
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(3 * time.Second)
		}
		chatClient.CloseSend()
	}()

	go func() {
		for {
			res, err := chatClient.Recv()
			if err != io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data %v", err)
			}
			log.Printf("Received response %v", res.GetMessage())
		}
		close(waitc)
	}()
	<-waitc
}
