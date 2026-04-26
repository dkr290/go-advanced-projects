package handlers


import pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"


type Server struct{
	pb.UnimplementedExecsServiceServer
	pb.UnimplementedStudentsServiceServer
	pb.UnimplementedTeachersServiceServer
}
