package handlers

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DecodeCursorToProto function to decode MongoDB cursor to any protobuf type
func DecodeCursorToProto[T any, P any](ctx context.Context,	cur *mongo.Cursor,	mapper func(*T) *P) ([]*P, error) {
	var results []*P

	for cur.Next(ctx) {
		data := new(T)
		err := cur.Decode(data)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"Error while decoding data: %v",
				err,
			)
		}
		results = append(results, mapper(data))
	}

	if err := cur.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Cursor error: %v", err)
	}

	// Return empty slice instead of nil for proper serialization
	if results == nil {
		results = []*P{}
	}

	return results, nil
}

// MapModelToPbModel converts any slice of model T to slice of proto U
func MapModelToPbModel[T, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

