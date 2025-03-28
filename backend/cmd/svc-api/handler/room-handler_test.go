package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/internals/tokener"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

func TestHandleGetRooms(t *testing.T) {
	config := NewConfig()
	tdb, app := Setup(mDatabase, config)

	var (
		user    = fixtures.AddUser(*tdb.Store, "test", "getRooms", false)
		hotel   = fixtures.AddHotel(*tdb.Store, "bar hotel", "a", 4, nil)
		room1   = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
		_       = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
		_       = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
		from    = time.Now().AddDate(0, 0, 1)
		till    = time.Now().AddDate(0, 0, 10)
		booking = fixtures.AddBooking(
			*tdb.Store,
			user.ID,
			room1.ID.Hex(),
			from,
			till,
		)
	)

	token, _ := tokener.GenerateJWT(user.ID.Hex(), user.IsAdmin, config)

	const target = "/v1/rooms"

	t.Run("get all the rooms with status occupied", func(t *testing.T) {

		t.Logf("booking: %+v", booking)
		testReq := utils.TestRequest{
			Method: "GET",
			Target: target,
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

		t.Logf("response: %s", string(data))

		var rooms []*types.Room
		if err = json.Unmarshal(data, &rooms); err != nil {
			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
		}

		if len(rooms) == 0 {
			t.Fatal("expected rooms but received nothing")
		}
		if len(rooms) != 1 {
			t.Fatalf("expected 1 occupied room, got %d", len(rooms))
		}

		if rooms[0].ID.Hex() != room1.ID.Hex() {
			t.Fatalf("expected room id: %s, but got: %s", room1.ID.Hex(), rooms[0].ID.Hex())
		}
	})

	t.Run("pagination of rooms", func(t *testing.T) {
		testReq := utils.TestRequest{
			Method: "GET",
			Target: target + "?limit=2&page=1",
			Token:  token,
		}

		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
		}

		var response *types.ResWithPaginate[types.ResNumericPaginate]
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Pagination.Limit != 2 {
			t.Fatalf("expected limit 2, got: %d", response.Pagination.Limit)
		}

		if response.Pagination.Page != 1 {
			t.Fatalf("expected page 1, got: %d", response.Pagination.Page)
		}

		if response.Pagination.Count != 3 {
			t.Fatalf("expected 1 count, got: %d", response.Pagination.Count)
		}
	})
}

// func TestHandleGetRooms(t *testing.T) {
// 	config := NewConfig()
// 	tdb, app := Setup(mDatabase, config)

// 	var (
// 		hotel = fixtures.AddHotel(*tdb.Store, "bar hotel", "a", 4, nil)
// 		room1 = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
// 		_     = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
// 		_     = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
// 	)

// 	const target = "/v1/rooms"

// 	t.Run("get all the rooms", func(t *testing.T) {
// 		testReq := utils.TestRequest{
// 			Method: "GET",
// 			Target: target,
// 		}

// 		resp, err := app.Test(testReq.NewRequestWithHeader())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if resp.StatusCode != http.StatusOK {
// 			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
// 		}

// 		var response *types.GenericResponse
// 		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 			t.Fatal(err)
// 		}

// 		data, err := json.Marshal(response.Data)
// 		if err != nil {
// 			t.Fatalf("Failed to marshal Data field: %v", err)
// 		}

// 		var rooms []*types.Room
// 		if err = json.Unmarshal(data, &rooms); err != nil {
// 			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
// 		}

// 		if len(rooms) == 0 {
// 			t.Fatal("expected rooms but received nothing")
// 		}
// 		if len(rooms) != 1 {
// 			t.Fatalf("expected 1 occupied room, got %d", len(rooms))
// 		}

// 		if rooms[0].ID.Hex() != room1.ID.Hex() {
// 			t.Fatalf("expected room id: %s, but got: %s", room1.ID.Hex(), rooms[0].ID.Hex())
// 		}
// 	})

// 	t.Run("pagination of rooms", func(t *testing.T) {
// 		testReq := utils.TestRequest{
// 			Method: "GET",
// 			Target: target + "?limit=2&page=1",
// 		}

