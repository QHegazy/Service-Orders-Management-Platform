package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing JWTs (use env variable in production)
var jwtKey = []byte(getEnv("JWT_SECRET", "my_secret_key"))

// Claims struct (you can add more fields if needed)
type EntityData struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Belong   string `json:"belong,omitempty"`
	Role     string `json:"role,omitempty"`
}

type Claims struct {
	Data EntityData

	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token
func GenerateToken(data EntityData, duration time.Duration, subject string) (string, error) {
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "backend",
			Subject:   subject, // "access" or "refresh"
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken validaates a token and returns the claims
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Small helper to read environment with default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
