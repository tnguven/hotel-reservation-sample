package store

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollection = "rooms"

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomsByHotelID(context.Context, string) ([]*types.Room, error)
	GetRooms(context.Context, *types.GetRoomsRequest) ([]*types.Room, int64, error)
}

type MongoRoomStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(db *repo.MongoDatabase, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		db:   db.GetDb(),
		coll: db.Coll(roomCollection),

		HotelStore: hotelStore,
	}
}

func (ms *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := ms.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = resp.InsertedID.(primitive.ObjectID)

	if err := ms.HotelStore.PutHotel(ctx, &types.UpdateHotelParams{RoomID: room.ID}, &room.HotelID); err != nil {
		return nil, err
	}

	return room, nil
}

func (ms *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, hotelID string) ([]*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	resp, err := ms.coll.Find(ctx, bson.M{"hotelID": oid})
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (ms *MongoRoomStore) GetRooms(ctx context.Context, qParams *types.GetRoomsRequest) ([]*types.Room, int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{}}},
		bson.D{{Key: "$facet", Value: bson.D{
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$skip", Value: &qParams.PaginationQuery.Page}},
				bson.D{{Key: "$limit", Value: &qParams.PaginationQuery.Limit}},
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
		Data       []*types.Room
		TotalCount int64
	}
	if err := cur.All(ctx, &aggResult); err != nil {
		return nil, 0, err
	}

	if len(aggResult) == 0 {
		return []*types.Room{}, 0, nil
	}

	return aggResult[0].Data, aggResult[0].TotalCount, nil
}

func (ms *MongoRoomStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", roomCollection)
	return ms.coll.Drop(ctx)
}
