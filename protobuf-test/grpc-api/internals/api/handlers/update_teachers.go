package handlers

import (
	"context"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/models"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/repositories/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	client, err := mongodb.CreateMongoClient(s.Log)
	if err != nil {
		return nil, s.Log.Errorf("Error mongodb connection %v", err)
	}
	defer client.Disconnect(ctx)

	teacherMapper := func(pbTeacher *pb.Teacher) *models.Teacher {
		objID, _ := bson.ObjectIDFromHex(pbTeacher.GetId())
		return &models.Teacher{
			ID: objID,
			FirstName: pbTeacher.FirstName,
			LastName:  pbTeacher.LastName,
			Email:     pbTeacher.Email,
			Class:     pbTeacher.Class,
			Subject:   pbTeacher.Subject,
		}
	}

	newTeachers := MapModelToPbModel(req.GetTeachers(), teacherMapper)

	for _, teacher := range newTeachers {
		filter := bson.D{{Key: "_id", Value: teacher.ID}}
		update := bson.D{{Key: "$set", Value: teacher}}
		_,  err = client.Database("school").
			Collection("teachers").
			UpdateOne(ctx, filter, update)
			if err != nil {
			return nil, status.Errorf(codes.Internal, "error updating teacher %v", err)
		}
	}

	return req, nil
}
