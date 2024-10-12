package tests

import (
	"medods_test_task/handlers"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRoute(t *testing.T) {
	PATH := "/api/v1/user/add"
	w := httptest.NewRecorder()

	exampleUser := handlers.UserRequest{
		Email:    "",
		Password: "",
	}

	req := sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Unathorized", func(t *testing.T) {
		assert.Equal(t, 401, w.Code)
		assert.Equal(t, "No authorization header provided", w.Body.String())
	})

	// Generate JWT tokens

	tokenPair, err := MockTokenController.NewJWT("test@gmail.com", "127.0.0.1")

	if err != nil {
		t.Errorf(" [Error] Failed to initialize jwt token controller. %s", err)
	}

	w = httptest.NewRecorder()

	exampleUser = handlers.UserRequest{
		Email:    "",
		Password: "",
	}

	req = sendRequest(exampleUser, PATH)
	req.Header.Add("Authorization", tokenPair.AccessToken)
	router.ServeHTTP(w, req)

	t.Run("No user provided", func(t *testing.T) {
		assert.Equal(t, 400, w.Code)
	})

	w = httptest.NewRecorder()

	exampleUser = handlers.UserRequest{
		Email:    "test2@gmail.com",
		Password: "password",
	}

	req = sendRequest(exampleUser, PATH)
	req.Header.Add("Authorization", tokenPair.AccessToken)
	router.ServeHTTP(w, req)

	t.Run("Correct request", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}
