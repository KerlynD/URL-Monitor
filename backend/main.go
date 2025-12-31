package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/db"
	"github.com/KerlynD/URL-Monitor/backend/logging"
	"github.com/KerlynD/URL-Monitor/backend/metrics"
	"github.com/KerlynD/URL-Monitor/backend/routes"
	"github.com/KerlynD/URL-Monitor/backend/worker"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	/*
		Main entry to the backend:
			1. Detect if running in Docker
			2. Init Logger
			3. Init DB
			4. Init Metrics
			5. Start Monitor Checker
			6. Setup Routes
			7. Configure HTTP
			8. Start Server in goroutine
			9. Shutdown on interupt
	*/

	// Detect Docker Container
	datadogHost := getDatadogHost()

	// Init Logger
	logFilePath := getLogPath()
	err := logging.InitLogger(logFilePath)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logging.Close()

	// Init DB
	dbPath := "db/monitor.db"
	err = db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize tracer
	tracer.Start(
		tracer.WithService("url-monitor"),
		tracer.WithEnv(getEnv()),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithAgentAddr(datadogHost+":8126"),
		tracer.WithAnalytics(true),
		tracer.WithRuntimeMetrics(),
	)
	defer tracer.Stop()

	// Init Metrics
	metricsHost := datadogHost + ":8125"
	err = metrics.InitMetrics(metricsHost)
	if err != nil {
		log.Printf("Failed to init metrics: %v", err)
	} else {
		log.Println("Metrics initialized")
		defer metrics.CloseMetrics()
	}

	// Start Monitor Checker
	checkInterval := 30 * time.Second
	worker.StartMonitorChecker(checkInterval)

	// Setup Routes
	handler := routes.SetupServer()

	port := "8080"

	// Configure HTTP
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Start Server in goroutine
	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Shutdown on interupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down")
}

// Helper function to determine the Datadog Agent Address
func getDatadogHost() string {
	// Check ENV
	if host := os.Getenv("DD_AGENT_HOST"); host != "" {
		return host
	}

	// Default to host.docker.internal
	return "host.docker.internal"
}

// Helper function to determine the Datadog Environment
func getEnv() string {
	if env := os.Getenv("DD_ENV"); env != "" {
		return env
	}
	return "dev"
}

// Helper function to get the log file path
func getLogPath() string {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "logs/url-monitor.log"
	}
	// On host: use parent directory's logs folder
	return "../logs/url-monitor.log"
}
