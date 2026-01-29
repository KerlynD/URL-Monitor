package middleware

import (
	"net/http"
	"os"
	"strings"
)

// Function to wrap handlers with CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Get allowed origins from environment variable, or use default
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		var allowedOrigins []string
		
		if allowedOriginsEnv != "" {
			// Split comma-separated origins from env
			allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		} else {
			// Default allowed origins for development
			allowedOrigins = []string{
				"http://localhost:3000",
				"https://url-monitor-gamma.vercel.app",
			}
		}

		origin := r.Header.Get("Origin")
		
		// Check if origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				if origin == "" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}