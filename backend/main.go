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
			1. Init DB
			2. Init Logger
			3. Init Metrics
			3. Start Monitor Checker
			4. Setup Routes
			5. Configure HTTP
			6. Start Server in goroutine
			7. Shutdown on interupt
	*/
	// Init Logger (path relative to backend directory)
	logFilePath := "../logs/url-monitor.log"
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
        tracer.WithEnv("dev"),                       
        tracer.WithServiceVersion("1.0.0"),          
        tracer.WithAgentAddr("localhost:8126"),      
        tracer.WithAnalytics(true),                  
        tracer.WithRuntimeMetrics(),                 
    )
    defer tracer.Stop()

	// Init Metrics
	err = metrics.InitMetrics("127.0.0.1:8125")
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
