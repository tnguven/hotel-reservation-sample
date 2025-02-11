package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tnguven/hotel-reservation-app/cmd/api/handler"
	"github.com/tnguven/hotel-reservation-app/internals/config"
	"github.com/tnguven/hotel-reservation-app/internals/middleware"
	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/server"
	"github.com/tnguven/hotel-reservation-app/internals/store"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

// @description Sample API

// @BasePath /api/v1

// @schemes http https
// @produce application/json
// @consumes application/json

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}

	configs := config.New().
		WithEnv(env).
		// WithDbUserName("admin").
		// WithDbPassword("secret").
		Validate()

	var (
		rootCtx          = context.Background()
		client, database = repo.NewMongoClient(rootCtx, configs)
		route            = server.NewServer(configs.Log, configs.Env)
		userStore        = store.NewMongoUserStore(database)
		hotelStore       = store.NewMongoHotelStore(database)
		roomStore        = store.NewMongoRoomStore(database, hotelStore) // TODO refactor this shenanigan
		bookingStore     = store.NewMongoBookingStore(database, roomStore)
	)

	handlers := handler.NewHandler(&store.Stores{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	})

	validator, err := middleware.NewValidator()
	if err != nil {
		panic(err)
	}

	handlers.Register(route, configs, validator)

	go func() {
		if err := route.Listen(configs.Port); err != nil && err != http.ErrServerClosed {
			log.Panicf("⚠️ server listen error: %s", err)
		}
	}()

	utils.GraceFullyShutDown(rootCtx, func(shutdownCtx context.Context) {
		defer func() {
			client.Disconnect(rootCtx)
		}()

		if err := route.Shutdown(); err != nil {
			log.Printf("server shutdown failed: %s", err)
		}
	})
}
