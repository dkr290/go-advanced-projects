package main

import (
	"context"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
)




func deleteBlog(client proto.BlogServiceClient, id string) *proto.BlogId {


	log.Println("Readblog function was running")
	req := &proto.BlogId{Id: id}
	res ,err := client.DeleteBlog(context.Background(), req)
	if err != nil {
    log.Printf("Error happened while deleting %v\n", err)
	}

	log.Printf("Blog was deleted %v\n", res)

	return  req
}

