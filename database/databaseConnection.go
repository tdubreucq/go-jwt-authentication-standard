package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = dbConnection()

func dbConnection() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to Mongodb.")

	return client
}

func dbDisconnection(mongoClient *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient.Disconnect(ctx)
}

func OpenCollection(mongoClient *mongo.Client, collectionName string) *mongo.Collection {
	return mongoClient.Database("gogame").Collection(collectionName);
}