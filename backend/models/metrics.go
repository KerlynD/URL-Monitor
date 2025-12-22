package models

import (
	"time"

)

type MonitorEntry struct {
	ID string `json:"id"`
	URL string `json:"url"`
	CheckInterval int `json:"check_interval"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MonitorResult struct {
	StatusCode int `json:"status_code"`
	ResponseTime time.Duration `json:"response_time"`
	IsUp bool `json:"is_up"`
	Error string `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

type MonitorWithStatus struct {
	MonitorEntry
	LastResult *MonitorResult `json:"last_result"`
}