package middleware

import (
	"net/http"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/metrics"
)

// MetricsMiddleware to track request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	/*
		This middleware tracks request metrics and status codes.
		It increments the request count, tracks durations, status codes, and errors.
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Skip metrics if client not initialized
		if metrics.Client == nil {
			next.ServeHTTP(w, r)
			return
		}

		// Increment request count
		metrics.Client.Incr("http.request.count",
			[]string{
				"method:" + r.Method,
				"path:" + r.URL.Path,
			}, 1.0,
		)

		// Create response wrapper to track status codes
		wrapper := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		// Track Durations
		duration := time.Since(startTime)
		metrics.Client.Timing("http.request.duration",
			duration,
			[]string{
				"method:" + r.Method,
				"path:" + r.URL.Path,
				"status:" + http.StatusText(wrapper.statusCode),
			}, 1.0)

		// Track Status Codes
		metrics.Client.Incr("http.response.status",
			[]string{
				"status_code:" + http.StatusText(wrapper.statusCode),
			}, 1.0)

		// Track Errors
		if wrapper.statusCode >= 400 {
			metrics.Client.Incr("http.response.error",
				[]string{
					"method:" + r.Method,
					"endpoint:" + r.URL.Path,
					"status:" + http.StatusText(wrapper.statusCode),
				}, 1.0)
		}
	})
}

// Wrapper for ResponseWriter to track status codes
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
