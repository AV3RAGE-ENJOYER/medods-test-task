package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/user/ping", nil)
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
	req, _ = http.NewRequest("GET", "/api/v1/user/ping", nil)
	req.Header.Add("Authorization", tokenPair.AccessToken)
	router.ServeHTTP(w, req)

	t.Run("StatusCode", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})

	t.Run("ResponseBody", func(t *testing.T) {
		assert.Equal(t, "pong", w.Body.String())
	})
}
