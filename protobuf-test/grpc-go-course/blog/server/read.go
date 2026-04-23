package main

import (
	"context"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (s *Server) ReadBlog(ctx context.Context, in *proto.BlogId) (*proto.Blog, error) {
	log.Printf("ReadBlog was invoked with %v\n", in)

	// 1. Convert the string ID from the request to a MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}

	// 2. Create a struct to unpack the MongoDB document
	// Usually, this matches your proto definition but includes BSON tags
	data := &BlogItem{}
	filter := bson.M{"_id": oid}

	// 3. Execute FindOne
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot find blog with specified ID: %v", err,
		)
	}

	// 4. Map the custom struct to your gRPC proto message

	return documentToBlog(data),nil
	}
