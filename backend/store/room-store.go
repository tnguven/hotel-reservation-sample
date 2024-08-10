package store

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollection = "rooms"

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, string) ([]*types.Room, error)
}

type MongoRoomStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(db *mongo.Database, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		db:   db,
		coll: db.Collection(roomCollection),

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

func (ms *MongoRoomStore) GetRooms(ctx context.Context, hotelID string) ([]*types.Room, error) {
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

func (ms *MongoRoomStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", roomCollection)
	return ms.coll.Drop(ctx)
}
