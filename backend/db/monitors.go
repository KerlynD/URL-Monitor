package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/KerlynD/URL-Monitor/backend/models"
	tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

/*
Function to save monitor into sqlite DB
*/
func SaveMonitor(entry models.MonitorEntry) error {
	span := tracer.StartSpan("db.save_monitor",
		tracer.SpanType("sql"),
		tracer.ResourceName("INSERT INTO monitors"),
	)
	defer span.Finish()

	query := `
	INSERT OR REPLACE INTO monitors (id, url, check_interval, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		entry.ID,
		entry.URL,
		entry.CheckInterval,
		entry.CreatedAt,
		entry.UpdatedAt,
	)

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		return fmt.Errorf("eror saving monitor to db: %w", err)
	}

	log.Printf("Monitor %s saved successfully", entry.ID)
	return nil
}

/*
Function to get a single monitor from the DB
*/
func GetMonitor(id string) (models.MonitorEntry, error) {
	span := tracer.StartSpan("db.get_monitor",
		tracer.SpanType("sql"),
		tracer.ResourceName("SELECT * FROM monitors WHERE id = ?"),
	)
	defer span.Finish()

	query := `SELECT id, url, check_interval, created_at, updated_at 
              FROM monitors 
              WHERE id = ?`

	var monitor models.MonitorEntry

	row := db.QueryRow(query, id)

	err := row.Scan( // <- Scans row to make sure struct fields match & puts into fields
		&monitor.ID,
		&monitor.URL,
		&monitor.CheckInterval,
		&monitor.CreatedAt,
		&monitor.UpdatedAt,
	)

	if err == sql.ErrNoRows || err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		return models.MonitorEntry{}, fmt.Errorf("error querying db for monitor: %w", err)
	}

	return monitor, nil
}

/*
Function to get all current URLs we are monitoring
*/
func GetAllMonitors() ([]models.MonitorEntry, error) {
	span := tracer.StartSpan("db.get_all_monitors",
		tracer.SpanType("sql"),
		tracer.ResourceName("SELECT * FROM monitors"),
	)
	defer span.Finish()

	query := `SELECT id, url, check_interval, created_at, updated_at 
              FROM monitors`

	rows, err := db.Query(query)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		return nil, fmt.Errorf("error querying db for monitors: %w", err)
	}
	defer rows.Close()

	var monitors []models.MonitorEntry

	for rows.Next() {
		var monitor models.MonitorEntry

		err := rows.Scan(
			&monitor.ID,
			&monitor.URL,
			&monitor.CheckInterval,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
		)

		if err != nil {
			span.SetTag("error", true)
			span.SetTag("error.message", err.Error())
			return nil, fmt.Errorf("error scanning monitor: %w", err)
		}

		monitors = append(monitors, monitor)
	}

	if err = rows.Err(); err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		return nil, fmt.Errorf("error iterating through monitors: %w", err)
	}

	log.Printf("Monitors retrieved successfully")
	return monitors, nil
}
