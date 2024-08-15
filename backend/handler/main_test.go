package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db      *mongo.Database
	configs *config.Configs
)

func TestMain(m *testing.M) {
	configs = config.New().
		WithDbUserName("admin").
		WithDbPassword("secret").
		WithDbCreateIndex(true).
		WithDbName("hotel_io_test").
		Validate()

	var client *mongo.Client
	client, db = utils.NewDb(configs)

	list, err := db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(list)

	// CLeanup the existing docs before testing
	// we can't drop the collections because we need the created indexes
	for _, coll := range list {
		db.Collection(coll).DeleteMany(context.TODO(), bson.M{})
	}

	exitCode := m.Run()

	if err := db.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Failed to close the database connection: %v", err)
	}

	os.Exit(exitCode)
}
