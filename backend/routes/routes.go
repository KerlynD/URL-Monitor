package routes

import (
	"net/http"

	"github.com/KerlynD/URL-Monitor/backend/handlers"
	"github.com/KerlynD/URL-Monitor/backend/middleware"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func RegisterRoutes() *http.ServeMux {
	/*
		This function registers the routes for the server
	*/
	mux := http.NewServeMux()

	mux.HandleFunc("POST /monitor", handlers.CreateMonitor)
	mux.HandleFunc("GET /monitor", handlers.ListMonitors)
	mux.HandleFunc("GET /monitor/{id}", handlers.GetMonitor)
	mux.HandleFunc("POST /monitor/{id}/check", handlers.TriggerCheck)

	return mux
}

func SetupServer() http.Handler {
	/*
		This function sets up the server and registers the routes
	*/
	mux := RegisterRoutes()

	var handler http.Handler = mux

	// Wrap handler with Datadog tracing
	handler = httptrace.WrapHandler(handler, "url-monitor", "/")

	handler = middleware.JSONMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.MetricsMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)

	return handler
}