// 		resp, err := app.Test(testReq.NewRequestWithHeader())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if resp.StatusCode != http.StatusOK {
// 			t.Fatalf("expected 200 status code but received %d", resp.StatusCode)
// 		}

// 		var response *types.GenericResponse
// 		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 			t.Fatal(err)
// 		}

// 		if response.PaginationResponse == nil {
// 			t.Fatalf("expected pagination response")
// 		}

// 		if response.PaginationResponse.Limit != 2 {
// 			t.Fatalf("expected limit 2, got: %d", response.PaginationResponse.Limit)
// 		}
// 		if response.PaginationResponse.Page != 1 {
// 			t.Fatalf("expected page 1, got: %d", response.PaginationResponse.Page)
// 		}
// 		if response.PaginationResponse.Count != 1 {
// 			t.Fatalf("expected 1 count, got: %d", response.PaginationResponse.Count)
// 		}
// 	})
// }

// func TestHandleBookRoom(t *testing.T) {
// 	config := NewConfig()
// 	tdb, app := Setup(mDatabase, config)

// 	var (
// 		user     = fixtures.AddUser(*tdb.Store, "book", "room", false)
// 		hotel    = fixtures.AddHotel(*tdb.Store, "bar hotel", "a", 4, nil)
// 		room     = fixtures.AddRoom(*tdb.Store, types.FamilyRoomType, hotel.ID, 10.99)
// 		token, _ = tokener.GenerateJWT(user.ID.Hex(), user.IsAdmin, config)
// 		target   = fmt.Sprintf("/v1/rooms/%s/booking", room.ID.Hex())
// 	)

// 	t.Run("book a room successfully", func(t *testing.T) {
// 		from := time.Now().AddDate(0, 0, 1)
// 		till := from.AddDate(0, 0, 5)
// 		params := types.BookingParam{
// 			FromDate:    from,
// 			TillDate:    till,
// 			CountPerson: 2,
// 		}
// 		b, _ := json.Marshal(params)
// 		testReq := utils.TestRequest{
// 			Method:  "POST",
// 			Target:  target,
// 			Token:   token,
// 			Payload: bytes.NewReader(b),
// 		}
// 		resp, err := app.Test(testReq.NewRequestWithHeader())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if resp.StatusCode != fiber.StatusCreated {
// 			t.Fatalf("expected 201 status code but received %d", resp.StatusCode)
// 		}

// 		var response *types.GenericResponse
// 		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 			t.Fatal(err)
// 		}

// 		data, err := json.Marshal(response.Data)
// 		if err != nil {
// 			t.Fatalf("Failed to marshal Data field: %v", err)
// 		}

// 		var booking *types.Booking
// 		if err = json.Unmarshal(data, &booking); err != nil {
// 			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
// 		}

// 		if booking.UserID != user.ID {
// 			t.Fatalf("expected user id %s, got %s", user.ID.Hex(), booking.UserID.Hex())
// 		}

// 		if booking.RoomID != room.ID {
// 			t.Fatalf("expected room id %s, got %s", room.ID.Hex(), booking.RoomID.Hex())
// 		}
// 	})
// 	t.Run("Validation", func(t *testing.T) {
// 		from := time.Now().AddDate(0, 0, -1)
// 		till := from.AddDate(0, 0, 5)

// 		emptyInput := &types.BookingParam{}
// 		pastInput := &types.BookingParam{
// 			FromDate:    from,
// 			TillDate:    till,
// 			CountPerson: 2,
// 		}

// 		invalidCount := &types.BookingParam{
// 			FromDate:    time.Now().AddDate(0, 0, 1),
// 			TillDate:    time.Now().AddDate(0, 0, 5),
// 			CountPerson: 0,
// 		}

