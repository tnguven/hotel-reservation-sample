package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/tnguven/hotel-reservation-app/internal/repo"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const bookingCollection = "bookings"

type BookingStore interface {
	Dropper

	InsertBooking(context.Context, *types.BookingParam) (*types.Booking, error)
	GetBookingsByRoomID(context.Context, *types.BookingParam) ([]*types.Booking, error)
	GetBookingsByID(context.Context, string) (*types.Booking, error)
	GetBookingsAsAdmin(context.Context) ([]*types.Booking, error)
	GetBookingsAsUser(context.Context, *types.User) ([]*types.Booking, error)
	CancelBookingByUserID(context.Context, string, primitive.ObjectID) error
	CancelBookingByAdmin(context.Context, string) error
}

type MongoBookingStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	RoomStore
}

func NewMongoBookingStore(mongodb *repo.MongoDatabase, roomStore RoomStore) *MongoBookingStore {
	return &MongoBookingStore{
		db:   mongodb.GetDb(),
		coll: mongodb.Coll(bookingCollection),

		RoomStore: roomStore,
	}
}

func (ms *MongoBookingStore) InsertBooking(ctx context.Context, params *types.BookingParam) (*types.Booking, error) {
	booking, err := types.NewBookingFromParams(params)
	if err != nil {
		return nil, fmt.Errorf("NewBookingFromParams failed %w", err)
	}

	// Start a session
	session, err := ms.db.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction %w", err)
	}
	defer session.EndSession(ctx)

	var ErrRoomNotAvailable = errors.New("room is not available")

	// Define transaction function
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		roomIsAvailable, err := ms.GetBookingsByRoomID(ctx, params)
		if err != nil {
			return nil, err
		}
		if len(roomIsAvailable) > 0 {
			return nil, ErrRoomNotAvailable
		}

		insertedBooking, err := ms.coll.InsertOne(sessCtx, booking)
		if err != nil {
			return nil, err
		}

		return insertedBooking, nil
	}

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Run transaction
	// WithTransaction will rollback if mongo returns an error
	result, err := session.WithTransaction(
		ctx,
		callback,
		txnOptions,
	)
	if err != nil {
		log.Println("Transaction error:", err)
		if errors.Is(err, ErrRoomNotAvailable) {
			return nil, types.NewError(err, http.StatusConflict, "room is not abailable")
		}
		return nil, types.NewError(err, http.StatusInternalServerError, "error processing booking transaction")
	}

	booking.ID = result.(*mongo.InsertOneResult).InsertedID.(primitive.ObjectID)
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

func (ms *MongoBookingStore) GetBookingsAsAdmin(ctx context.Context) ([]*types.Booking, error) {
	// TODO add pagination
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

func (ms *MongoBookingStore) GetBookingsAsUser(ctx context.Context, user *types.User) ([]*types.Booking, error) {
	cur, err := ms.coll.Find(ctx, bson.M{"userID": user.ID})
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
