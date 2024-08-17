package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/cmd/api/handler"
	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"github.com/tnguven/hotel-reservation-app/internals/utils"
)

func TestHandleAuthenticate(t *testing.T) {
	tdb, app := handler.Setup(db, false, configs)
	const target = "/v1/auth"

	t.Run("Validations", func(t *testing.T) {
		t.Parallel()
		invalidEmail := &types.AuthParams{
			Email:    "foo_bar@invalid",
			Password: "foo_bar",
		}
		invalidPassword := &types.AuthParams{
			Email:    "foo_bar@valid.com",
			Password: "foobar",
		}
		invalidBoth := &types.AuthParams{
			Email:    "foo_bar@invalid",
			Password: "foobar",
		}

		type test struct {
			desc     string
			input    *types.AuthParams
			expected *utils.Error
			status   int
		}

		inputs := []test{
			{
				desc:  "should return invalid email error",
				input: invalidEmail,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"Email": "email - invalid",
						},
					},
				},
				status: 400,
			},
			{
				desc:  "should return invalid password error",
				input: invalidPassword,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"Password": "min - invalid",
						},
					},
				},
				status: 400,
			},
			{
				desc:  "should return invalid password error",
				input: invalidBoth,
				expected: &utils.Error{
					GenericResponse: &utils.GenericResponse{
						Status: 400,
						Msg:    "Bad Request",
						Errors: map[string]interface{}{
							"Email":    "email - invalid",
							"Password": "min - invalid",
						},
					},
				},
				status: 400,
			},
		}

		for _, test := range inputs {
			b, _ := json.Marshal(test.input)
			testReq := utils.TestRequest{
				Method:  "POST",
				Target:  target,
				Payload: bytes.NewReader(b),
			}
			req := testReq.NewRequestWithHeader()
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			t.Run(fmt.Sprintf("should return %d status code", test.status), func(t *testing.T) {
				t.Parallel()
				if resp.StatusCode != test.status {
					t.Errorf("expected status code %d but return %d", test.status, resp.StatusCode)
				}
			})

			t.Run(test.desc, func(t *testing.T) {
				t.Parallel()
				var body utils.Error
				if errDecode := json.NewDecoder(resp.Body).Decode(&body); errDecode != nil {
					t.Fatal(errDecode)
				}

				if !reflect.DeepEqual(&body, test.expected) {
					t.Errorf("expected response %+v, but got %+v", test.expected, body)
				}
			})
		}
	})

	t.Run("success_with_correct_password", func(t *testing.T) {
		t.Parallel()
		user := fixtures.AddUser(*tdb.Store, "auth", "success", false)
		params := types.AuthParams{
			Email:    "auth_success@test.com",
			Password: "auth_success",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  target,
			Payload: bytes.NewReader(b),
		}
		req := testReq.NewRequestWithHeader()
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
		}

		var result utils.GenericResponse
		if errDecode := json.NewDecoder(resp.Body).Decode(&result); errDecode != nil {
			t.Fatal(errDecode)
		}

		authResponseData, err := json.Marshal(result.Data)
		if err != nil {
			t.Fatalf("Failed to marshal Data field: %v", err)
		}

		var authResponse handler.AuthResponse
		if err = json.Unmarshal(authResponseData, &authResponse); err != nil {
			t.Fatalf("Failed to unmarshal Data field into AuthResponse: %v", err)
		}
		if authResponse.Token == "" {
			t.Fatal("expected the JWT token to be present in the auth response")
		}
		if authResponse.User.Email != "auth_success@test.com" {
			t.Fatalf("Expected user email 'auth_success@test.com' but got %s", authResponse.User.Email)
		}

		// set the encrypted password to an empty string, because we do not return that in any
		user.EncryptedPassword = ""
		if !reflect.DeepEqual(user, authResponse.User) {
			t.Fatal("expected the user to be the inserted user")
		}
	})

	t.Run("failure_with_wrong_password", func(t *testing.T) {
		t.Parallel()
		fixtures.AddUser(*tdb.Store, "unAuth", "user", false)
		params := types.AuthParams{
			Email:    "unAuth_user@test.com",
			Password: "wrong-valid-password",
		}
		b, _ := json.Marshal(params)
		testReq := utils.TestRequest{
			Method:  "POST",
			Target:  target,
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

		var result utils.Error
		if decodeErr := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(decodeErr)
		}
		if result.Status != fiber.StatusBadRequest {
			t.Fatalf("expected to get status code %d but received: %d", fiber.StatusBadRequest, result.Status)
		}
		if result.Msg != "invalid credentials" {
			t.Fatalf("expected to get msg invalid credential but received: %s", result.Msg)
		}
	})
}

func TestHandleSignin(t *testing.T) {
	tdb, app := handler.Setup(db, false, configs)
	const target = "/v1/auth/signin"
	invalidMaxCharName := strings.Repeat("a", 49)

	t.Run("validations", func(t *testing.T) {
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
			resp, reqErr := app.Test(testReq.NewRequestWithHeader())
			if reqErr != nil {
				t.Fatal(reqErr)
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

	t.Run("Signin user", func(t *testing.T) {
		t.Parallel()
		params := types.CreateUserParams{
			Email:     "sing_in@test.com",
			FirstName: "sign",
			LastName:  "in",
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
		t.Parallel()
		fixtures.AddUser(*tdb.Store, "signnew", "user", false)
		params := types.CreateUserParams{
			Email:     "signnew_user@test.com",
			FirstName: "Test",
			LastName:  "Bar",
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
