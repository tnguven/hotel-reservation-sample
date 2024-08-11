package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
)

func insertTestUser(t *testing.T, userStore store.UserStore) *types.CreateUserParams {
	userParam := types.CreateUserParams{
		Email:     "Test",
		FirstName: "FirstName",
		LastName:  "LastName",
		Password:  "secret1234",
	}
	user, err := types.NewUserFromParams(userParam)
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return &userParam
}

func TestHandleAuthenticate(t *testing.T) {
	tdb, _, app, handlers := Setup()
	defer tdb.tearDown(t)

	insertedUser := insertTestUser(t, tdb.User)

	app.Post("/", handlers.HandleAuthenticate)

	t.Run("happy path", func(t *testing.T) {
		params := types.AuthParams{
			Email:    insertedUser.Email,
			Password: insertedUser.Password,
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  "/",
			Payload: bytes.NewReader(b),
		}
		req := testReq.NewRequestWithHeader()
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Fatal("expected http status of 200 but got %s", resp.StatusCode)
		}

		var result AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Error(err)
		}

	})

}
