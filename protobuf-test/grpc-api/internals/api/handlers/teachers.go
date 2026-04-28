package handlers

import (
	"context"

	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/repositories/mongodb"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	client, err := mongodb.CreateMongoClient(s.Log)
	if err != nil {
		return nil, s.Log.Errorf("Error mongodb connection %v", err)
	}
	defer client.Disconnect(ctx)
}
