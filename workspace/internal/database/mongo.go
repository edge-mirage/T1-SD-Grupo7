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
	Token   *mongo.Collection
)

func Init(uri string, database string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	localClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	client = localClient

	Clients = client.Database(database).Collection("clients")
	Users = client.Database(database).Collection("users")
	Token = client.Database(database).Collection("token")

	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	return err
}

func Close() error {
	return client.Disconnect(context.Background())
}
