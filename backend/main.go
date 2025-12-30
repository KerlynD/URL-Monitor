package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/db"
	"github.com/KerlynD/URL-Monitor/backend/routes"
	"github.com/KerlynD/URL-Monitor/backend/worker"
)

func main() {
	/*
		Main entry to the backend:
			1. Init DB
			2. Start Monitor Checker
			3. Setup Routes
			4. Configure HTTP
			5. Start Server in goroutine
			6. Shutdown on interupt
	*/

	dbPath := "db/monitor.db"
	err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	checkInterval := 2 * time.Minute
	worker.StartMonitorChecker(checkInterval)

	handler := routes.SetupServer()

	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server starting on http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down")
}
