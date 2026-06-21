package main

import (
	"crypto/rand"
	"math/big"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func generateAccessToken(hmacSecret []byte, userID string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:   "auth-service",
			Subject:  userID,
			Audience: []string{"user"},
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "auth_" + userID,
		},
	)

	return token.SignedString(hmacSecret)
}

func generateRefreshToken() (string, error) {
	token := make([]byte, 32)

	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		token[i] = alphabet[n.Int64()]
	}

	return string(token), nil
}
