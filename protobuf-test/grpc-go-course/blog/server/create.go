package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBlog(ctx context.Context, in *proto.Blog) (*proto.BlogId, error) {
	log.Printf("The createblog function was invoked with %v", in)

	data := BlogItem{
		AuthorId: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error %v\n", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
     return  nil ,status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to oid %v\n", err))
	}

 return &proto.BlogId{
	 Id: oid.Hex(),
 },nil

}
