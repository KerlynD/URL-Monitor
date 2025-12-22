package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/models"
)

/*
Function to save result into sqlite DB
*/
func SaveResult(monitorID string, result models.MonitorResult) error {
	query := `
    INSERT INTO results (monitor_id, status_code, response_time, is_up, error, timestamp)
    VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		monitorID,
		result.StatusCode,
		result.ResponseTime.Milliseconds(),
		result.IsUp,
		result.Error,
		result.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("error saving result to db: %w", err)
	}

	log.Printf("Result %s saved successfully", monitorID)
	return nil
}

func GetLatestResult(monitorID string) (models.MonitorResult, error) {
	query := `
    SELECT status_code, response_time, is_up, error, timestamp
    FROM results
    WHERE monitor_id = ?
    ORDER BY timestamp DESC
    LIMIT 1`

	var result models.MonitorResult
	var responseTimeMs int64

	row := db.QueryRow(query, monitorID)

	err := row.Scan(
		&result.StatusCode,
		&responseTimeMs,
		&result.IsUp,
		&result.Error,
		&result.Timestamp,
	)

	if err == sql.ErrNoRows || err != nil {
		return models.MonitorResult{}, fmt.Errorf("error querying db for latest result: %w", err)
	}

	result.ResponseTime = time.Duration(responseTimeMs) * time.Millisecond

	log.Printf("Latest result %s retrieved successfully", monitorID)
	return result, nil
}
