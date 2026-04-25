package main

import (
	"context"
	"fmt"

	pb "github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ListBlog(
	in *emptypb.Empty, stream grpc.ServerStreamingServer[pb.Blog],
) error {
	fmt.Printf("The listblog server side function was invoked %v\n", in)

	// Get all blog collections
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return status.Errorf(codes.Internal, "Unknown internal error %v\n", err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &BlogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				"Error While decoding data from the string %v\n",
				err,
			)
		}
		err = stream.Send(documentToBlog(data))
		if err != nil {
			return status.Errorf(codes.Internal, "Error while sending data %v\n", err)
		}

	}
	if err := cur.Err(); err != nil {
		return status.Errorf(codes.Internal, "Cursor error %v\n", err)
	}

	return nil
}
