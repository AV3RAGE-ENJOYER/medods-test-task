package tests

import (
	"medods_test_task/handlers"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginRoute(t *testing.T) {
	PATH := "/api/v1/auth/login"
	w := httptest.NewRecorder()

	exampleUser := handlers.UserRequest{
		Email:    "test1@gmail.com",
		Password: "test",
	}

	req := sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Non existing user", func(t *testing.T) {
		assert.Equal(t, 404, w.Code)
		t.Log(w.Body.String())
	})

	w = httptest.NewRecorder()
	exampleUser = handlers.UserRequest{
		Email:    "test@gmail.com",
		Password: "test1",
	}

	req = sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Incorrect password", func(t *testing.T) {
		assert.Equal(t, 401, w.Code)
	})

	w = httptest.NewRecorder()
	exampleUser = handlers.UserRequest{
		Email:    "test@gmail.com",
		Password: "test",
	}

	req = sendRequest(exampleUser, PATH)
	router.ServeHTTP(w, req)

	t.Run("Existing user", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}
