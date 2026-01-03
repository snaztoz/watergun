package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/snaztoz/watergun"
	"github.com/stretchr/testify/assert"
)

func TestProtectionForAdminRoutes(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()

	req.Header.Add("Authorization", "Bearer "+watergun.AdminAPIKey())

	adminRoutesAuth(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	).ServeHTTP(rr, req)

	assert.NotEqual(t, 403, rr.Code)
}

func TestMissingBearerTokenForAdminRoutes(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()

	adminRoutesAuth(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	).ServeHTTP(rr, req)

	assert.Equal(t, 403, rr.Code)
}

func TestIncorrectBearerTokenForAdminRoutes(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()

	req.Header.Add("Authorization", "Bearer\t") // use tab instead of space

	adminRoutesAuth(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	).ServeHTTP(rr, req)

	assert.Equal(t, 403, rr.Code)
}
