package main

import (
	"context"
	"log"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	log.Printf("Updating blog function is here %v\n", in)

	// 1. Convert the string ID from the request to a MongoDB ObjectID
	oid, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot parse ID",
		)
	}
	data := &BlogItem{
		AuthorId: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": data}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Could not find record to update",
		)
	}

	if res.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Cannot find blog with id")
	}

	return &emptypb.Empty{}, nil
}
