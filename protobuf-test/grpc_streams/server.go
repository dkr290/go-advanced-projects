package main

import (
	"log"
	"net"
	"time"

	mainpb "grpc_stream/proto/gen"

	"google.golang.org/grpc"
)

type server struct {
	mainpb.UnimplementedCalculatorServiceServer
}

func (s *server) GenerateFibunacci(
	req *mainpb.GenerateFibunacciRequest,
	stream mainpb.CalculatorService_GenerateFibunacciServer,
) error {
	n := req.N
	a, b := 0, 1
	for i := 0; i < int(n); i++ {
		err := stream.Send(&mainpb.GenerateFibunacciResponse{
			Number: int32(a),
		})
		if err != nil {
			return err
		}
		a, b := b, a+b
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	mainpb.RegisterCalculatorServiceServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
