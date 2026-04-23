package main

import (
	"context"
	"log"
	"net"

	"github.com/dkr290-go-advanced-projects/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var (
	addr       string = "0.0.0.0:50051"
	collection *mongo.Collection
)

type Server struct {
	proto.UnimplementedBlogServiceServer
}

func connectMongo() *mongo.Client {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI("mongodb://mongoadmin:mongoadmin@localhost:27017/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	mongoCl := connectMongo()
	collection = mongoCl.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %v", err)
	}

	log.Printf("Listening on %s\n", addr)

	grpcSerer := grpc.NewServer()

	proto.RegisterBlogServiceServer(grpcSerer, &Server{})
	if err := grpcSerer.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %c\n", err)
	}
}
