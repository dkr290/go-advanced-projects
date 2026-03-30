package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/dkr290/go-advanced-projects/protobuf-test/grpc-payment-app-serverstream/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


var addr = "localhost:50051"

func main()  {
	
	conn, err  := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to the server", err)
	}
	defer conn.Close()

	client := proto.NewPamentServiceClient(conn)

	ctx,cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	req := &proto.PaymentRequest{
		UserId: "Kosta",
	}

	stream , err := client.Payment(ctx, req)
	if err != nil {
		log.Fatalln("Error while calling the payment server",err)
	}
	for {
		msg,err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Stream ended with error", err)
		}
		log.Printf("Payment %s | %.2f %s | %s | %d", msg.PaymentId,msg.Amount,msg.Currency,msg.Status,msg.Timestamp)
	}
}
