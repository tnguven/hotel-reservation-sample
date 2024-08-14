package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	mid "github.com/tnguven/hotel-reservation-app/handler/middleware"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	withLog            = false
	invalidMaxCharName = "TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT"
)

func TestPostUser(t *testing.T) {
	_, app, handlers, validator := Setup(db, false)

	app.Post("/", mid.WithValidation(validator, InsertUserRequestSchema), handlers.HandlePostUser)

	t.Run("Validations", func(t *testing.T) {
		t.Parallel()
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
			FirstName: invalidMaxCharName,
			LastName:  invalidMaxCharName,
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
				status:   400,
				expected: `{"errors":{"Email":"required","FirstName":"required","LastName":"required","Password":"required"}}`,
			},
			{
				expect:   "Should return invalid email field error",
				input:    invalidEmail,
				status:   400,
				expected: `{"errors":{"Email":"email - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName and lastName minimum field error",
				input:    invalidMinNames,
				status:   400,
				expected: `{"errors":{"FirstName":"min - invalid","LastName":"min - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName and lastName maximum field error",
				input:    invalidMaxNames,
				status:   400,
				expected: `{"errors":{"FirstName":"max - invalid","LastName":"max - invalid"}}`,
			},
			{
				expect:   "Should return invalid firstName and lastName maximum field error",
				input:    invalidAlphaNames,
				status:   400,
				expected: `{"errors":{"FirstName":"alpha - invalid","LastName":"alpha - invalid"}}`,
			},
			{
				expect:   "Should return invalid password min field error",
				input:    invalidPassword,
				status:   400,
				expected: `{"errors":{"Password":"min - invalid"}}`,
			},
		}

		for _, tc := range tests {
			b, _ := json.Marshal(tc.input)
			testReq := utils.TestRequest{
				Method:  "POST",
				Target:  "/",
				Payload: bytes.NewReader(b),
			}
			resp, err := app.Test(testReq.NewRequestWithHeader())
			if err != nil {
				t.Error(err)
			}

			t.Run(fmt.Sprintf("should return %d status code", tc.status), func(t *testing.T) {
				t.Parallel()
				if resp.StatusCode != tc.status {
					t.Errorf("expected status code %d but return %d", tc.status, resp.StatusCode)
				}
			})

			t.Run(tc.expect, func(t *testing.T) {
				t.Parallel()
				body := make([]byte, resp.ContentLength)
				resp.Body.Read(body)

				if string(body) != tc.expected {
					t.Errorf("should return %s but received %s", tc.expected, string(body))
				}
			})
		}
	})

	email := "insert_user@test.com"

	t.Run("Insert user", func(t *testing.T) {
		params := types.CreateUserParams{
			Email:     email,
			FirstName: "Tan",
			LastName:  "Foo",
			Password:  "1234567",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  "/",
			Payload: bytes.NewReader(b),
		}
		res, err := app.Test(testReq.NewRequestWithHeader())
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
			Email:     email,
			FirstName: "same",
			LastName:  "email",
			Password:  "1234567",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  "/",
			Payload: bytes.NewReader(b),
		}
		res, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != 409 {
			t.Errorf("expected 409 conflict status but received %d", res.StatusCode)
		}
	})
}

