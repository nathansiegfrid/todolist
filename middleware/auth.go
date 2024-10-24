package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type tokenErrorContextKey struct{}

type tokenVerifier interface {
	VerifyToken(signedToken string) (jwt.Claims, error)
}

// VerifyAuth middleware verifies the Authorization header and extracts user ID from the token.
func VerifyAuth(tokenVerifier tokenVerifier) func(http.Handler) http.Handler {
	verifyRequest := func(r *http.Request) (uuid.UUID, error) {
		authHeaderValue := r.Header.Get("Authorization")
		if authHeaderValue == "" {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Missing authorization header.")
		}

		// Check if Authorization header has "Bearer" prefix and extract the token.
		token := strings.TrimPrefix(authHeaderValue, "Bearer ")
		if token == authHeaderValue {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Invalid authorization header.")
		}

		// Validate the token and extract the claims.
		claims, err := tokenVerifier.VerifyToken(token)
		if err != nil {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Invalid token.")
		}

		// Check if token is expired.
		exp, _ := claims.GetExpirationTime()
		if time.Now().After(exp.Time) {
			// TODO: Add a way to refresh the token. Maybe a separate API endpoint?
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Token is expired.")
		}

		// Get user ID from subject.
		sub, _ := claims.GetSubject()
		uid, _ := uuid.Parse(sub)
		if uid == uuid.Nil {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "User ID not found in token.")
		}
		return uid, nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			uid, err := verifyRequest(r)
			if err != nil {
				ctx = context.WithValue(ctx, tokenErrorContextKey{}, err)
			} else {
				ctx = service.ContextWithUserID(ctx, uid)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth middleware returns an error response if the request is not authenticated.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errCtxVal := r.Context().Value(tokenErrorContextKey{})
		if errCtxVal != nil {
			err, _ := errCtxVal.(error)
			service.WriteError(w, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
