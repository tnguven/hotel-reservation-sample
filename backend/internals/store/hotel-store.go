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
	GetHotels(context.Context, *types.HotelQueryParams) ([]*types.Hotel, int64, error)
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

func (ms *MongoHotelStore) GetHotels(ctx context.Context, qParams *types.HotelQueryParams) ([]*types.Hotel, int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{}}},
		bson.D{{Key: "$facet", Value: bson.D{
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$skip", Value: &qParams.Page}},
				bson.D{{Key: "$limit", Value: &qParams.Limit}},
			}},
			{Key: "totalCount", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "data", Value: 1},
			{Key: "totalCount", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$totalCount.count", 0}},
			}},
		}}},
	}
	cur, err := ms.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	var aggResult []struct {
		Data       []*types.Hotel `bson:"data"`
		TotalCount int64          `bson:"totalCount"`
	}

	if err := cur.All(ctx, &aggResult); err != nil {
		return nil, 0, err
	}

	if len(aggResult) == 0 {
		return []*types.Hotel{}, 0, nil
	}

	return aggResult[0].Data, aggResult[0].TotalCount, nil
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