func TestHandleGetUser(t *testing.T) {
	tdb, app, handlers, _ := Setup(db, false)

	var (
		firstName = "get"
		lastName  = "userbyid"
	)

	app.Get("/:id", handlers.HandleGetUser)

	t.Run("get user by ID", func(t *testing.T) {
		user := fixtures.AddUser(*tdb.Store, firstName, lastName, false)
		testReq := utils.TestRequest{
			Method: "GET",
			Target: fmt.Sprintf("/%s", user.ID.Hex()),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		var fetchedUser types.User
		json.NewDecoder(resp.Body).Decode(&fetchedUser)

		if len(fetchedUser.EncryptedPassword) > 0 {
			t.Errorf("should not include EncryptedPassword in json response")
		}
		if fetchedUser.IsAdmin != false {
			t.Errorf("should create isAdmin false but received isAdmin true")
		}
		if fetchedUser.FirstName != firstName {
			t.Errorf("expected firstName %s but got %s", "aa", fetchedUser.FirstName)
		}
		if fetchedUser.LastName != lastName {
			t.Errorf("expected lastName %s but got %s", "bb", fetchedUser.LastName)
		}
		if fetchedUser.Email != "get_userbyid@test.com" {
			t.Errorf("expected Email %s but got %s", "get_userbyid@test.com", fetchedUser.Email)
		}
	})
}

func TestHandlePutUser(t *testing.T) {
	tdb, app, handlers, validator := Setup(db, false)

	app.Put("/:id", mid.WithValidation(validator, UpdateUserRequestSchema), handlers.HandlePutUser)

	t.Run("Validations", func(t *testing.T) {
		t.Parallel()
		invalidMinFields := types.UpdateUserParams{
			FirstName: "T",
			LastName:  "F",
		}
		invalidMaxFields := types.UpdateUserParams{
			FirstName: invalidMaxCharName,
			LastName:  invalidMaxCharName,
		}
		invalidAlphaFields := types.UpdateUserParams{
			FirstName: "Test1",
			LastName:  "Test2",
		}
		validParams := types.UpdateUserParams{
			FirstName: "Foo",
			LastName:  "Bar",
		}

		type test struct {
			id       string
			input    types.UpdateUserParams
			expect   string
			expected string
			status   int
		}

		tests := []test{
			{
				id:       primitive.NewObjectID().Hex(),
				expect:   "Should return invalid firstName and lastName minimum field error",
				input:    invalidMinFields,
				status:   400,
				expected: `{"errors":{"FirstName":"min - invalid","LastName":"min - invalid"}}`,
			},
			{
				id:       primitive.NewObjectID().Hex(),
				expect:   "Should return invalid firstName and lastName maximum field error",
				input:    invalidMaxFields,
				status:   400,
				expected: `{"errors":{"FirstName":"max - invalid","LastName":"max - invalid"}}`,
			},
			{
				id:       primitive.NewObjectID().Hex(),
				expect:   "Should return invalid firstName and lastName maximum field error",
				input:    invalidAlphaFields,
				status:   400,
				expected: `{"errors":{"FirstName":"alpha - invalid","LastName":"alpha - invalid"}}`,
			},
			{
				id:       "invalidId",
				input:    validParams,
				expect:   "must return required and invalid fields",
				status:   400,
				expected: `{"errors":{"ID":"id - invalid"}}`,
			},
		}

		for _, tc := range tests {
			b, _ := json.Marshal(tc.input)
			testReq := utils.TestRequest{
				Method:  "PUT",
				Target:  fmt.Sprintf("/%s", tc.id),
				Payload: bytes.NewReader(b),
			}
			resp, err := app.Test(testReq.NewRequestWithHeader())
			if err != nil {
				t.Fatal(err)
			}

			t.Run(fmt.Sprintf("should return %d status code", tc.status), func(t *testing.T) {
				t.Parallel()
				if resp.StatusCode != tc.status {
					t.Fatalf("expected status code %d but return %d", tc.status, resp.StatusCode)
				}
			})

			t.Run(tc.expect, func(t *testing.T) {
				t.Parallel()
				body := make([]byte, resp.ContentLength)
				resp.Body.Read(body)

				if string(body) != tc.expected {
					t.Errorf("should return %s but received %s", tc.expected, string(body))
				}
			})
		}
	})

	t.Run("update user", func(t *testing.T) {
		user := fixtures.AddUser(*tdb.Store, "update", "user", false)
		params := types.UpdateUserParams{
			FirstName: "user",
			LastName:  "update",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "PUT",
			Target:  fmt.Sprintf("/%s", user.ID.Hex()),
			Payload: bytes.NewReader(b),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["updated"] != user.ID.Hex() {
			t.Errorf("expecting user id %s to be updated, but got %v", user.ID.Hex(), response["updated"])
		}
	})

	t.Run("return error if the id does not exist", func(t *testing.T) {
		params := types.UpdateUserParams{
			FirstName: "AAB",
			LastName:  "Bar",
		}
		obi := primitive.NewObjectID().Hex()
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "PUT",
			Target:  fmt.Sprintf("/%s", obi),
			Payload: bytes.NewReader(b),
		}
		res, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		expectedError := fmt.Sprintf("no user found with id %s", obi)

		if response["error"] != expectedError {
			t.Errorf("expecting error %s but received %s", expectedError, response["error"])
		}
	})
}
