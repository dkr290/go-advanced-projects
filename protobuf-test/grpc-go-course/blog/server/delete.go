package main

import (
	"context"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)




	func(s *Server)DeleteBlog(ctx context.Context, in *pb.BlogId) (*emptypb.Empty, error){

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
	filter := bson.M{"_id": oid}

	// 3. Execute FindOne
	res,err  := collection.DeleteOne(ctx, filter)
	if err != nil{
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot delete the blog: %v", err,
		)
	}
	if res.DeletedCount == 0 {
    return nil,status.Errorf(codes.Internal, "Cannot find blog Id to delete")
	}
return &emptypb.Empty{},nil



	}
