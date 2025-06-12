package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Uses    int `json:"uses"`
	MaxUses int `json:"max_uses"`
	jwt.RegisteredClaims
}

const (
	tokenExpiration = 1 * time.Hour
	maxTokenUses    = 5
	secretKey       = "your-256-bit-secret"
)

func GenerateToken() (string, error) {
	claims := &Claims{
		Uses:    0,
		MaxUses: maxTokenUses,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (bool, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return false, nil, err
	}

	if !token.Valid {
		return false, nil, nil
	}

	if claims.Uses >= claims.MaxUses {
		return false, claims, nil
	}

	claims.Uses++
	return true, claims, nil
}
