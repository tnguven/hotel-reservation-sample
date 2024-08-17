package store

import (
	"context"
	"fmt"
	"log"

	"github.com/tnguven/hotel-reservation-app/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelCollection = "hotels"

type HotelStore interface {
	Dropper

	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	PutHotel(context.Context, *types.UpdateHotelParams, *primitive.ObjectID) error
	GetHotels(context.Context) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoHotelStore(db *mongo.Database) *MongoHotelStore {
	return &MongoHotelStore{
		db:   db,
		coll: db.Collection(hotelCollection),
	}
}

func (ms *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := ms.coll.Find(ctx, bson.M{}) // TODO add limit
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (ms *MongoHotelStore) GetHotelByID(ctx context.Context, hotelID string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	var hotel types.Hotel

	if err := ms.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (ms *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := ms.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (ms *MongoHotelStore) PutHotel(ctx context.Context, params *types.UpdateHotelParams, hotelId *primitive.ObjectID) error {
	result, err := ms.coll.UpdateOne(ctx, bson.M{"_id": hotelId}, params.ToBsonMap())
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no hotel found with id %s", hotelId)
	}

	return nil
}

func (ms *MongoHotelStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", hotelCollection)
	return ms.coll.Drop(ctx)
}
