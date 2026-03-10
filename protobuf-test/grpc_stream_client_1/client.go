package main

import (
	"context"
	"log"
	"time"

	mainpb "grpc_client_stream_1/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Cannot create grpc client ", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := mainpb.NewChatServiceClient(conn)
	// Start the streaming RPC
	chatClient, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	waitChan := make(chan struct{})

	// Send some predefined messages
	messages := []string{
		"Hello from client!",
		"This is a test message.",
		"Another message here.",
		"Final message.",
	}

	go func() {
		for _, msg := range messages {
			log.Println("sending message", msg)
			req := mainpb.ChatRequest{
				Message: msg,
			}
			if err := chatClient.Send(&req); err != nil {
				log.Fatalf("Error sending the message %v\n", msg)
			}
			time.Sleep(500 * time.Millisecond)
		}
		chatClient.CloseSend()
		close(waitChan)
	}()
	<-waitChan
}
