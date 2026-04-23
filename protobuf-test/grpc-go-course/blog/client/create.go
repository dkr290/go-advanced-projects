package main

import (
	"context"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
)

func createBlog(ctx context.Context, client proto.BlogServiceClient) string {
	log.Println("The createBlog function is started")

	blogRequest := proto.Blog{
		AuthorId: "Clement",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	res, err := client.CreateBlog(context.Background(), &blogRequest)
	if err != nil {
		log.Fatalf("Unexpected error %v\n", err)
	}

	log.Printf("Blog has been created %s\n", res.Id)
	return res.Id
}
