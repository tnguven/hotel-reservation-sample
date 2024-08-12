package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/store"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
)

func insertTestUser(t *testing.T, userStore store.UserStore) *types.User {
	userParam := types.CreateUserParams{
		Email:     "test@test.com",
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

	return user
}

func TestHandleAuthenticate(t *testing.T) {
	t.Parallel()

	tdb, _, app, handlers, _ := Setup()
	defer tdb.tearDown(t)

	insertedUser := insertTestUser(t, tdb.User)

	app.Post("/", handlers.HandleAuthenticate)

	params := types.AuthParams{
		Email:    "test@test.com",
		Password: "secret1234",
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
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response")
	}
	// set the encrypted password to an empty string, because we do not return that in any
	// JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, result.User) {
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestHandleAuthenticateFailure(t *testing.T) {
	t.Parallel()

	tdb, _, app, handlers, _ := Setup()
	defer tdb.tearDown(t)

	app.Post("/", handlers.HandleAuthenticate)

	params := types.AuthParams{
		Email:    "test@test.com",
		Password: "wrong",
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

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var result genericResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if result.Type != "error" {
		t.Fatalf("expected to get type error but received: %s", result.Type)
	}

	if result.Msg != "invalid credential" {
		t.Fatalf("expected to get msg invalid credential but received: %s", result.Msg)
	}

}
