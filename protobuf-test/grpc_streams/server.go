package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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
		log.Println("send number: ", a)
		a, b = b, a+b
		time.Sleep(time.Second)
	}
	return nil
}

func (s *server) SendNumbers(stream mainpb.CalculatorService_SendNumbersServer) error {
	var sum int32
	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mainpb.SendNumbersResponse{
				Sum: sum,
			})
		}
		if err != nil {
			return err
		}
		log.Println("Number:", req.GetNumber())
		sum += req.GetNumber()
	}
}

func (s *server) Chat(stream mainpb.CalculatorService_ChatServer) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		// receiving messagess from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Println("Received message:", req.GetMessage())

		fmt.Print("Enter response:")
		str, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input := strings.TrimSpace(str)

		// sending messagess to the stream
		err = stream.Send(&mainpb.ChatResponse{
			Message: input,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	mainpb.RegisterCalculatorServiceServer(grpcServer, &server{})
	log.Println("The server started")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
