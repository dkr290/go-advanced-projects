package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func doList(c pb.BlogServiceClient) {
	fmt.Println("doList was invoked")

	stream, err := c.ListBlog(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog: %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v\n", err)
		}

		log.Printf("Blog: %v\n", msg)
	}
}

