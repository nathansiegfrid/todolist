package api

import (
	"net/http"
	"sort"
	"strings"
)

// MethodHandler maps HTTP methods to handler functions.
type MethodHandler map[string]http.HandlerFunc

func (hmap MethodHandler) HandlerFunc() http.HandlerFunc {
	// Define HTTP methods that can be used as keys in MethodHandler.
	// Unknown keys will return false by default.
	validKeys := map[string]bool{
		"HEAD":    true,
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"PATCH":   true,
		"DELETE":  true,
		"OPTIONS": false,
		"CONNECT": false,
		"TRACE":   false,
	}

	// If HEAD handler is not defined, use GET handler.
	if hmap["HEAD"] == nil {
		hmap["HEAD"] = hmap["GET"]
	}

	// Get list of allowed methods.
	methods := make([]string, 0, len(hmap))
	for k, v := range hmap {
		// Check if method is valid and handler is not nil.
		if validKeys[k] && v != nil {
			methods = append(methods, k)
		} else {
			delete(hmap, k)
		}
	}
	sort.Strings(methods) // Unnecessary, but makes it consistent across app restarts.
	allowHeaderValue := strings.Join(methods, ", ")

	return func(w http.ResponseWriter, r *http.Request) {
		serveHTTP := hmap[r.Method]
		if serveHTTP == nil {
			// If no handler found, respond with 405 Method Not Allowed.
			w.Header().Set("Allow", allowHeaderValue)
			MethodNotAllowed(w, r)
			return
		}
		serveHTTP(w, r)
	}
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	err := Error(http.StatusNotFound, "Resource not found.")
	WriteError(w, err)
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	err := Error(http.StatusMethodNotAllowed, "Method not allowed.")
	WriteError(w, err)
}
