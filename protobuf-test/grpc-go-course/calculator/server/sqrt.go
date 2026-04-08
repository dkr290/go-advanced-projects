package main

import (
	"context"
	"log"
	"math"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Sqrt(ctx context.Context, in *proto.SqrtRequest) (*proto.SqrtResponse, error) {
	log.Printf("Sqrt was iunvoked with %v\n", in)

	number := in.Number
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "received negative number")
	}

	return &proto.SqrtResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
}
