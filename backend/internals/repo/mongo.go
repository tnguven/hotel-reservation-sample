package repo

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/internals/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func getMongoClient(ctx context.Context, config *config.Configs) *mongo.Client {
	if mongoClient != nil {
		return mongoClient
	}

	clientOpts := options.Client().ApplyURI(config.DbURI)

	if config.DbUserName != "" && config.DbPassword != "" {
		clientOpts.SetAuth(options.Credential{
			Username: config.DbUserName,
			Password: config.DbPassword,
		})
	}

	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to the mongodb")

	return mongoClient
}

func NewMongoClient(ctx context.Context, config *config.Configs) (*mongo.Client, *mongo.Database) {
	client := getMongoClient(ctx, config)
	return client, client.Database(config.DbName)
}
