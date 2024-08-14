package handler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	mid "github.com/tnguven/hotel-reservation-app/handler/middleware"
	"github.com/tnguven/hotel-reservation-app/server"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAuthenticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return user, nil
}

type TestDb struct {
	Store *store.Stores
	db    *mongo.Database
}

type CollectionType string

const (
	UsersColl    CollectionType = "users"
	BookingsColl CollectionType = "bookings"
	HotelsColl   CollectionType = "hotels"
	RoomsColl    CollectionType = "rooms"
)

func (tdb *TestDb) tearDown(t *testing.T, collections []CollectionType) {
	ctx := context.Background()
	var wg sync.WaitGroup
	errChan := make(chan error, len(collections))

	for _, coll := range collections {
		wg.Add(1)
		switch coll {
		case UsersColl:
			go func() {
				defer wg.Done()
				if err := tdb.Store.User.Drop(ctx); err != nil {
					errChan <- err
					return
				}
			}()
		case BookingsColl:
			go func() {
				defer wg.Done()
				if err := tdb.Store.Hotel.Drop(ctx); err != nil {
					errChan <- err
					return
				}
			}()
		case HotelsColl:
			go func() {
				defer wg.Done()
				if err := tdb.Store.Room.Drop(ctx); err != nil {
					errChan <- err
					return
				}
			}()
		case RoomsColl:
			go func() {
				defer wg.Done()
				if err := tdb.Store.Booking.Drop(ctx); err != nil {
					errChan <- err
					return
				}
			}()
		default:
			log.Fatal("unknown collection")
		}
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Setup(db *mongo.Database, withLog bool) (*TestDb, *fiber.App, *Handler, *mid.Validator) {
	hotelStore := store.NewMongoHotelStore(db)
	roomStore := store.NewMongoRoomStore(db, hotelStore)
	bookingStore := store.NewMongoBookingStore(db, roomStore)

	tdb := &TestDb{
		Store: &store.Stores{
			User:    store.NewMongoUserStore(db),
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		},
		db: db,
	}

	return tdb, server.New(withLog), NewHandler(tdb.Store), mid.NewValidator()
}
