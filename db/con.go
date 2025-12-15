package db

import (
	"ciphera-api/env"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() error {
	mongoURI := env.GetDatabaseServerUri()
	fmt.Println(mongoURI)
	if mongoURI == "" {
		return errors.New("database uri is required")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(`error occurred while connecting to database.`)
		fmt.Println(err)
		return errors.New("mongo connection failed")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return errors.New("mongo ping failed")
	}
	log.Println("connected to mongodb")
	Client = client
	return nil
}

func Collection(collection string) *mongo.Collection {
	return Client.Database(env.GetDatabaseName()).Collection(collection)
}
