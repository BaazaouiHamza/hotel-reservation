package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/baazaouihamza/hotel-reservation/db"
	"github.com/baazaouihamza/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "hamza",
		LastName:  "baazaoui",
		Email:     "hamza@baazaoui.com",
		Password:  "1234567",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "hamza@baazaoui.com",
		Password: "1234567",
	}

	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("exprected http status 200 but god %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected JWT token to be present in the auth response")
	}

	// set the encrypted password to an empty string, because we omit this field when we return it in any
	// JSON response
	insertedUser.EncryptedPasswod = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatalf("expected user to be the inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "hamza@baazaoui.com",
		Password: "1234567notcorrect",
	}

	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("exprected http status 400 but god %d", resp.StatusCode)
	}

	var genResp genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be valid <invalid credentials> but got %s", genResp.Msg)
	}
}
