package main

import (
	"context"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
)



func ReadBlog(client proto.BlogServiceClient, id string) *proto.Blog {


	log.Println("Readblog function was running")
	req := &proto.BlogId{Id: id}
	res ,err := client.ReadBlog(context.Background(), req)
	if err != nil {
    log.Printf("Error happened while reading %v\n", err)
	}

	log.Printf("Blog was read %v\n", res)

	return  res
	
}

