package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoMatchClient(uri string) (*MongoClient, error) {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(timeoutContext, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(timeoutContext, nil); err != nil {
		return nil, fmt.Errorf("error pinging mongo server: %w", err)
	}

	database := client.Database("gohoot")

	return &MongoClient{
		Client:   client,
		Database: database,
	}, nil
}
