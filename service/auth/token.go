package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nathansiegfrid/todolist-go/service"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret []byte) *JWTService {
	return &JWTService{secret}
}

func (s *JWTService) GenerateToken(userID string, duration time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *JWTService) VerifyToken(signedToken string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, service.Error(http.StatusUnauthorized, "Invalid token.")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !(ok && token.Valid) {
		return nil, service.Error(http.StatusUnauthorized, "Invalid token.")
	}
	return claims, nil
}
