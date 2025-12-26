package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/snaztoz/watergun"
)

func TestProtectionForAdminRoutes(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()

	req.Header.Add("Authorization", "Bearer "+watergun.AdminAPIKey())

	adminRoutesAuth(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	).ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusForbidden {
		t.Errorf(
			"handler returned wrong status code: got %v",
			status,
		)
	}
}

func TestMissingBearerTokenForAdminRoutes(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()

	adminRoutesAuth(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf(
			"handler returned wrong status code: want %v, got %v",
			http.StatusForbidden,
			status,
		)
	}
}

func TestIncorrectBearerTokenForAdminRoutes(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin", nil)
	rr := httptest.NewRecorder()

	req.Header.Add("Authorization", "Bearer\t") // use tab instead of space

	adminRoutesAuth(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf(
			"handler returned wrong status code: want %v, got %v",
			http.StatusForbidden,
			status,
		)
	}
}
