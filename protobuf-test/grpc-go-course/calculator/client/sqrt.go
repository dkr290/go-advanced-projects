package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func doSqrt(c proto.SqrtServiceClient, n int32) {
	fmt.Println("doSqrt was invoked")
	res, err := c.Sqrt(context.Background(), &proto.SqrtRequest{
		Number: n,
	})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			log.Printf("Error message from server %v\n", e.Message())
			log.Printf("Error code from server %v\n", e.Code())

			if e.Code() == codes.InvalidArgument{
         log.Println("It probably negative nbumber has been send ")
				 return
			}

		} else {
			log.Fatalf("Non grpc error %v\n", err)
		}
	}

	log.Printf("Sqrt %f\n", res.GetResult())
}
