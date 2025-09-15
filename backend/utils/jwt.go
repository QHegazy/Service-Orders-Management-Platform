package utils

import (
	"backend/internal/redis"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	goredis "github.com/redis/go-redis/v9"
)

var jwtKey = []byte(getEnv("JWT_SECRET", "my_secret_key"))

type EntityData struct {
	ID       string   `json:"id"`
	Username string   `json:"username,omitempty"`
	Belong   []string `json:"belong,omitempty"`
	Role     string   `json:"role,omitempty"`
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

func ValidateToken(tokenStr string) (*Claims, error) {
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
	isBlackListed, err := CheckBlacklistToken(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to check token : %w", err)
	}

	if isBlackListed {
		return nil, errors.New("invalid token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func DecodeJwtClaims(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenStr, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func AddBlackListToken(token string, minutes int64) error {
	log.Printf("JWT - Adding token to blacklist for %d minutes", minutes)
	ctx := redis.Ctx
	err := redis.Rdb.Set(ctx, token, "blacklisted", time.Duration(minutes)*time.Minute).Err()
	if err != nil {
		log.Printf("JWT - Failed to blacklist token: %v", err)
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	log.Printf("JWT - Token successfully added to blacklist")
	return nil
}

func CheckBlacklistToken(token string) (bool, error) {
	ctx := redis.Ctx
	val, err := redis.Rdb.Get(ctx, token).Result()
	if err == goredis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}
	return val == "blacklisted", nil
}
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
