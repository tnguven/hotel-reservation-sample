package store

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollection = "bookings"

type BookingStore interface {
	Dropper

	InsertBooking(context.Context, *types.BookingParam) (*types.Booking, error)
}

type MongoBookingStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	RoomStore
}

func NewMongoBookingStore(db *mongo.Database, roomStore RoomStore) *MongoBookingStore {
	return &MongoBookingStore{
		db:   db,
		coll: db.Collection(bookingCollection),

		RoomStore: roomStore,
	}
}

func (ms *MongoBookingStore) InsertBooking(ctx context.Context, params *types.BookingParam) (*types.Booking, error) {
	booking, err := types.NewBookingFromParams(params)
	if err != nil {
		return nil, err
	}

	resp, err := ms.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (ms *MongoBookingStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", bookingCollection)
	return ms.coll.Drop(ctx)
}
