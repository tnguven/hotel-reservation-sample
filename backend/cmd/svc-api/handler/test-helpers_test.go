package handler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/cmd/svc-api/handler"
	mid "github.com/tnguven/hotel-reservation-app/internal/middleware"
	"github.com/tnguven/hotel-reservation-app/internal/repo"
	"github.com/tnguven/hotel-reservation-app/internal/server"
	"github.com/tnguven/hotel-reservation-app/internal/store"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"github.com/tnguven/hotel-reservation-app/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestDb struct {
	Store *store.Stores
	db    *mongo.Database
}

type CollectionType string

func (tdb *TestDb) TearDown(t *testing.T) {
	ctx := context.Background()
	errChan := make(chan error, 4)
	events := []func(){
		func() {
			if err := tdb.Store.User.Drop(ctx); err != nil {
				errChan <- err
			}
		},
		func() {
			if err := tdb.Store.Hotel.Drop(ctx); err != nil {
				errChan <- err
			}
		},
		func() {
			if err := tdb.Store.Room.Drop(ctx); err != nil {
				errChan <- err
			}
		},
		func() {
			if err := tdb.Store.Booking.Drop(ctx); err != nil {
				errChan <- err
			}
		},
	}

	for event := range utils.Parallel(events) {
		event()
	}

	close(errChan)

	for err := range errChan {
		if err != nil {
			t.Fatal(err)
		}
	}
}

// mimic the real implementation to test the integration also
func Setup(db *repo.MongoDatabase, configs *TestConfigs) (*TestDb, *fiber.App) {
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
		db: db.GetDb(),
	}

	app := server.NewServer(configs)

	validator, _ := mid.NewValidator()
	handlers := handler.NewHandler(tdb.Store)
	handlers.Register(app, configs, validator)

	return tdb, app
}

func GetAutheVnticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return user, nil
}
