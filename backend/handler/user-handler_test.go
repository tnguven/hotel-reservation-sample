package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/router"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	withLog    = false
	collection = "users"
)

type testDb struct {
	store.UserStore
}

func (tdb *testDb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() (*testDb, *mongo.Collection, *fiber.App, *Handler) {
	db := utils.NewDb()
	coll := db.Collection(collection)

	tdb := &testDb{
		UserStore: store.NewMongoUserStore(db),
	}

	app := router.New(withLog)
	handlers := NewHandler(tdb.UserStore)

	return tdb, coll, app, handlers
}

func TestPostUser(t *testing.T) {
	tdb, _, app, handlers := setup()
	defer tdb.tearDown(t)

	app.Post("/", handlers.HandlePostUser)

	t.Run("Validations", func(t *testing.T) {
		type test struct {
			expect   string
			input    types.CreateUserParams
			expected string
			status   int
		}

		partialInput := types.CreateUserParams{}
		invalidEmail := types.CreateUserParams{
			Email:     "invalid-email",
			FirstName: "Tan",
			LastName:  "Foo",
			Password:  "1234567",
		}
		invalidMinNames := types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "T",
			LastName:  "F",
			Password:  "1234567",
		}
		invalidMaxNames := types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT",
			LastName:  "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			Password:  "1234567",
		}
		invalidAlphaNames := types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "Test test",
			LastName:  "Foo foo",
			Password:  "1234567",
		}
		invalidPassword := types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "Test",
			LastName:  "Foo",
			Password:  "123456",
		}

		tests := []test{
			{
				expect:   "Should return all required fields error",
				input:    partialInput,
				status:   422,
				expected: `{"errors":{"Email":"required","FirstName":"required","LastName":"required","Password":"required"}}`,
			},
			{
				expect:   "Should return invalid email field error",
				input:    invalidEmail,
				status:   422,
				expected: `{"errors":{"Email":"email - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName minimum field error",
				input:    invalidMinNames,
				status:   422,
				expected: `{"errors":{"FirstName":"min - invalid","LastName":"min - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName maximum field error",
				input:    invalidMaxNames,
				status:   422,
				expected: `{"errors":{"FirstName":"max - invalid","LastName":"max - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName maximum field error",
				input:    invalidAlphaNames,
				status:   422,
				expected: `{"errors":{"FirstName":"alpha - invalid","LastName":"alpha - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName maximum field error",
				input:    invalidPassword,
				status:   422,
				expected: `{"errors":{"Password":"min - invalid"}}`,
			},
		}

		for _, tc := range tests {
			b, _ := json.Marshal(tc.input)
			req := utils.NewRequestWithHeader("POST", "/", bytes.NewReader(b))
			resp, err := app.Test(req)
			if err != nil {
				t.Error(err)
			}

			t.Run(fmt.Sprintf("should return %d status code", tc.status), func(t *testing.T) {
				if resp.StatusCode != tc.status {
					t.Errorf("expected status code %d but return %d", tc.status, resp.StatusCode)
				}
			})

			t.Run(tc.expect, func(t *testing.T) {
				body := make([]byte, resp.ContentLength)
				resp.Body.Read(body)

				if string(body) != tc.expected {
					t.Errorf("should return %s but received %s", tc.expected, string(body))
				}
			})
		}
	})

	t.Run("Insert user", func(t *testing.T) {
		params := types.CreateUserParams{
			Email:     "some@test.com",
			FirstName: "Tan",
			LastName:  "Foo",
			Password:  "1234567",
		}
		b, _ := json.Marshal(params)

		req := utils.NewRequestWithHeader("POST", "/", bytes.NewReader(b))
		res, err := app.Test(req)
		if err != nil {
			t.Error(err)
		}

		var createdUser types.User

		json.NewDecoder(res.Body).Decode(&createdUser)
		if len(createdUser.ID) == 0 {
			t.Errorf("expecting a user id to be set")
		}
		if len(createdUser.EncryptedPassword) > 0 {
			t.Errorf("should not include EncryptedPassword in json response")
		}
		if createdUser.FirstName != params.FirstName {
			t.Errorf("expected firstName %s but got %s", params.FirstName, createdUser.FirstName)
		}
		if createdUser.LastName != params.LastName {
			t.Errorf("expected lastName %s but got %s", params.LastName, createdUser.LastName)
		}
		if createdUser.Email != params.Email {
			t.Errorf("expected Email %s but got %s", params.Email, createdUser.Email)
		}
	})

	t.Run("Not insert user with existing email", func(t *testing.T) {
		params := types.CreateUserParams{
			Email:     "some@test.com",
			FirstName: "Test",
			LastName:  "Bar",
			Password:  "1234567",
		}
		b, _ := json.Marshal(params)

		req := utils.NewRequestWithHeader("POST", "/", bytes.NewReader(b))
		res, err := app.Test(req)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != 409 {
			t.Errorf("expected 409 conflict status but received %d", res.StatusCode)
		}
	})
}

func TestHandleGetUser(t *testing.T) {
	tdb, coll, app, handlers := setup()
	defer tdb.tearDown(t)

	newUsers := []interface{}{
		types.User{FirstName: "AA", LastName: "AA", Email: "aa@test.com", EncryptedPassword: "encrypted"},
		types.User{FirstName: "BB", LastName: "BB", Email: "bb@test.com", EncryptedPassword: "encrypted"},
	}

	fixtureUser, _ := newUsers[0].(types.User)

	result, err := coll.InsertMany(context.TODO(), newUsers)
	if err != nil {
		t.Error(err)
	}

	app.Get("/:id", handlers.HandleGetUser)

	t.Run("Validations userId is ObjectId", func(t *testing.T) {
		type test struct {
			id       string
			expect   string
			expected string
			status   int
		}

		tests := []test{
			{
				id:       "invalidId",
				expect:   "must return invalid ObjectId",
				status:   422,
				expected: `{"errors":{"ID":"id - invalid"}}`,
			},
		}

		for _, tc := range tests {
			req := utils.NewRequestWithHeader("GET", fmt.Sprintf("/%s", tc.id), nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Error(err)
			}

			t.Run(fmt.Sprintf("should return %d status code", tc.status), func(t *testing.T) {
				if resp.StatusCode != tc.status {
					t.Errorf("expected status code %d but return %d", tc.status, resp.StatusCode)
				}
			})

			t.Run(tc.expect, func(t *testing.T) {
				body := make([]byte, resp.ContentLength)
				resp.Body.Read(body)

				if string(body) != tc.expected {
					t.Errorf("should return %s but received %s", tc.expected, string(body))
				}
			})
		}
	})

	t.Run("get user by ID", func(t *testing.T) {
		objectId := result.InsertedIDs[0].(primitive.ObjectID)
		req := utils.NewRequestWithHeader("GET", fmt.Sprintf("/%s", objectId.Hex()), nil)
		res, err := app.Test(req)
		if err != nil {
			t.Error(err)
		}

		var user types.User

		json.NewDecoder(res.Body).Decode(&user)
		if user.ID.Hex() != objectId.Hex() {
			t.Errorf("expecting a user id %s received %s", objectId.Hex(), user.ID.Hex())
		}
		if len(user.EncryptedPassword) > 0 {
			t.Errorf("should not include EncryptedPassword in json response")
		}
		if user.FirstName != fixtureUser.FirstName {
			t.Errorf("expected firstName %s but got %s", "AA", user.FirstName)
		}
		if user.LastName != fixtureUser.LastName {
			t.Errorf("expected lastName %s but got %s", fixtureUser.LastName, user.LastName)
		}
		if user.Email != fixtureUser.Email {
			t.Errorf("expected Email %s but got %s", fixtureUser.Email, user.Email)
		}
	})
}
