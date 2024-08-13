package store

import (
	"context"
	"fmt"
	"log"

	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollection = "bookings"

type BookingStore interface {
	Dropper

	InsertBooking(context.Context, *types.BookingParam) (*types.Booking, error)
	GetBookingsByRoomID(context.Context, *types.BookingParam) ([]*types.Booking, error)
	GetBookingsByID(context.Context, string) (*types.Booking, error)
	GetBookings(context.Context) ([]*types.Booking, error)
	CancelBookingByUserID(context.Context, string, primitive.ObjectID) error
	CancelBookingByAdmin(context.Context, string) error
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

func (ms *MongoBookingStore) GetBookingsByRoomID(ctx context.Context, params *types.BookingParam) ([]*types.Booking, error) {
	roomOID, err := primitive.ObjectIDFromHex(params.RoomID)
	if err != nil {
		return nil, err
	}

	cur, err := ms.coll.Find(ctx, bson.M{
		"roomID":   roomOID,
		"fromDate": bson.M{"$gte": params.FromDate},
		"tillDate": bson.M{"$lte": params.TillDate},
	})
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (ms *MongoBookingStore) GetBookingsByID(ctx context.Context, id string) (*types.Booking, error) {
	bookingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking *types.Booking
	if err := ms.coll.FindOne(ctx, bson.M{"_id": bookingID}).Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (ms *MongoBookingStore) GetBookings(ctx context.Context) ([]*types.Booking, error) {
	cur, err := ms.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (ms *MongoBookingStore) CancelBookingByUserID(ctx context.Context, bookingId string, userId primitive.ObjectID) error {
	bookingOID, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		return err
	}
	param := types.CancelBookingParam{Canceled: true}
	resp, err := ms.coll.UpdateOne(ctx, bson.M{"_id": bookingOID, "userID": userId}, bson.M{
		"$set": param.ToBsonMap(),
	})
	if err != nil {
		return err
	}
	if resp.MatchedCount == 0 {
		return fmt.Errorf("booking not found")
	}
	if resp.ModifiedCount == 0 {
		return fmt.Errorf("booking already canceled")
	}

	return nil
}

func (ms *MongoBookingStore) CancelBookingByAdmin(ctx context.Context, bookingId string) error {
	bookingOID, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		return err
	}
	param := types.CancelBookingParam{Canceled: true}
	resp, err := ms.coll.UpdateOne(ctx, bson.M{"_id": bookingOID}, bson.M{"$set": param.ToBsonMap()})
	if err != nil {
		return err
	}
	if resp.MatchedCount == 0 {
		return fmt.Errorf("booking not found")
	}
	if resp.ModifiedCount == 0 {
		return fmt.Errorf("booking already canceled")
	}

	return nil
}

func (ms *MongoBookingStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", bookingCollection)
	return ms.coll.Drop(ctx)
}
