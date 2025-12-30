package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/KerlynD/URL-Monitor/backend/db"
    "github.com/KerlynD/URL-Monitor/backend/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestMain: Entry point for testing - setups and tears down
func TestMain(m *testing.M) {
	// Init DB for testing
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

// Function to test CreateMonitor handler with valid input
func TestCreateMonitor_Success(t *testing.T) {
	// Create test data and request, perform request, and validate
	setupTestDB(t)

	reqBody := map[string]interface{}{
		"url": "https://www.example.com",
		"check_interval": 60,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/monitor", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateMonitor(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.MonitorEntry
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, reqBody["url"], response.URL)
	assert.Equal(t, reqBody["check_interval"], response.CheckInterval)
	assert.NotEmpty(t, response.ID)
}

func TestCreateMonitor_InvalidURL(t *testing.T) {
	// Create test data and request, perform request, and validate that the URL is invalid
	setupTestDB(t)

	reqBody := map[string]interface{}{
		"url": "invalid-url",
		"check_interval": 60,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/monitor", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateMonitor(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid URL")
}

func TestCreateMonitor_MalformedJSON(t *testing.T) {
	// Create test data and request, perform request, and validate that the JSON is malformed
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodPost, "/monitor", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateMonitor(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListMonitors_Empty(t *testing.T) {
	// Create test data and request, perform request, and validate that the list is empty
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/monitor", nil)
	rr := httptest.NewRecorder()

	ListMonitors(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var monitors []models.MonitorWithStatus
	json.NewDecoder(rr.Body).Decode(&monitors)
	assert.Empty(t, monitors)
}

func TestListMonitors_WithData(t *testing.T) {
	// Create test data and request, perform request, and validate that the list is not empty
	setupTestDB(t)

	testMonitor := models.MonitorEntry{
		ID: "test-monitor",
		URL: "https://www.example.com",
		CheckInterval: 60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.SaveMonitor(testMonitor)

	req := httptest.NewRequest(http.MethodGet, "/monitor", nil)
	rr := httptest.NewRecorder()

	ListMonitors(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var monitors []models.MonitorWithStatus
	json.NewDecoder(rr.Body).Decode(&monitors)
	require.Len(t, monitors, 1)
	assert.Equal(t, testMonitor.URL, monitors[0].URL)
}

func TestGetMonitor_Success(t *testing.T) {
	// Create test data and request, perform request, and validate that the monitor is retrieved successfully
	setupTestDB(t)

	testMonitor := models.MonitorEntry{
		ID: "test-monitor",
		URL: "https://www.datadoghq.com",
		CheckInterval: 60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.SaveMonitor(testMonitor)

	req := httptest.NewRequest(http.MethodGet, "/monitor/test-monitor", nil)
	req.SetPathValue("id", "test-monitor")
	rr := httptest.NewRecorder()

	GetMonitor(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.MonitorWithStatus
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, testMonitor.URL, response.URL)
	assert.Equal(t, testMonitor.ID, response.ID)
}

func TestGetMonitor_NotFound(t *testing.T) {
	// Create test data and request, perform request, and validate that the monitor is not found
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/monitor/not-real-monitor", nil)
	req.SetPathValue("id", "not-real-monitor")
	rr := httptest.NewRecorder()

	GetMonitor(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestTriggerCheck_Success(t *testing.T) {
	// Create test data and request, perform request, and validate that the check is triggered successfully
	setupTestDB(t)

	testMonitor := models.MonitorEntry{
		ID: "test-monitor",
		URL: "https://www.datadoghq.com",
		CheckInterval: 60,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.SaveMonitor(testMonitor)

	req := httptest.NewRequest(http.MethodPost, "/monitor/test-monitor/check", nil)
	req.SetPathValue("id", "test-monitor")
	rr := httptest.NewRecorder()

	TriggerCheck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	
	var result models.MonitorResult
	json.NewDecoder(rr.Body).Decode(&result)
	assert.NotZero(t, result.Timestamp)
	assert.NotZero(t, result.ResponseTime)
}


// <----- HELPER FUNCTIONS ----->

// Test PerformCheck with valid URL
func TestPerformCheck_ValidURL(t *testing.T) {
    testURL := "https://www.datadoghq.com"
    
    result := PerformCheck(testURL)
    
    assert.True(t, result.IsUp, "Datadog website should be up")
    assert.Equal(t, 200, result.StatusCode)
    assert.Greater(t, result.ResponseTime, time.Duration(0))
    assert.Empty(t, result.Error)
}

// Test PerformCheck with invalid URL
func TestPerformCheck_InvalidURL(t *testing.T) {
    testURL := "https://this-domain-definitely-does-not-exist-12345.com"
    
    result := PerformCheck(testURL)
    
    assert.False(t, result.IsUp, "Invalid domain should be down")
    assert.NotEmpty(t, result.Error)
}

// Test PerformCheck with unreachable server
func TestPerformCheck_UnreachableServer(t *testing.T) {
    testURL := "http://localhost:9999"
    
    result := PerformCheck(testURL)
    
    assert.False(t, result.IsUp)
    assert.NotEmpty(t, result.Error)
}

// <---- VALIDATION CASES ----->

func TestCreateMonitor_ValidationCases(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    map[string]interface{}
        expectedStatus int
        expectedError  string
    }{
        {
            name: "Valid HTTP URL",
            requestBody: map[string]interface{}{
                "url":            "http://example.com",
                "check_interval": 60,
            },
            expectedStatus: http.StatusCreated,
            expectedError:  "",
        },
        {
            name: "Valid HTTPS URL",
            requestBody: map[string]interface{}{
                "url":            "https://example.com",
                "check_interval": 60,
            },
            expectedStatus: http.StatusCreated,
            expectedError:  "",
        },
        {
            name: "Missing scheme",
            requestBody: map[string]interface{}{
                "url":            "example.com",
                "check_interval": 60,
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid URL",
        },
        {
            name: "FTP URL (not allowed)",
            requestBody: map[string]interface{}{
                "url":            "ftp://example.com",
                "check_interval": 60,
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Invalid URL",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
			setupTestDB(t)
			
            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest(http.MethodPost, "/monitor", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            rr := httptest.NewRecorder()
            
            CreateMonitor(rr, req)
            
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            if tt.expectedError != "" {
                var response map[string]string
                json.NewDecoder(rr.Body).Decode(&response)
                assert.Contains(t, response["error"], tt.expectedError)
            }
        })
    }
}