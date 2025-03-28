package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/internals/tokener"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

func TestHandleGetBookings(t *testing.T) {
	config := NewConfig()
	tdb, app := Setup(mDatabase, config)

	var (
		user  = fixtures.AddUser(*tdb.Store, "get", "booking", false)
		hotel = fixtures.AddHotel(*tdb.Store, "bar hotel", "a", 4, nil)
		room  = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
		from  = time.Now()
		till  = from.AddDate(0, 0, 5)
		_     = fixtures.AddBooking(*tdb.Store, user.ID, room.ID.Hex(), from, till)
	)

	t.Run("restrict_for_non_admin_user", func(t *testing.T) {
		token, _ := tokener.GenerateJWT(user.ID.Hex(), user.IsAdmin, config)
		testReq := utils.TestRequest{
			Method: "GET",
			Target: "/v1/admin/bookings",
			Token:  token,
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusForbidden {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}
	})

	t.Run("get_the_bookings_as_admin", func(t *testing.T) {
		admin := fixtures.AddUser(*tdb.Store, "admin", "booking", true)
		token, _ := tokener.GenerateJWT(admin.ID.Hex(), admin.IsAdmin, config)
		testReq := utils.TestRequest{
			Method: "GET",
			Target: "/v1/admin/bookings",
			Token:  token,
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}

		var response *types.ResGeneric
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		data, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("Failed to marshal Data field: %v", err)
		}

		var bookings []*types.Booking
		if err = json.Unmarshal(data, &bookings); err != nil {
			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
		}

		if len(bookings) == 0 {
			t.Fatal("expected booking but received nothing")
		}
	})

	t.Run("get_bookings as user", func(t *testing.T) {
		token, _ := tokener.GenerateJWT(user.ID.Hex(), user.IsAdmin, config)
		testReq := utils.TestRequest{
			Method: "GET",
			Target: "/v1/bookings",
			Token:  token,
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}

		var response *types.ResGeneric
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		data, err := json.Marshal(response.Data)
		if err != nil {
			t.Fatalf("Failed to marshal Data field: %v", err)
		}

		var bookings []*types.Booking
		if err = json.Unmarshal(data, &bookings); err != nil {
			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
		}

		if len(bookings) == 0 {
			t.Fatal("expected booking but received nothing")
		}
	})
}
