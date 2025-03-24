package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/tnguven/hotel-reservation-app/cmd/svc-api/handler"
	"github.com/tnguven/hotel-reservation-app/internal/middleware"
	"github.com/tnguven/hotel-reservation-app/internal/must"
	"github.com/tnguven/hotel-reservation-app/internal/repo"
	"github.com/tnguven/hotel-reservation-app/internal/server"
	"github.com/tnguven/hotel-reservation-app/internal/store"
	"github.com/tnguven/hotel-reservation-app/internal/utils"
)

// @description Sample API

// @BasePath /api/v1

// @schemes http https
// @produce application/json
// @consumes application/json

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("warn can not load .env file")
	}

	var (
		rootCtx      = context.Background()
		configs      = NewConfig().Validate().Debug()
		mongodb      = repo.NewMongoDatabase(rootCtx, configs)
		route        = server.NewServer(configs)
		userStore    = store.NewMongoUserStore(mongodb)
		hotelStore   = store.NewMongoHotelStore(mongodb)
		roomStore    = store.NewMongoRoomStore(mongodb, hotelStore) // TODO refactor this shenanigan
		bookingStore = store.NewMongoBookingStore(mongodb, roomStore)
	)

	handlers := handler.NewHandler(&store.Stores{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	})

	validator := must.Panic(middleware.NewValidator())
	handlers.Register(route, configs, validator)

	go func() {
		if err := route.Listen(configs.ListenAddr()); err != nil && err != http.ErrServerClosed {
			log.Panicf("⚠️ server listen error: %s", err)
		}
	}()

	utils.GraceFullyShutDown(rootCtx, func(shutdownCtx context.Context) {
		defer func() {
			mongodb.CloseConnection(rootCtx)
		}()

		if err := route.Shutdown(); err != nil {
			log.Printf("server shutdown failed: %s", err)
		}
	})
}
