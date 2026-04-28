package handlers

import (
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/pkg/config"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/pkg/utils"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
)

type Server struct {
	pb.UnimplementedExecsServiceServer
	pb.UnimplementedStudentsServiceServer
	pb.UnimplementedTeachersServiceServer
	Log utils.Logger
	Cfg *config.Config // or your config type
}
