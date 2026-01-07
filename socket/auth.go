package socket

import (
	"crypto"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(token string, key crypto.PublicKey) (string, error) {
	parsed, err := jwt.Parse(
		token,
		func(t *jwt.Token) (any, error) {
			return key, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return "", err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Failed to parse JWT claims")
	}

	return claims.GetSubject()
}
