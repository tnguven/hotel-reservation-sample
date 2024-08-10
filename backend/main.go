package main

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/handler"
	"github.com/tnguven/hotel-reservation-app/router"
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
		Validate()

	ctx := context.Background()

	database := db.New(ctx, configs)

	withLog := true
	route := router.New(withLog)

	var (
		userStore  = store.NewMongoUserStore(database)
		hotelStore = store.NewMongoHotelStore(database)
		roomStore  = store.NewMongoRoomStore(database, hotelStore) // TODO refactor this shenanigan
	)

	handlers := handler.NewHandler(&store.Stores{
		Hotel: hotelStore,
		Room:  roomStore,
		User:  userStore,
	})

	handlers.Register(route)

	log.Fatal(route.Listen(":5000"))
}
