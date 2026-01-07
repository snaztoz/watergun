package socket

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	privateKey := generateKey()

	t.Run("normal case", func(t *testing.T) {
		token := createToken(
			jwt.MapClaims{
				"sub": "some-user-id",
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Minute).Unix(),
			},
			privateKey,
		)

		subject, err := Authenticate(token, privateKey.Public())

		assert.Nil(t, err)
		assert.Equal(t, "some-user-id", subject)
	})

	t.Run("with no `exp` claim", func(t *testing.T) {
		token := createToken(
			jwt.MapClaims{
				"sub": "some-user-id",
				"iat": time.Now().Unix(),
			},
			privateKey,
		)

		_, err := Authenticate(token, privateKey.Public())

		assert.NotNil(t, err)
	})

	t.Run("with no `sub` claim", func(t *testing.T) {
		token := createToken(
			jwt.MapClaims{
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Minute).Unix(),
			},
			privateKey,
		)

		sub, _ := Authenticate(token, privateKey.Public())

		assert.Empty(t, sub)
	})

	t.Run("with expired token", func(t *testing.T) {
		token := createToken(
			jwt.MapClaims{
				"sub": "some-user-id",
				"iat": time.Now().Add(-1 * time.Minute).Unix(),
				"exp": time.Now().Add(-30 * time.Second).Unix(),
			},
			privateKey,
		)

		_, err := Authenticate(token, privateKey.Public())

		assert.NotNil(t, err)
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
