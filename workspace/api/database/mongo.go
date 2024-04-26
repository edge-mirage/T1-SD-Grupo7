package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client  *mongo.Client
	Clients *mongo.Collection
	Users   *mongo.Collection
)

func Init(uri string, database string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	localClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	client = localClient

	Clients = client.Database(database).Collection("clients")
	Users = client.Database(database).Collection("users")

	// Send a ping to confirm a successful connection
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	return err
}

func Close() error {
	return client.Disconnect(context.Background())
}
