package main

import (
	"context"
	"log"
	"os"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/handler"
	"github.com/tnguven/hotel-reservation-app/server"
	"github.com/tnguven/hotel-reservation-app/store"
)

// @description Sample API

// @BasePath /api/v1

// @schemes http https
// @produce application/json
// @consumes application/json

func main() {
	configs := config.New().
		WithDbUserName("admin").
		WithDbPassword("secret").
		WithDbName("hotel_io").
		WithDbCreateIndex(os.Getenv("CREATE_INDEX") == "true").
		Validate()

	ctx := context.Background()
	_, database := db.New(ctx, configs)

	var (
		withLog      = true
		route        = server.New(withLog)
		userStore    = store.NewMongoUserStore(database)
		hotelStore   = store.NewMongoHotelStore(database)
		roomStore    = store.NewMongoRoomStore(database, hotelStore) // TODO refactor this shenanigan
		bookingStore = store.NewMongoBookingStore(database, roomStore)
	)

	handlers := handler.NewHandler(&store.Stores{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	})

	handlers.Register(route, configs)

	log.Fatal(route.Listen(":5000"))
}
