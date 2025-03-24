package fixtures

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/tnguven/hotel-reservation-app/internal/store"
	"github.com/tnguven/hotel-reservation-app/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// helper to avoid duplication errors
func isDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func AddUser(store store.Stores, fname, lname string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		Email:     fmt.Sprintf("%s_%s@test.com", strings.ToLower(fname), strings.ToLower(lname)),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", strings.ToLower(fname), strings.ToLower(lname)),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil && !isDup(err) {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store store.Stores, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomIDS := rooms
	if roomIDS == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := &types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), hotel)
	if err != nil && !isDup(err) {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store store.Stores, roomType types.RoomType, hotelId primitive.ObjectID, basePrice float64) *types.Room {
	room := &types.Room{
		Type:      roomType,
		HotelID:   hotelId,
		BasePrice: basePrice,
	}

	insertedRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil && !isDup(err) {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(
	store store.Stores,
	uid primitive.ObjectID,
	rid string,
	from, till time.Time,
) *types.Booking {
	booking := &types.BookingParam{
		UserID:      uid,
		RoomID:      rid,
		FromDate:    from,
		TillDate:    till,
		CountPerson: rngInt(1, 8),
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil && !isDup(err) {
		log.Println(err)
		return nil
	}

	return insertedBooking
}

func rngInt(min int, max int) int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Intn(max-min+1) + min
}
