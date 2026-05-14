package helpers

import (
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/models"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"go.mongodb.org/mongo-driver/v2/bson"

)

func TeacherMapper(pbTeacher *pb.Teacher) *models.Teacher {
	objID, _ := bson.ObjectIDFromHex(pbTeacher.GetId())
	return &models.Teacher{
		ID:        objID,
		FirstName: pbTeacher.FirstName,
		LastName:  pbTeacher.LastName,
		Email:     pbTeacher.Email,
		Class:     pbTeacher.Class,
		Subject:   pbTeacher.Subject,
	}
}
// MapModelToPbModel converts any slice of model T to slice of proto U
func MapModelToPbModel[T, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

