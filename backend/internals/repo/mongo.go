package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/tnguven/hotel-reservation-app/internals/config"
	"github.com/tnguven/hotel-reservation-app/internals/must"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mdClient *mongo.Client

type MongoDatabase struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDatabase(ctx context.Context, conf *config.Configs) *MongoDatabase {
	if mdClient == nil {
		mdClient = must.Panic(mongo.Connect(
			ctx,
			options.Client().ApplyURI(fmt.Sprintf("%s/%s", conf.MongoDbURI, conf.MongoDbName)),
		))
	}

	db := &MongoDatabase{
		client: mdClient,
		db:     mdClient.Database(conf.MongoDbName),
	}

	db.CheckConnection(ctx)

	return db
}

func (d *MongoDatabase) CloseConnection(ctx context.Context) {
	log.Println("shutting mongodb")
	if err := d.client.Disconnect(ctx); err != nil {
		log.Printf("close mongodb failed: %v", err)
	}
}

func (d *MongoDatabase) Coll(name string) *mongo.Collection {
	coll := d.db.Collection(name)
	return coll
}

func (d *MongoDatabase) CheckConnection(ctx context.Context) {
	// Send a ping to confirm a successful connection
	if err := d.db.RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}

	log.Println("✅ connected to MongoDB")
}

func (d *MongoDatabase) GetDb() *mongo.Database {
	return d.db
}
