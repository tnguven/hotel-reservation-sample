package handler

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	mid "github.com/tnguven/hotel-reservation-app/handler/middleware"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
)

func TestHandleGetBookings(t *testing.T) {
	tdb, app, handlers, _ := setup(db, false)

	app.Get("/admin", mid.JWTAuthentication(tdb.Store.User, configs), mid.WithAdminAuth, handlers.HandleGetBookingsAsAdmin)
	app.Get("/user", mid.JWTAuthentication(tdb.Store.User, configs), handlers.HandleGetBookingsAsUser)

	var (
		user  = fixtures.AddUser(*tdb.Store, "get", "booking", false)
		hotel = fixtures.AddHotel(*tdb.Store, "bar hotel", "a", 4, nil)
		room  = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
		from  = time.Now()
		till  = from.AddDate(0, 0, 5)
		_     = fixtures.AddBooking(*tdb.Store, user.ID, room.ID.Hex(), from, till)
	)

	t.Run("restrict_for_non_admin_user", func(t *testing.T) {
		testReq := utils.TestRequest{
			Method: "GET",
			Target: "/admin",
			Token:  utils.GenerateJWT(user.ID.Hex(), user.IsAdmin, configs),
		}

		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusForbidden {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}
	})

	t.Run("no restriction_for_admin_user", func(t *testing.T) {
		var (
			admin = fixtures.AddUser(*tdb.Store, "admin", "booking", true)
		)

		testReq := utils.TestRequest{
			Method: "GET",
			Target: "/admin",
			Token:  utils.GenerateJWT(admin.ID.Hex(), admin.IsAdmin, configs),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}

		var bookings []*types.Booking
		if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
			t.Fatal(err)
		}

		if len(bookings) == 0 {
			t.Fatal("expected booking got nothing")
		}
	})

	t.Run("get_bookings as user", func(t *testing.T) {
		testReq := utils.TestRequest{
			Method: "GET",
			Target: "/user",
			Token:  utils.GenerateJWT(user.ID.Hex(), user.IsAdmin, configs),
		}

		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}

		var bookings []*types.Booking
		if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
			t.Fatal(err)
		}

		if len(bookings) == 0 {
			t.Fatal("expected booking got nothing")
		}
	})
}
