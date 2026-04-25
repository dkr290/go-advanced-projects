package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
)

func updateBlogClient(client pb.BlogServiceClient, id string) error {
	fmt.Println("The update blog client function")

	blog := &pb.Blog{
		Id:       id,
		AuthorId: "Not Clement",
		Title:    "A new title",
		Content:  "Content Of the first blog with some awsome additions",
	}

	_, err := client.UpdateBlog(context.Background(), blog)
	if err != nil {
		return fmt.Errorf("Error happening while updating %v\n", err)
	}

	log.Println("Blog was updated")

	return nil
}
