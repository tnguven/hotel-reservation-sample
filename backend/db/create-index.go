package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
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
