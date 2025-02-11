package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/internals/config"
	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/store"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx     = context.Background()
	dbStore store.Stores
	client  *mongo.Client
)

func init() {
	configs := config.New().
		// WithDbUserName("admin").
		// WithDbPassword("secret").
		Validate()

	mClient, database := repo.NewMongoClient(ctx, configs)
	client = mClient

	userStore := store.NewMongoUserStore(database)
	hotelStore := store.NewMongoHotelStore(database)
	roomStore := store.NewMongoRoomStore(database, hotelStore)
	bookingStore := store.NewMongoBookingStore(database, roomStore)

	dbStore = store.Stores{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	userStore.Drop(ctx)
	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)
	bookingStore.Drop(ctx)
}

func main() {
	defer client.Disconnect(context.TODO())
	admin := fixtures.AddUser(dbStore, "Tan", "Guven", true)
	fmt.Println("admin => ", admin.ID)
	user := fixtures.AddUser(dbStore, "Can", "Guven", false)
	fmt.Println("user => ", user.ID)
	fixtures.AddUser(dbStore, "Fatos", "Guven", false)
	fixtures.AddUser(dbStore, "Leo", "Guven", false)

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

	for _, h := range hotels {
		hotel := fixtures.AddHotel(dbStore, h[0], h[1], randIntGenerator(1, 10), nil)
		rooms := []types.Room{
			{
				Type:      types.FamilyRoomType,
				BasePrice: randomFloatGenerator(100.99, 200.99),
			},
			{
				Type:      types.SuiteRoomType,
				BasePrice: randomFloatGenerator(200.99, 250.99),
			},
			{
				Type:      types.FamilySuitRoomType,
				BasePrice: randomFloatGenerator(250.99, 300.99),
			},
			{
				Type:      types.HoneyMoonRoomType,
				BasePrice: randomFloatGenerator(300.99, 350.99),
			},
			{
				Type:      types.KingRoomType,
				BasePrice: randomFloatGenerator(350.99, 700.99),
			},
		}

		for _, room := range rooms {
			wg.Add(1)
			go func() {
				defer wg.Done()
				insertedRoom := fixtures.AddRoom(dbStore, room.Type, hotel.ID, room.BasePrice)
				if randIntGenerator(1, 10) > 5 {
					booked := fixtures.AddBooking(dbStore, user.ID, insertedRoom.ID.Hex(), time.Now(), time.Now().AddDate(0, 0, randIntGenerator(1, 10)))
					fmt.Println("booking =>", booked.ID)
				}
			}()
		}
	}

	wg.Wait()
}

func randomFloatGenerator(min float64, max float64) float64 {
	randFloat := rand.Float64()*(max-min) + min
	return math.Trunc(randFloat*100) / 100 // right padding 2
}

func randIntGenerator(min int, max int) int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(max-min+1) + min
}
