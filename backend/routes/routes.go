package routes

import (
	"net/http"

	"github.com/KerlynD/URL-Monitor/backend/handlers"
	"github.com/KerlynD/URL-Monitor/backend/middleware"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /monitor", handlers.CreateMonitor)
	mux.HandleFunc("GET /monitor", handlers.ListMonitors)
	mux.HandleFunc("GET /monitor/{id}", handlers.GetMonitor)
	mux.HandleFunc("POST /monitor/{id}/check", handlers.TriggerCheck)

	return mux
}

func SetupServer() http.Handler {
	mux := RegisterRoutes()

	var handler http.Handler = mux

	handler = middleware.JSONMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.MetricsMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)

	return handler
}
