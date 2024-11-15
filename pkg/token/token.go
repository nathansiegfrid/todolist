package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/pkg/response"
)

var (
	errTokenInvalid = response.Error(http.StatusUnauthorized, "Token verification failed.")
	errTokenExpired = response.Error(http.StatusUnauthorized, "Token has expired.")
	errTokenSubject = response.Error(http.StatusUnauthorized, "Token subject is not a valid UUID.")
)

type JWTAuth struct {
	secret []byte
}

func NewJWTAuth(secret []byte) *JWTAuth {
	return &JWTAuth{secret}
}

func (auth *JWTAuth) GenerateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(auth.secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (auth *JWTAuth) VerifyToken(signedToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.secret, nil
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
	userID, _ := uuid.Parse(sub)
	if userID == uuid.Nil {
		return uuid.Nil, errTokenSubject
	}
	return userID, nil
}
