package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ctx        = context.Background()
	hotelStore store.HotelStore
	roomStore  store.RoomStore
)

func init() {
	configs := config.New().
		WithDbUserName("admin").
		WithDbPassword("secret").
		Validate()

	database := db.New(ctx, configs)

	hotelStore = store.NewMongoHotelStore(database)
	roomStore = store.NewMongoRoomStore(database, hotelStore)

	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)
}

func main() {
	hotels := [][]string{
		{"Hilton", "Germany"},
		{"Hilton", "United Kingdom"},
		{"Hilton", "France"},
		{"Hilton", "Ankara"},
		{"Hilton", "Germany"},

		{"Sheraton", "Germany"},
		{"Sheraton", "United Kingdom"},
		{"Sheraton", "France"},
		{"Sheraton", "Ankara"},
		{"Sheraton", "Germany"},

		{"Swissotel", "Germany"},
		{"Swissotel", "United Kingdom"},
		{"Swissotel", "France"},
		{"Swissotel", "Ankara"},
		{"Swissotel", "Germany"},
	}

	for _, h := range hotels {
		seedHotel(h[0], h[1])
	}
}

func seedHotel(hotelName string, location string) {
	hotel := types.Hotel{
		Name:     hotelName,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Type:      types.FamilyRoomType,
			BasePrice: randomFloatGenerator(100.99, 200.99),
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.FamilySuitRoomType,
			BasePrice: randomFloatGenerator(200.99, 250.99),
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.SuiteRoomType,
			BasePrice: randomFloatGenerator(250.99, 300.99),
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.HoneyMoonRoomType,
			BasePrice: randomFloatGenerator(300.99, 350.99),
			HotelID:   insertedHotel.ID,
		},
		{
			Type:      types.KingRoomType,
			BasePrice: randomFloatGenerator(350.99, 700.99),
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

func randomFloatGenerator(min float64, max float64) float64 {
	randFloat := rand.Float64()*(max-min) + min
	return math.Trunc(randFloat*100) / 100 // right padding 2
}
