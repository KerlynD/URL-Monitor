package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(dbPath string) error {
	/*
		Function to initialize the db at backend start
	*/
	var err error

	log.Printf("Starting connection to DB")

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("error opening the DB connection pool: %w", err)
	}

	// Test
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}

	err = createTables()
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	log.Printf("Successfully opened DB connection")
	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func createTables() error {
	/*
		Function to create tables if they dont exist (Initializing)
	*/
	monitorsTable := `
    CREATE TABLE IF NOT EXISTS monitors (
        id TEXT PRIMARY KEY,
        url TEXT NOT NULL,
        check_interval INTEGER NOT NULL,
        created_at DATETIME,
        updated_at DATETIME
    );`

	resultsTable := `
    CREATE TABLE IF NOT EXISTS results (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        monitor_id TEXT NOT NULL,
        status_code INTEGER,
        response_time INTEGER,
        is_up BOOLEAN,
        error TEXT,
        timestamp DATETIME,
        FOREIGN KEY (monitor_id) REFERENCES monitors(id)
    );`

	// Execute
	_, err := db.Exec(monitorsTable)
	if err != nil {
		return fmt.Errorf("error creating monitors table: %w", err)
	}

	_, err = db.Exec(resultsTable)
	if err != nil {
		return fmt.Errorf("error creating results table: %w", err)
	}

	log.Println("Tables created successfully")
	return nil
}
