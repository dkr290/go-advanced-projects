package main

import (
	"context"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = "localhost:50051"
	tls  = false
)

func main() {
	var opts []grpc.DialOption

	if tls {
		certfile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certfile, "localhost")
		if err != nil {
			log.Fatalf("Failed loading certificates %v", err)
		}
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	} else {
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to the server %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)

	id := createBlog(context.Background(), client)
	ReadBlog(client, id)
	ReadBlog(client, "nonExistingID")
	err = updateBlogClient(client, id)
	if err != nil {
		log.Fatalf("errror %v\n", err)
	}
	doList(client)
	deleteBlog(client, id)
}
