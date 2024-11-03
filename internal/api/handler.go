package api

import (
	"net/http"
	"slices"
	"strings"

	"github.com/samber/lo"
)

// MethodHandler maps HTTP methods to handler functions.
type MethodHandler map[string]http.HandlerFunc

func (hmap MethodHandler) HandlerFunc() http.HandlerFunc {
	// If HEAD handler is not defined, use GET handler.
	if hmap["HEAD"] == nil {
		hmap["HEAD"] = hmap["GET"]
	}

	// Filter out invalid methods and nil handlers.
	hmap = lo.PickByKeys(hmap, []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"})
	hmap = lo.OmitBy(hmap, func(_ string, v http.HandlerFunc) bool { return v == nil })

	// Get list of allowed methods.
	methods := lo.Keys(hmap)
	slices.Sort(methods) // Make the order consistent across app restarts.
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
