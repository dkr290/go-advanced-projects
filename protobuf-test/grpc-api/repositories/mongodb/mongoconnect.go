package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateMongoClient() (client *mongo.Client, err error) {
	client, err = mongo.Connect(options.Client().ApplyURI("mongodb://mongoadmin:mongoadmin@localhost:27017/"))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb %v", err)
	}
	for i := range 10 {

		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Printf("Error connecting to mongodb %v\n", err)
		}
		i++

		if i > 10 {
			return nil, fmt.Errorf("error connecting to mongodb after 10 retries %v", err)
		}
	}
	log.Println("Connected to mongodb")
	return
}
