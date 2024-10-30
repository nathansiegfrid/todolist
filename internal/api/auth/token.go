package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/internal/api"
)

var (
	errTokenInvalid = api.Error(http.StatusUnauthorized, "Token verification failed.")
	errTokenExpired = api.Error(http.StatusUnauthorized, "Token has expired.")
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret []byte) *JWTService {
	return &JWTService{secret}
}

func (s *JWTService) GenerateToken(subject uuid.UUID, duration time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   subject.String(),
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

func (s *JWTService) VerifyToken(signedToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, errTokenInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, errTokenInvalid
	}

	exp, _ := claims.GetExpirationTime() // Error value is always nil.
	if time.Now().After(exp.Time) {
		return uuid.Nil, errTokenExpired
	}

	sub, _ := claims.GetSubject() // Error value is always nil.
	uid, _ := uuid.Parse(sub)
	return uid, nil
}
