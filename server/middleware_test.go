package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/snaztoz/watergun/serverctx"
)

func TestAccessTokenParser(t *testing.T) {
	t.Run("with bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		rr := httptest.NewRecorder()

		key := "some-key-here"

		req.Header.Add("Authorization", "Bearer "+key)

		accessTokenParser(allowQueryParamToken)(
			http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				assert.Equal(t, key, r.Context().Value(serverctx.AccessTokenKey))
			}),
		).ServeHTTP(rr, req)
	})

	t.Run("with query param token", func(t *testing.T) {
		key := "some-key-here"

		req := httptest.NewRequest("GET", "/socket?token="+key, nil)
		rr := httptest.NewRecorder()

		accessTokenParser(allowQueryParamToken)(
			http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				assert.Equal(t, key, r.Context().Value(serverctx.AccessTokenKey))
			}),
		).ServeHTTP(rr, req)
	})

	t.Run("with no bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		rr := httptest.NewRecorder()

		accessTokenParser(allowQueryParamToken)(
			http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
		).ServeHTTP(rr, req)

		assert.NotEqual(t, 403, rr.Result().StatusCode)
	})

	t.Run("with unpermitted query param token", func(t *testing.T) {
		key := "some-key-here"

		req := httptest.NewRequest("GET", "/socket?token="+key, nil)
		rr := httptest.NewRecorder()

		accessTokenParser(!allowQueryParamToken)(
			http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {}),
		).ServeHTTP(rr, req)

		assert.NotEqual(t, 403, rr.Result().StatusCode)
	})
}

func TestSocketRouteAuth(t *testing.T) {
	privateKey := generateKey()

	t.Run("with bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/socket", nil)
		rr := httptest.NewRecorder()

		token := createToken(
			jwt.MapClaims{
				"sub": "some-user-id",
				"exp": time.Now().Add(time.Minute).Unix(),
			},
			privateKey,
		)

		req.Header.Add("Authorization", "Bearer "+token)

		accessTokenParser(allowQueryParamToken)(
			socketRouteAuth(privateKey.Public())(
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "some-user-id", r.Context().Value(serverctx.UserIDKey))
				}),
			),
		).ServeHTTP(rr, req)

		assert.NotEqual(t, 403, rr.Result().StatusCode)
	})

	t.Run("with no bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/socket", nil)
		rr := httptest.NewRecorder()

		accessTokenParser(allowQueryParamToken)(
			socketRouteAuth(privateKey.Public())(
				http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
			),
		).ServeHTTP(rr, req)

		assert.Equal(t, 403, rr.Result().StatusCode)
	})

	t.Run("with expired bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/socket", nil)
		rr := httptest.NewRecorder()

		token := createToken(
			jwt.MapClaims{
				"sub": "some-user-id",
				"exp": time.Now().Add(-1 * time.Minute).Unix(),
			},
			privateKey,
		)

		req.Header.Add("Authorization", "Bearer "+token)

		accessTokenParser(allowQueryParamToken)(
			socketRouteAuth(privateKey.Public())(
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {}),
			),
		).ServeHTTP(rr, req)

		assert.Equal(t, 403, rr.Result().StatusCode)
	})
}

func TestAdminRoutesAuth(t *testing.T) {
	t.Run("with bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		rr := httptest.NewRecorder()

		req.Header.Add("Authorization", "Bearer "+adminAPIKey())

		accessTokenParser(allowQueryParamToken)(
			adminRoutesAuth(
				http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
			),
		).ServeHTTP(rr, req)

		assert.NotEqual(t, 403, rr.Result().StatusCode)
	})

	t.Run("with no bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		rr := httptest.NewRecorder()

		accessTokenParser(allowQueryParamToken)(
			adminRoutesAuth(
				http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
			),
		).ServeHTTP(rr, req)

		assert.Equal(t, 403, rr.Result().StatusCode)
	})

	t.Run("with invalid bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin", nil)
		rr := httptest.NewRecorder()

		req.Header.Add("Authorization", "Bearer\t") // use tab instead of space

		accessTokenParser(allowQueryParamToken)(
			adminRoutesAuth(
				http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
			),
		).ServeHTTP(rr, req)

		assert.Equal(t, 403, rr.Result().StatusCode)
	})
}

func createToken(claims jwt.Claims, privateKey *ecdsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func generateKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return key
}
