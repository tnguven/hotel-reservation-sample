package handler

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/utils"
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

	exitCode := m.Run()

	if err := db.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Failed to close the database connection: %v", err)
	}

	os.Exit(exitCode)
}
