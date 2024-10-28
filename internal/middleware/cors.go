package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func CORSAllowOrigins(origins ...string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-CSRF-Token", "X-Request-ID"},
		AllowCredentials: true,
	})
}
