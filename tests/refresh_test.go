package tests

import (
	"context"
	"medods_test_task/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshRoute(t *testing.T) {
	PATH := "/api/v1/auth/refresh"
	w := httptest.NewRecorder()

	// Generate refresh token

	refreshToken, err := MockTokenController.NewRefreshToken()

	if err != nil {
		t.Error("Failed to generate refresh token")
	}

	newRefreshTokenHash, err := utils.HashBCrypt([]byte(refreshToken))

	if err != nil {
		t.Error("Failed to hash refresh token")
	}

	// Add correct owner

	err = MockDB.AddRefreshToken(
		context.Background(),
		"test@gmail.com",
		string(newRefreshTokenHash),
		MockTokenController.RefreshTokenTTL,
		"127.0.0.1",
	)

	if err != nil {
		t.Error("Failed to add refresh token to database")
	}

	req, _ := http.NewRequest(
		"POST",
		PATH,
		nil,
	)
	q := req.URL.Query()
	q.Add("refresh_token", refreshToken)
	// Test for invalid owner
	q.Add("email", "test1@gmail.com")
	req.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, req)

	t.Run("Invalid owner", func(t *testing.T) {
		assert.Equal(t, 404, w.Code)
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		"POST",
		PATH,
		nil,
	)
	router.ServeHTTP(w, req)

	t.Run("No parameters", func(t *testing.T) {
		assert.Equal(t, 400, w.Code)
	})

	// Generate invalid token

	invalidToken, err := MockTokenController.NewRefreshToken()

	if err != nil {
		t.Error("Failed to hash refresh token")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		"POST",
		PATH,
		nil,
	)
	q = req.URL.Query()
	q.Add("refresh_token", invalidToken)
	q.Add("email", "test@gmail.com")
	req.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, req)

	t.Run("Invalid token", func(t *testing.T) {
		assert.Equal(t, 401, w.Code)
	})

	// Correct token and owner

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		"POST",
		PATH,
		nil,
	)
	q = req.URL.Query()
	q.Add("refresh_token", refreshToken)
	q.Add("email", "test@gmail.com")
	req.URL.RawQuery = q.Encode()
	router.ServeHTTP(w, req)

	t.Run("Correct token", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}
