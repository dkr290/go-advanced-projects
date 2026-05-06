// Package handlers...
package handlers

import (
	"context"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/models"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/repositories/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	// Create a MongoDB client for this request
	client, err := mongodb.CreateMongoClient(s.Log)
	if err != nil {
		return nil, s.Log.Errorf("Error mongodb connection %v", err)
	}
	defer client.Disconnect(ctx)

	// Step 1: Convert protobuf Teacher messages to model Teacher structs
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
	}

	// Step 2: Insert each teacher into MongoDB and build the response
	var insertedTeachers []*pb.Teacher

	for _, teacher := range newTeachers {
		// Insert the teacher document; MongoDB generates an ObjectID
		result, err := client.Database("school").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, s.Log.Errorf("Error inserting teachers %v", err)
		}
		// Extract the MongoDB ObjectID from the insertion result
		objectID, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			return nil, s.Log.Errorf("unexpected ID type: %T", result.InsertedID)
		}

		// Store the generated ID back into the teacher model (optional, for later use)
		teacher.ID = objectID
		// Build the protobuf response teacher with all fields including the new ID
		insertedTeachers = append(insertedTeachers, &pb.Teacher{
			Id:        objectID.Hex(),
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			Subject:   teacher.Subject,
			Class:     teacher.Class,
			Email:     teacher.Email,
		})

	}
	s.Log.Info("Teachers inserted in mongodb")

	return &pb.Teachers{Teachers: insertedTeachers}, nil
}
