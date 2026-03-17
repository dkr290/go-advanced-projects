package main

import (
	"context"
	"log"
	"time"

	mainpipb "grpc_client/proto/gen"
	firewellpb "grpc_client/proto/gen/firewell"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

func main() {
	cert := "cert.pem"

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalln("Failed to load the certificate")
	}

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	)
	if err != nil {
		log.Fatalln("Did not connect", err)
	}

	defer conn.Close()

	client := mainpipb.NewCalculateServiceClient(conn)

	greeterClient := mainpipb.NewGreeterServiceClient(conn)

	// firewClient := firewellpb.NewAufwiedersehenClient(conn)
	firewClient := mainpipb.NewBidFirewellServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	re := mainpipb.AddRequest{
		A: 10,
		B: 20,
	}
	md := metadata.Pairs("authorization", "Bearer=password", "test", "testing")
	ctx = metadata.NewOutgoingContext(ctx, md)
	resp, err := client.Add(ctx, &re)
	if err != nil {
		log.Fatal("Inavalid request", err)
	}

	greetReq := mainpipb.GreetRequest{
		Name: "John",
	}
	greetResp, err := greeterClient.Greet(ctx, &greetReq)
	if err != nil {
		log.Fatal("Inavalid greet resp", err)
	}

	goodByeReq := firewellpb.BidGoodByeRequest{
		Name: "Jane",
	}

	goodByeResp, err := firewClient.BidGoodBye(ctx, &goodByeReq)
	if err != nil {
		log.Fatal("Cannot say goodbye:", err)
	}

	log.Println("The sum is:", resp.GetSim())

	log.Println("The greet message is", greetResp.GetMessage())
	log.Println("Goodbye message is", goodByeResp.GetMessage())
	state := conn.GetState()
	log.Println("Connaction state:", state)
}
