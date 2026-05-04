package handlers

import (
	"context"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/models"
	pb "github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/proto/gen"
	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/repositories/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetTeachers(
	ctx context.Context,
	req *pb.GetTeachersRequest,
) (*pb.Teachers, error) {
	// Create a MongoDB client for this request
	client, err := mongodb.CreateMongoClient(s.Log)
	if err != nil {
		return nil, s.Log.Errorf("Error mongodb connection %v", err)
	}
	defer client.Disconnect(ctx)

	// filtering , getting the filters from the request
	filter := dynamicFilter(req)
  opts := buildSortOptions(req.SortBy)	
// Build Dynamic Filter from non empty request fields



	cur, err := client.Database("school").Collection("teachers").Find(ctx, filter,opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unknown internal error %v\n", err)
	}
	defer cur.Close(context.Background())
	var teachers []*pb.Teacher
	// Decode the data from mongodb and rthen populate fields in protobuf
	for cur.Next(ctx) {
		data := &models.Teacher{}
		err := cur.Decode(data)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"Error While decoding data from the string %v\n",
				err,
			)
		}
		teachers = append(teachers, &pb.Teacher{
			Id:        data.ID.Hex(),
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     data.Email,
			Class:     data.Class,
			Subject:   data.Subject,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Cursor error: %v", err)
	}
	// Return empty slice instead of nil for proper serialization
	if teachers == nil {
		teachers = []*pb.Teacher{}
	}
	s.Log.Info("Teachers fetched from mongodb")
	return &pb.Teachers{Teachers: teachers}, nil
}

func dynamicFilter(req *pb.GetTeachersRequest) bson.M {
	filter := bson.M{}
	if req.Teacher != nil {
		if req.Teacher.FirstName != "" {
			filter["first_name"] = req.Teacher.FirstName
		}
		if req.Teacher.LastName != "" {
			filter["last_name"] = req.Teacher.LastName
		}
		if req.Teacher.Email != "" {
			filter["email"] = req.Teacher.Email
		}
		if req.Teacher.Class != "" {
			filter["class"] = req.Teacher.Class
		}
		if req.Teacher.Subject != "" {
			filter["subject"] = req.Teacher.Subject
		}
	}

	return filter
}

func buildSortOptions(sortFields []*pb.SortField) *options.FindOptionsBuilder {
	var sortOptions bson.D
	if len(sortFields) == 0 {
		return nil
	}
	for _, sf := range sortFields {
		order := 1
		if sf.GetOrder() == pb.Order_ORDER_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: sf.Field, Value: order})

	}
	return options.Find().SetSort(sortOptions)
}
