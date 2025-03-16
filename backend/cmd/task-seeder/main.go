package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/tnguven/hotel-reservation-app/db"
	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/store"
	"github.com/tnguven/hotel-reservation-app/internals/types"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		ctx          = context.Background()
		configs      = NewConfig().Validate().Debug()
		mongodb      = repo.NewMongoDatabase(ctx, configs)
		userStore    = store.NewMongoUserStore(mongodb)
		hotelStore   = store.NewMongoHotelStore(mongodb)
		roomStore    = store.NewMongoRoomStore(mongodb, hotelStore)
		bookingStore = store.NewMongoBookingStore(mongodb, roomStore)
	)

	defer func() {
		mongodb.CloseConnection(ctx)
		fmt.Println("Shutting down...")
	}()

	dbStore := store.Stores{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	userStore.Drop(ctx)
	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)
	bookingStore.Drop(ctx)

	admin := fixtures.AddUser(dbStore, "Test", "test", true)
	fmt.Println("admin => ", admin.ID)
	user := fixtures.AddUser(dbStore, "Test1", "Test2", false)
	fmt.Println("user => ", user.ID)
	fixtures.AddUser(dbStore, "Test3", "Test3", false)
	fixtures.AddUser(dbStore, "Test4", "Test4", false)

	hotels := [][]string{
		{"Hilton", "Germany"},
		{"Hilton", "United Kingdom"},
		{"Hilton", "France"},
		{"Hilton", "Ankara"},
		{"Hilton", "Turkey"},
		{"Sheraton", "Germany"},
		{"Sheraton", "United Kingdom"},
		{"Sheraton", "France"},
		{"Sheraton", "Ankara"},
		{"Sheraton", "Turkey"},
		{"Swissotel", "Germany"},
		{"Swissotel", "United Kingdom"},
		{"Swissotel", "France"},
		{"Swissotel", "Ankara"},
		{"Swissotel", "Turkey"},
		{"X", "Germany"},
		{"X", "United Kingdom"},
		{"X", "France"},
		{"X", "Ankara"},
		{"X", "Turkey"},
		{"Y", "Germany"},
		{"Y", "United Kingdom"},
		{"Y", "France"},
		{"Y", "Ankara"},
		{"Y", "Turkey"},
		{"Z", "Germany"},
		{"Z", "United Kingdom"},
		{"Z", "France"},
		{"Z", "Ankara"},
		{"Z", "Turkey"},
	}

	wg := sync.WaitGroup{}
	bookedIds := []string{}

	for _, h := range hotels {
		hotel := fixtures.AddHotel(dbStore, h[0], h[1], rngInt(1, 10), nil)
		rooms := []types.Room{
			{
				Type:      types.FamilyRoomType,
				BasePrice: rngFloat(100.99, 200.99),
			},
			{
				Type:      types.SuiteRoomType,
				BasePrice: rngFloat(200.99, 250.99),
			},
			{
				Type:      types.FamilySuitRoomType,
				BasePrice: rngFloat(250.99, 300.99),
			},
			{
				Type:      types.HoneyMoonRoomType,
				BasePrice: rngFloat(300.99, 350.99),
			},
			{
				Type:      types.KingRoomType,
				BasePrice: rngFloat(350.99, 700.99),
			},
		}

		for _, room := range rooms {
			wg.Add(1)
			go func() {
				defer wg.Done()
				insertedRoom := fixtures.AddRoom(dbStore, room.Type, hotel.ID, room.BasePrice)
				if rngInt(1, 10) > 5 {
					booked := fixtures.AddBooking(
						dbStore,
						user.ID,
						insertedRoom.ID.Hex(),
						time.Now().AddDate(0, 0, rngInt(0, 10)),
						time.Now().AddDate(0, 0, rngInt(11, 22)),
					)
					if booked != nil {
						bookedIds = append(bookedIds, fmt.Sprintf("%s\n", booked.ID.Hex()))
					}
				}
			}()
		}
	}

	wg.Wait()
	fmt.Println(bookedIds)
	db.CreateIndexes(ctx, mongodb.GetDb())
}

func rngFloat(min float64, max float64) float64 {
	randFloat := rand.Float64()*(max-min) + min
	return math.Trunc(randFloat*100) / 100 // right padding 2
}

func rngInt(min int, max int) int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(max-min+1) + min
}
