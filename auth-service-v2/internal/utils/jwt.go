package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret string
	ttl    time.Duration
}

func NewJWTManagerFromEnv() *JWTManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	ttl := 60
	if v := os.Getenv("JWT_TTL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ttl = n
		}
	}
	return &JWTManager{secret: secret, ttl: time.Minute * time.Duration(ttl)}
}

func NewJWTManager(secret string, ttlMins int) *JWTManager {
	return &JWTManager{secret: secret, ttl: time.Minute * time.Duration(ttlMins)}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (j *JWTManager) Generate(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%v", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "auth-service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
