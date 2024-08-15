package db

import (
	"context"
	"log"
	"sync"

	"github.com/tnguven/hotel-reservation-app/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoClient(ctx context.Context, config *config.Configs) *mongo.Client {
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

	if config.CreateIndex {
		log.Println("creating indexes...")
		if err := createIndexes(ctx, client.Database(config.DbName)); err != nil {
			log.Fatal(err)
		}
		log.Println("indexes are created")
	}

	return client
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)
	go createBookingIndexes(ctx, db, &wg, errChan)
	go createUsersIndexes(ctx, db, &wg, errChan)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func createBookingIndexes(ctx context.Context, db *mongo.Database, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	bookingCollection := db.Collection("bookings")
	roomIDIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "roomID", Value: 1}, // Index in ascending order
		},
	}
	userIDIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userID", Value: 1}, // Index in ascending order
		},
	}

	indexModels := []mongo.IndexModel{roomIDIndexModel, userIDIndexModel}

	for _, model := range indexModels {
		_, err := bookingCollection.Indexes().CreateOne(ctx, model)
		if err != nil {
			errChan <- err
			return
		}
	}

	log.Println("Created index bookings.roomID, bookings.userID fields")
}

func createUsersIndexes(ctx context.Context, db *mongo.Database, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	userCollection := db.Collection("users")
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1}, // Index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	indexModels := []mongo.IndexModel{emailIndexModel}

	for _, model := range indexModels {
		_, err := userCollection.Indexes().CreateOne(ctx, model)
		if err != nil {
			errChan <- err
			return
		}
	}

	log.Println("Created index users.email fields")
}

func New(ctx context.Context, config *config.Configs) (*mongo.Client, *mongo.Database) {
	client := getMongoClient(ctx, config)
	return client, client.Database(config.DbName)
}
