package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type tokenErrorContextKey struct{}

type tokenVerifier interface {
	VerifyToken(signedToken string) (string, error)
}

// VerifyAuth middleware verifies the Authorization header and extracts user ID from the token.
func VerifyAuth(tokenVerifier tokenVerifier) func(http.Handler) http.Handler {
	verifyRequest := func(r *http.Request) (uuid.UUID, error) {
		authHeaderValue := r.Header.Get("Authorization")
		if authHeaderValue == "" {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Authorization header not found.")
		}

		// Check if Authorization header has "Bearer" prefix and extract the token.
		token := strings.TrimPrefix(authHeaderValue, "Bearer ")
		if token == authHeaderValue {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Authorization header is not a Bearer token.")
		}

		// Verify the token and extract the subject.
		sub, err := tokenVerifier.VerifyToken(token)
		if err != nil {
			return uuid.Nil, service.Error(http.StatusUnauthorized, err.Error())
		}

		uid, _ := uuid.Parse(sub)
		if uid == uuid.Nil {
			return uuid.Nil, service.Error(http.StatusUnauthorized, "Token subject is not a valid UUID.")
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
		// TODO: Add permission-based authorization.
		// 0 = no access
		// 1 = read access
		// 2 = read-create access
		// 3 = read-create-update-delete access
		// 4 = admin access (grant/revoke access to non-admin users)
		// 5 = owner access (grant/revoke access, change owner)
		next.ServeHTTP(w, r)
	})
}
