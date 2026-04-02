package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/greet/proto"
)


func doGreetEveryone(client pb.GreetServiceClient) {
log.Println("do GreetEveryone was invoked")

stream , err := client.GreetEveryone(context.Background())
if err != nil {
	log.Fatalf("Error while creating stream:", err)
}

reqs := []*pb.GreetRequest{

	{FirstName: "Claude"},
	{FirstName: "Woolf"},
	{FirstName: "Test"},

}

waitC := make(chan struct{})

go func() {
	for _, req := range reqs {
		log.Printf("Send request %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	stream.CloseSend()
}()


go func() {

	for {

	  res := &pb.GreetResponse{}	
		err := stream.RecvMsg(res)
		 

		if err == io.EOF{
			break
		}
		if err != nil {
			log.Printf("Error while receiving  %v\n", err)
			break
		}
		log.Printf("Received %v\n", res.Result)
	}
	close(waitC)



}()

<-waitC

}
