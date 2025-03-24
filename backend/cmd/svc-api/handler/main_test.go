package handler_test

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/internal/configure"
	"github.com/tnguven/hotel-reservation-app/internal/must"
	"github.com/tnguven/hotel-reservation-app/internal/repo"
	"github.com/tnguven/hotel-reservation-app/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type TestConfigs struct {
	configure.Common
	configure.DbConfig
	configure.Server
	configure.Secrets
	configure.Session

	mongoDbURI   string
	mongoDbName  string
	jwtSecret    string
	tokenExpHour int64
	listenAddr   string
	log          bool
	env          string
}

func (conf *TestConfigs) WithMongoDbURI(dbURI string) *TestConfigs {
	conf.mongoDbURI = dbURI
	return conf
}

func NewConfig() *TestConfigs {
	return &TestConfigs{
		mongoDbName:  cmp.Or(os.Getenv("MONGO_DATABASE"), "hotel_io_test"),
		jwtSecret:    "top_secret",
		tokenExpHour: 2,
		listenAddr:   ":5000",
		env:          "test",
		log:          true,
	}
}

func (conf *TestConfigs) DbName() string {
	return conf.mongoDbName
}

func (conf *TestConfigs) DbURI() string {
	return conf.mongoDbURI
}

func (conf *TestConfigs) DbUriWithDbName() string {
	return fmt.Sprintf("%s/%s", conf.DbURI(), conf.DbName())
}

func (conf *TestConfigs) JWTSecret() string {
	return conf.jwtSecret
}

func (conf *TestConfigs) TokenExpHour() int64 {
	return conf.tokenExpHour
}

func (conf *TestConfigs) ListenAddr() string {
	return conf.listenAddr
}

func (conf *TestConfigs) GoEnv() string {
	return conf.env
}

func (conf *TestConfigs) WithLog() bool {
	return conf.log
}

var (
	mDatabase *repo.MongoDatabase
)

func TestMain(m *testing.M) {
	var exitCode int
	ctx := context.Background()
	mongoDBContainer := must.Panic(
		mongodb.Run(ctx, "mongo:8", mongodb.WithReplicaSet("rs_test")), // necessary for transactions
	)
	defer func() {
		if err := mongoDBContainer.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate container: %s", err)
		}
		os.Exit(exitCode)
	}()

	endpoint := must.Panic(mongoDBContainer.ConnectionString(ctx))
	testConfig := NewConfig().WithMongoDbURI(endpoint)
	mDatabase = utils.NewDb(testConfig)

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
