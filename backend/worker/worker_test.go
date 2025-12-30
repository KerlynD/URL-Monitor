package worker

import (
	"testing"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/db"
	"github.com/KerlynD/URL-Monitor/backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	db.InitDB(":memory:")
	defer db.CloseDB()
	m.Run()
}

// setupTestDB initializes a clean database for each test
func setupTestDB(t *testing.T) {
	t.Helper()
	db.CloseDB()
	err := db.InitDB(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.CloseDB()
	})
}

// <----- MAIN TEST FUNCTIONS ----->

func TestCheckAllMonitors_Success(t *testing.T) {
	// Create test data and request, perform request, and validate that the check is triggered successfully
	setupTestDB(t)

	monitor1 := models.MonitorEntry{
		ID:            "monitor1",
		URL:           "https://www.datadoghq.com",
		CheckInterval: 60,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	monitor2 := models.MonitorEntry{
		ID:            "monitor2",
		URL:           "https://www.google.com",
		CheckInterval: 60,
		CreatedAt:     time.Now(),
	}

	db.SaveMonitor(monitor1)
	db.SaveMonitor(monitor2)

	checkAllMonitors()

	result1, err := db.GetLatestResult("monitor1")
	require.NoError(t, err)
	assert.True(t, result1.IsUp)

	result2, err := db.GetLatestResult("monitor2")
	require.NoError(t, err)
	assert.True(t, result2.IsUp)
}

func TestStartMonitorChecker_Success(t *testing.T) {
	// Check without crashing
	setupTestDB(t)

	monitor := models.MonitorEntry{
		ID:            "monitor1",
		URL:           "https://www.datadoghq.com",
		CheckInterval: 60,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	db.SaveMonitor(monitor)

	StartMonitorChecker(500 * time.Millisecond)

	// Wait for at least one check cycle to complete
	time.Sleep(750 * time.Millisecond)

	result, err := db.GetLatestResult("monitor1")
	require.NoError(t, err)
	assert.NotZero(t, result.Timestamp)

	// Background goroutine will continue running briefly after test ends
}
