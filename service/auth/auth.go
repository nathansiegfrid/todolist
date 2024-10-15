package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nathansiegfrid/todolist-go/service"
)

type Service struct {
	jwtSecret []byte
}

func NewService(jwtSecret []byte) *Service {
	return &Service{jwtSecret}
}

func (s *Service) GenerateToken(userID string, duration time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *Service) VerifyToken(signedToken string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
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
