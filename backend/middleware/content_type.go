package middleware

import (
	"net/http"
)

// JSONMiddleware so we can get the right content type header for all resp
func JSONMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}