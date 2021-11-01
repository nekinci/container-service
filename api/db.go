package api

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
)

var mongoOnce sync.Once
var ConnString = os.Getenv("MONGO_CONN_STRING")
var clientInstance *mongo.Client

func GetMongoClient() *mongo.Client {

	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(ConnString)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatalf("%v", err)
		}
		clientInstance = client
	})

	return clientInstance
}

func GetDatabase(name string) *mongo.Database {
	return GetMongoClient().Database(name)
}

func GetContainerDatabase() *mongo.Database {
	return GetDatabase("ContainerService")
}
