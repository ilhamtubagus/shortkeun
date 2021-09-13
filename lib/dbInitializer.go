package lib

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func InitDatabaseClient() *mongo.Client {
	if mongoClient == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
		if err != nil {
			log.Fatal("Error while initializing database " + err.Error())
		}
		mongoClient = client
		return mongoClient
	}
	return mongoClient
}
