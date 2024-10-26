package middleware

import (
	"net/http"
	"strings"
)

func Heartbeat(endpoint string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.EqualFold(r.URL.Path, endpoint) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("200 OK"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
