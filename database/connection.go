package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"
)

var (
	database *mongo.Database
)

// Connect with database
func Connect(db string, url string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
		return nil

	} else {
		database = client.Database(db)
		return client

	}
}

func Get() *mongo.Database {
	return database
}
