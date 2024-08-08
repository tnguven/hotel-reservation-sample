package db

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoClient(ctx context.Context, config config.Configs) *mongo.Client {
	clientOpts := options.Client().ApplyURI(config.DbURI)

	if config.DbUserName != "" && config.DbPassword != "" {
		clientOpts.SetAuth(options.Credential{
			Username: config.DbUserName,
			Password: config.DbPassword,
		})
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to the mongodb")

	if err := createIndexes(ctx, client.Database(config.DbName)); err != nil {
		log.Fatal(err)
	}

	return client
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1}, // Index in ascending order on the 'email' field
		},
		Options: options.Index().SetUnique(true),
	}

	indexModels := []mongo.IndexModel{emailIndexModel}

	for _, model := range indexModels {
		_, err := collection.Indexes().CreateOne(ctx, model)
		if err != nil {
			log.Println("Could not create index:", err)
			return err
		}
	}

	log.Println("Created index fields")
	return nil
}

func New(ctx context.Context, config config.Configs) *mongo.Database {
	client := getMongoClient(ctx, config)
	return client.Database(config.DbName)
}
