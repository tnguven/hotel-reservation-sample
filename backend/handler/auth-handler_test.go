package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tnguven/hotel-reservation-app/db/fixtures"
	"github.com/tnguven/hotel-reservation-app/types"
	"github.com/tnguven/hotel-reservation-app/utils"
)

func TestHandleAuthenticate(t *testing.T) {
	tdb, app := setup(db, false, configs)
	const target = "/v1/auth"

	t.Run("Validations", func(t *testing.T) {
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
			expect   string
			input    *types.AuthParams
			expected string
			status   int
		}
		inputs := []test{
			{
				expect:   "should return invalid email error",
				input:    invalidEmail,
				expected: `{"errors":{"Email":"email - invalid"}}`,
				status:   400,
			},
			{
				expect:   "should return invalid password error",
				input:    invalidPassword,
				expected: `{"errors":{"Password":"min - invalid"}}`,
				status:   400,
			},
			{
				expect:   "should return invalid password error",
				input:    invalidBoth,
				expected: `{"errors":{"Email":"email - invalid","Password":"min - invalid"}}`,
				status:   400,
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
				if resp.StatusCode != test.status {
					t.Errorf("expected status code %d but return %d", test.status, resp.StatusCode)
				}
			})

			t.Run(test.expect, func(t *testing.T) {
				body := make([]byte, resp.ContentLength)
				resp.Body.Read(body)

				if string(body) != test.expected {
					t.Errorf("should return %s but received %s", test.expected, string(body))
				}
			})
		}
	})

	t.Run("success_with_correct_password", func(t *testing.T) {
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
		user.EncryptedPassword = ""
		if !reflect.DeepEqual(user, result.User) {
			t.Fatal("expected the user to be the inserted user")
		}
	})

	t.Run("failure_with_wrong_password", func(t *testing.T) {
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
	})
}

func TestHandleSignin(t *testing.T) {
	_, app := setup(db, false, configs)
	const target = "/v1/auth/signin"

	t.Run("validate singin inputs", func(t *testing.T) {
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
	email := "sing_in@test.com"
	t.Run("Signin user", func(t *testing.T) {
		params := types.CreateUserParams{
			Email:     email,
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
		res, err := app.Test(testReq.NewRequestWithHeader())
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != 409 {
			t.Errorf("expected 409 conflict status but received %d", res.StatusCode)
		}
	})
}
