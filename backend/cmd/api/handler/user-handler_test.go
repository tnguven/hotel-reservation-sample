package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tnguven/hotel-reservation-app/cmd/api/handler"
	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPostUser(t *testing.T) {
	tdb, app := handler.Setup(db, false, configs)
	invalidMaxCharName := strings.Repeat("a", 49)
	const target = "/v1/users"

	t.Run("Validations", func(t *testing.T) {
		t.Parallel()
		type test struct {
			desc     string
			input    *types.CreateUserParams
			expected *utils.Error
			status   int
		}

		partialInput := &types.CreateUserParams{}
		invalidEmail := &types.CreateUserParams{
			Email:     "invalid-email",
			FirstName: "Tan",
			LastName:  "Foo",
			Password:  "1234567",
		}
		invalidMinNames := &types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "T",
			LastName:  "F",
			Password:  "1234567",
		}
		invalidMaxNames := &types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: invalidMaxCharName,
			LastName:  invalidMaxCharName,
			Password:  "1234567",
		}
		invalidAlphaNames := &types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "Test test",
			LastName:  "Foo foo",
			Password:  "1234567",
		}
		invalidPassword := &types.CreateUserParams{
			Email:     "test@test.com",
			FirstName: "Test",
			LastName:  "Foo",
			Password:  "123456",
		}

		tests := []test{
			{
				desc:   "Should return all required fields error",
				input:  partialInput,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"Email":     "required",
							"FirstName": "required",
							"LastName":  "required",
							"Password":  "required",
						},
					},
				},
			},
			{
				desc:   "Should return invalid email field error",
				input:  invalidEmail,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"Email": "email - invalid",
						},
					},
				},
			},
			{
				desc:   "Should return invalid firstName and lastName minimum field error",
				input:  invalidMinNames,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"FirstName": "min - invalid",
							"LastName":  "min - invalid",
						},
					},
				},
			},
			{
				desc:   "Should return invalid firstName and lastName maximum field error",
				input:  invalidMaxNames,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"FirstName": "max - invalid",
							"LastName":  "max - invalid",
						},
					},
				},
			},
			{
				desc:   "Should return invalid firstName and lastName maximum field error",
				input:  invalidAlphaNames,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"FirstName": "alpha - invalid",
							"LastName":  "alpha - invalid",
						},
					},
				},
			},
			{
				desc:   "Should return invalid password min field error",
				input:  invalidPassword,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"Password": "min - invalid",
						},
					},
				},
			},
		}

		for _, tc := range tests {
			b, _ := json.Marshal(tc.input)
			testReq := utils.TestRequest{
				Method:  "POST",
				Target:  target,
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

			t.Run(tc.desc, func(t *testing.T) {
				t.Parallel()
				var body utils.Error
				if errDecode := json.NewDecoder(resp.Body).Decode(&body); errDecode != nil {
					t.Fatal(errDecode)
				}

				if !reflect.DeepEqual(&body, tc.expected) {
					t.Errorf("expected response %+v, but got %+v", tc.expected, body)
				}
			})
		}
	})

	t.Run("Insert user", func(t *testing.T) {
		params := types.CreateUserParams{
			Email:     "insert_user@test.com",
			FirstName: "Tan",
			LastName:  "Foo",
			Password:  "1234567",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  target,
			Payload: bytes.NewReader(b),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		var result utils.GenericResponse
		if errDecode := json.NewDecoder(resp.Body).Decode(&result); errDecode != nil {
			t.Fatal(errDecode)
		}

		data, err := json.Marshal(result.Data)
		if err != nil {
			t.Fatalf("Failed to marshal Data field: %v", err)
		}
		var createdUser *types.User
		if err = json.Unmarshal(data, &createdUser); err != nil {
			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
		}

		if len(createdUser.ID.Hex()) == 0 {
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
		fixtures.AddUser(*tdb.Store, "insertnew", "user", false)
		params := types.CreateUserParams{
			Email:     "insertnew_user@test.com",
			FirstName: "same",
			LastName:  "email",
			Password:  "1234567",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  target,
			Payload: bytes.NewReader(b),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != 409 {
			t.Errorf("expected 409 conflict status but received %d", resp.StatusCode)
		}
	})
}

func TestHandleGetUser(t *testing.T) {
	tdb, app := handler.Setup(db, false, configs)
	var (
		firstName = "get"
		lastName  = "userbyid"
	)
	const target = "/v1/users"

	t.Run("get user by ID", func(t *testing.T) {
		user := fixtures.AddUser(*tdb.Store, firstName, lastName, false)
		testReq := utils.TestRequest{
			Method: "GET",
			Target: fmt.Sprintf("%s/%s", target, user.ID.Hex()),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		var result utils.GenericResponse
		if errDecode := json.NewDecoder(resp.Body).Decode(&result); errDecode != nil {
			t.Fatal(errDecode)
		}

		data, err := json.Marshal(result.Data)
		if err != nil {
			t.Fatalf("Failed to marshal Data field: %v", err)
		}
		var fetchedUser *types.User
		if err = json.Unmarshal(data, &fetchedUser); err != nil {
			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
		}

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
	tdb, app := handler.Setup(db, false, configs)
	invalidMaxCharName := strings.Repeat("a", 49)
	const target = "/v1/users"

	t.Run("Validations", func(t *testing.T) {
		t.Parallel()
		invalidMinFields := &types.UpdateUserParams{
			FirstName: "T",
			LastName:  "F",
		}
		invalidMaxFields := &types.UpdateUserParams{
			FirstName: invalidMaxCharName,
			LastName:  invalidMaxCharName,
		}
		invalidAlphaFields := &types.UpdateUserParams{
			FirstName: "Test1",
			LastName:  "Test2",
		}
		validParams := &types.UpdateUserParams{
			FirstName: "Foo",
			LastName:  "Bar",
		}

		type test struct {
			id       string
			input    *types.UpdateUserParams
			desc     string
			expected *utils.Error
			status   int
		}

		tests := []test{
			{
				id:     primitive.NewObjectID().Hex(),
				desc:   "Should return invalid firstName and lastName minimum field error",
				input:  invalidMinFields,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"FirstName": "min - invalid",
							"LastName":  "min - invalid",
						},
					},
				},
			},
			{
				id:     primitive.NewObjectID().Hex(),
				desc:   "Should return invalid firstName and lastName maximum field error",
				input:  invalidMaxFields,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"FirstName": "max - invalid",
							"LastName":  "max - invalid",
						},
					},
				},
			},
			{
				id:     primitive.NewObjectID().Hex(),
				desc:   "Should return invalid firstName and lastName maximum field error",
				input:  invalidAlphaFields,
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"FirstName": "alpha - invalid",
							"LastName":  "alpha - invalid",
						},
					},
				},
			},
			{
				id:     "invalidId",
				input:  validParams,
				desc:   "must return required and invalid fields",
				status: 400,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"ID": "id - invalid",
						},
					},
				},
			},
		}

		for _, tc := range tests {
			b, _ := json.Marshal(tc.input)
			testReq := utils.TestRequest{
				Method:  "PUT",
				Target:  fmt.Sprintf("%s/%s", target, tc.id),
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

			t.Run(tc.desc, func(t *testing.T) {
				t.Parallel()
				var body utils.Error
				if errDecode := json.NewDecoder(resp.Body).Decode(&body); errDecode != nil {
					t.Fatal(errDecode)
				}

				if !reflect.DeepEqual(&body, tc.expected) {
					t.Errorf("expected response %+v, but got %+v", tc.expected, body)
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
			Target:  fmt.Sprintf("%s/%s", target, user.ID.Hex()),
			Payload: bytes.NewReader(b),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Fatal(err)
		}

		var result utils.GenericResponse
		if errDecode := json.NewDecoder(resp.Body).Decode(&result); errDecode != nil {
			t.Fatal(errDecode)
		}

		if result.Msg != fmt.Sprintf("User %s updated", user.ID.Hex()) {
			t.Errorf("expecting user id %s to be updated, but got %v", user.ID.Hex(), result.Msg)
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
			Target:  fmt.Sprintf("%s/%s", target, obi),
			Payload: bytes.NewReader(b),
		}
		resp, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		var response utils.GenericResponse
		if errDecode := json.NewDecoder(resp.Body).Decode(&response); errDecode != nil {
			t.Fatal(errDecode)
		}

		expectedError := fmt.Sprintf("no user found with id %s", obi)

		if response.Msg != expectedError {
			t.Errorf("expecting error %s but received %s", expectedError, response.Msg)
		}
	})
}
