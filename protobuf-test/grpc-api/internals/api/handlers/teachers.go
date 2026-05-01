package handlers

import (
	"context"
	"fmt"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/models"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/repositories/mongodb"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	client, err := mongodb.CreateMongoClient(s.Log)
	if err != nil {
		return nil, s.Log.Errorf("Error mongodb connection %v", err)
	}
	defer client.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))

	for i, pbTeacher := range req.GetTeachers() {
		mt := models.Teacher{
			FirstName: pbTeacher.FirstName,
			LastName:  pbTeacher.LastName,
			Email:     pbTeacher.Email,
			Class:     pbTeacher.Class,
			Subject:   pbTeacher.Subject,
		}

		newTeachers[i] = &mt
		fmt.Println(mt)
	}
	// // Insert into MongoDB (adjust database/collection names)
	// collection := client.Database("school").Collection("teachers")
	// insertResult, err := collection.InsertMany(ctx, newTeachers)
	// if err != nil {
	// 	return nil, s.Log.Errorf("Error inserting teachers: %v", err)
	// }
	//
	// // Optionally update IDs in the response
	// // for i, id := range insertResult.InsertedIDs { ... }
	//
	return nil, nil
}