// 		type test struct {
// 			desc     string
// 			input    *types.BookingParam
// 			expected *types.Error
// 			status   int
// 		}
// 		tests := []test{
// 			{
// 				desc:  "should return all required fields error",
// 				input: emptyInput,
// 				expected: &types.Error{
// 					GenericResponse: &types.GenericResponse{
// 						Status: 400,
// 						Msg:    "Bad Request",
// 						Errors: map[string]interface{}{
// 							"FromDate":  "required",
// 							"TillDate":  "required",
// 							"NumPerson": "required",
// 						},
// 					},
// 				},
// 				status: 400,
// 			},
// 			{
// 				desc:  "should return past input error",
// 				input: pastInput,
// 				expected: &types.Error{
// 					GenericResponse: &types.GenericResponse{
// 						Status: 400,
// 						Msg:    "Bad Request",
// 						Errors: map[string]interface{}{
// 							"bookingRoomRequest": "cannot book a room in the past",
// 						},
// 					},
// 				},
// 				status: 400,
// 			},
// 			{
// 				desc:  "should return invalid count error",
// 				input: invalidCount,
// 				expected: &types.Error{
// 					GenericResponse: &types.GenericResponse{
// 						Status: 400,
// 						Msg:    "Bad Request",
// 						Errors: map[string]interface{}{
// 							"NumPerson": "min - invalid",
// 						},
// 					},
// 				},
// 				status: 400,
// 			},
// 		}
// 		for _, tc := range tests {
// 			b, _ := json.Marshal(tc.input)
// 			testReq := utils.TestRequest{
// 				Method:  "POST",
// 				Target:  target,
// 				Token:   token,
// 				Payload: bytes.NewReader(b),
// 			}
// 			resp, reqErr := app.Test(testReq.NewRequestWithHeader())
// 			if reqErr != nil {
// 				t.Fatal(reqErr)
// 			}

// 			t.Run(fmt.Sprintf("should return %d status code", tc.status), func(t *testing.T) {
// 				t.Parallel()
// 				if resp.StatusCode != tc.status {
// 					t.Errorf("expected status code %d but return %d", tc.status, resp.StatusCode)
// 				}
// 			})

// 			t.Run(tc.desc, func(t *testing.T) {
// 				t.Parallel()
// 				var body types.Error
// 				if errDecode := json.NewDecoder(resp.Body).Decode(&body); errDecode != nil {
// 					t.Fatal(errDecode)
// 				}

// 				if !reflect.DeepEqual(&body, tc.expected) {
// 					t.Errorf("expected response %+v, but got %+v", tc.expected, body)
// 				}
// 			})
// 		}
// 	})

// 	t.Run("unauth", func(t *testing.T) {
// 		from := time.Now().AddDate(0, 0, 1)
// 		till := from.AddDate(0, 0, 5)
// 		params := types.BookingParam{
// 			FromDate:    from,
// 			TillDate:    till,
// 			CountPerson: 2,
// 		}
// 		b, _ := json.Marshal(params)
// 		testReq := utils.TestRequest{
// 			Method:  "POST",
// 			Target:  target,
// 			Payload: bytes.NewReader(b),
// 		}
// 		resp, err := app.Test(testReq.NewRequestWithHeader())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if resp.StatusCode != fiber.StatusUnauthorized {
// 			t.Fatalf("expected 401 status code but received %d", resp.StatusCode)
// 		}
// 	})
// 	t.Run("should return invalid room id", func(t *testing.T) {
// 		from := time.Now().AddDate(0, 0, 1)
// 		till := from.AddDate(0, 0, 5)
// 		params := types.BookingParam{
// 			FromDate:    from,
// 			TillDate:    till,
// 			CountPerson: 2,
// 		}
// 		b, _ := json.Marshal(params)
// 		testReq := utils.TestRequest{
// 			Method:  "POST",
// 			Target:  fmt.Sprintf("/v1/rooms/%s/booking", primitive.NewObjectID().Hex()),
// 			Token:   token,
// 			Payload: bytes.NewReader(b),
// 		}
// 		resp, err := app.Test(testReq.NewRequestWithHeader())
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if resp.StatusCode != fiber.StatusInternalServerError {
// 			t.Fatalf("expected 500 status code but received %d", resp.StatusCode)
// 		}
// 	})
// }
