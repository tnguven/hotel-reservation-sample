package fixtures

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tnguven/hotel-reservation-app/internals/store"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store store.Stores, uid primitive.ObjectID, rid string, from, till time.Time) *types.Booking {
	booking := &types.BookingParam{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
