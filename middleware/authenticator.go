package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
)

type tokenVerifier interface {
	VerifyToken(signedToken string) (jwt.Claims, error)
}

func Authenticator(tokenVerifier tokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				service.WriteError(w, service.Error(http.StatusUnauthorized, "Missing Authorization header."))
				return
			}

			// Check if Authorization header has "Bearer" prefix and extract the token.
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				service.WriteError(w, service.Error(http.StatusUnauthorized, "Invalid Authorization header."))
				return
			}

			// Validate the token and extract the claims.
			claims, err := tokenVerifier.VerifyToken(token)
			if err != nil {
				service.WriteError(w, service.Error(http.StatusUnauthorized, "Invalid token."))
				return
			}

			// Get user ID from subject.
			sub, _ := claims.GetSubject()
			userID, _ := uuid.Parse(sub)
			if userID == uuid.Nil {
				service.WriteError(w, service.Error(http.StatusUnauthorized, "Malformed token."))
				return
			}

			ctx := service.ContextWithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
