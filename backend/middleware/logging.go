package middleware

import (
	"log"
	"net/http"
	"time"
	"fmt"
	tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		// Extract Trace Context
		span, _ := tracer.StartSpanFromContext(r.Context(), "middleware.logging")
		traceID := ""
		spanID := ""
		if span != nil {
			traceID = fmt.Sprintf("%d", span.Context().TraceID())
			spanID = fmt.Sprintf("%d", span.Context().SpanID())
		}

		log.Printf("[dd.trace_id=%s dd.span_id=%s] Started %s %s", 
            traceID, spanID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		log.Printf("[dd.trace_id=%s dd.span_id=%s] Completed %s %s in %v", 
            traceID, spanID, r.Method, r.URL.Path, duration)
	})
}
