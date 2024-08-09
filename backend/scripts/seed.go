package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()
	configs := config.New().
		WithDbUserName("admin").
		WithDbPassword("secret").
		Validate()

	database := db.New(ctx, configs)
	hotelStore := store.NewMongoHotelStore(database)
	roomStore := store.NewMongoRoomStore(database, hotelStore)

	hotel := types.Hotel{
		Name:     "Hilton",
		Location: "France",
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Type:      types.SingleBedRoomType,
			BasePrice: 99.9,
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.DoubleBedRoomType,
			BasePrice: 112.9,
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.SuiteRoomType,
			BasePrice: 222.9,
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.KingSuiteRoomType,
			BasePrice: 333.9,
			HotelID:   insertedHotel.ID,
		},
	}

	for _, room := range rooms {
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

}
