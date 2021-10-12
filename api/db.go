package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var mongoOnce sync.Once
var ConnString = `mongodb+srv://nekinci:%2BNiyazi678%2B@cluster0.p7nvm.mongodb.net/ContainerService?retryWrites=true&w=majority`
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
