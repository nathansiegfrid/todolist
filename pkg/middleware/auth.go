package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/pkg/request"
	"github.com/nathansiegfrid/todolist/pkg/response"
	"github.com/nathansiegfrid/todolist/pkg/token"
)

var (
	errHeaderMissing = response.Error(http.StatusUnauthorized, "Authorization header is missing.")
	errHeaderInvalid = response.Error(http.StatusUnauthorized, "Authorization header is not a Bearer token.")
)

type tokenErrorContextKey struct{}

// VerifyAuth middleware verifies the Authorization header and extracts user ID from the token.
func VerifyAuth(jwtService *token.JWTAuth) func(http.Handler) http.Handler {
	// TODO: Use cookie-based authentication for web clients.
	// Cookies support root domain and subdomain sharing.
	verifyRequest := func(r *http.Request) (uuid.UUID, error) {
		authHeaderValue := r.Header.Get("Authorization")
		if authHeaderValue == "" {
			return uuid.Nil, errHeaderMissing
		}
		token := strings.TrimPrefix(authHeaderValue, "Bearer ")
		if token == authHeaderValue {
			return uuid.Nil, errHeaderInvalid
		}
		return jwtService.VerifyToken(token)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID, err := verifyRequest(r)
			if err != nil {
				ctx = context.WithValue(ctx, tokenErrorContextKey{}, err)
			} else {
				ctx = request.ContextWithUserID(ctx, userID)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth middleware returns an error response if the request is not authenticated.
// It must be used after VerifyAuth middleware.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errCtxVal := r.Context().Value(tokenErrorContextKey{})
		if errCtxVal != nil {
			err, _ := errCtxVal.(error)
			response.WriteError(w, response.ErrorResponseFrom(err))
			return
		}
		next.ServeHTTP(w, r)
	})
}
