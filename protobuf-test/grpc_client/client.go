package main

import (
	"context"
	"log"
	"time"

	mainpipb "grpc_client/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	)
	if err != nil {
		log.Fatalln("Did not connect", err)
	}

	defer conn.Close()

	client := mainpipb.NewCalculateServiceClient(conn)

	greeterClient := mainpipb.NewGreeterServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	re := mainpipb.AddRequest{
		A: 10,
		B: 20,
	}
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

	log.Println("The sum is:", resp.GetSim())

	log.Println("The greet message is", greetResp.GetMessage())
	state := conn.GetState()
	log.Println("Connaction state:", state)
}
