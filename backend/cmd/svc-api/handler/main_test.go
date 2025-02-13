package handler_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/internals/config"
	"github.com/tnguven/hotel-reservation-app/internals/must"
	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	mDatabase *repo.MongoDatabase
	configs   *config.Configs
)

func TestMain(m *testing.M) {
	var exitCode int
	ctx := context.Background()
	mongoDBContainer := must.Panic(mongodb.Run(ctx, "mongo:8"))
	defer func() {
		if err := testcontainers.TerminateContainer(mongoDBContainer); err != nil {
			fmt.Printf("failed to terminate container: %s", err)
		}
		os.Exit(exitCode)
	}()
	endpoint := must.Panic(mongoDBContainer.ConnectionString(ctx))
	configs = config.New().
		WithMongoDbURI(endpoint).
		WithDbName("hotel_io_test").
		Validate().
		Debug()
	mDatabase = utils.NewDb(configs)

	list := must.Panic(mDatabase.GetDb().ListCollectionNames(context.TODO(), bson.M{}))
	db.CreateIndexes(context.TODO(), mDatabase.GetDb()) // need to create the indexes to prevent unique duplications

	// Cleanup the existing docs before testing
	// we can't drop the collections because we need the created indexes
	for _, coll := range list {
		must.Panic(mDatabase.Coll(coll).DeleteMany(context.TODO(), bson.M{}))
	}

	exitCode = m.Run()

	if err := mDatabase.GetDb().Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	mDatabase.CloseConnection(context.TODO())
}
