package mongodb

import (
	"context"
	"fmt"
	"os"

	"github.com/dkr290-go-advanced-projects/protobuf-test/grpc-api/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateMongoClient(log utils.Logger) (client *mongo.Client, err error) {
	client, err = mongo.Connect(
		options.Client().ApplyURI("mongodb://mongoadmin:mongoadmin@localhost:27017/"),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb %v", err)
	}
	for i := range 10 {

		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Error(fmt.Sprintf("Error connecting to mongodb %v\n", err))
			os.Exit(1)
		}
		i++

		if i > 10 {
			return nil, fmt.Errorf("error connecting to mongodb after 10 retries %v", err)
		}
	}
	log.Info("Connected to mongodb")
	return
}
