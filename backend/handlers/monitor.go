package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/db"
	"github.com/KerlynD/URL-Monitor/backend/metrics"
	"github.com/KerlynD/URL-Monitor/backend/models"
	"github.com/google/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

/*
Function to create a monitor that will monitor over the URL given in the
request body and add it to the db of monitors. This function returns
and HTTP response.
*/
func CreateMonitor(response http.ResponseWriter, request *http.Request) {
	/*
		This function parses the request body, validates the URL in the body
		generates a unique ID for the monitor, and saves it to the db
	*/

	span, _ := tracer.StartSpanFromContext(request.Context(), "handler.create_monitor")
	defer span.Finish()

	var req struct {
		URL           string `json:"url"`
		CheckInterval int    `json:"check_interval"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		// Return 400
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	span.SetTag("monitor.url", req.URL)
	span.SetTag("monitor.check_interval", req.CheckInterval)

	validationSpan := tracer.StartSpan("handler.url.validation", tracer.ChildOf(span.Context()))
	parsedURL, err := url.Parse(req.URL)
	validationSpan.Finish()

	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		span.SetTag("error", true)
		// Return 400
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Invalid URL format",
		})
		return
	}

	id := uuid.New().String()

	monitor := models.MonitorEntry{
		ID:            id,
		URL:           req.URL,
		CheckInterval: req.CheckInterval,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	saveSpan := tracer.StartSpan("db.save_monitor", tracer.ChildOf(span.Context()))
	err = db.SaveMonitor(monitor)
	saveSpan.Finish()

	if err != nil {
		// Return 500
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Failed to save monitor to DB",
		})
		return
	}

	// Track monitor creation
	if metrics.Client != nil {
		metrics.Client.Incr("monitors.created", nil, 1.0)

		monitors, _ := db.GetAllMonitors()
		metrics.Client.Gauge("monitors.total", float64(len(monitors)), nil, 1.0)
	}

	// Return 201
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(monitor)
}

/*
Function to list all current monitors requested by the client
*/
func ListMonitors(response http.ResponseWriter, request *http.Request) {
	/*
		This function requests all monitors from the database & builds a response status
		for each.
	*/

	span, _ := tracer.StartSpanFromContext(request.Context(), "handler.list_monitors")
	defer span.Finish()

	monitors, err := db.GetAllMonitors()
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		// Return 500
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Failed to fetch monitors from db",
		})
		return
	}


	var monitorResponses []models.MonitorWithStatus

	for _, monitor := range monitors {
		getLatestResultSpan := tracer.StartSpan("db.get_latest_result", tracer.ChildOf(span.Context()))
		result, err := db.GetLatestResult(monitor.ID)
		getLatestResultSpan.Finish()

		status := models.MonitorWithStatus{
			MonitorEntry: monitor,
			LastResult:   nil, // Default
		}

		if err == nil {
			status.LastResult = &result
		}

		monitorResponses = append(monitorResponses, status)
	}

	// Return 200
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(monitorResponses)
}

/*
Function to retrieve a URL monitor for the client
*/
func GetMonitor(response http.ResponseWriter, request *http.Request) {
	/*
		This function Extracts the ID from the request, makes a DB call to get
		the requested monitor, checks its latest result, and returns the status
		of that monitor.
	*/
	span, _ := tracer.StartSpanFromContext(request.Context(), "handler.get_monitor")
	defer span.Finish()

	id := request.PathValue("id")

	monitor, err := db.GetMonitor(id)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Monitor not found",
		})
		return
	}

	getLatestResultSpan := tracer.StartSpan("db.get_latest_result", tracer.ChildOf(span.Context()))
	result, err := db.GetLatestResult(id)
	getLatestResultSpan.Finish()

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())

		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Failed to get latest result",
		})
	}

	status := models.MonitorWithStatus{
		MonitorEntry: monitor,
		LastResult:   nil,
	}

	if err == nil {
		status.LastResult = &result
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(status)
}

func TriggerCheck(response http.ResponseWriter, request *http.Request) {
	/*
		This function extracts the ID from the request, tries to get the monitor,
		performs an HTTP check with performCheck() (helper below) and saves result
		to database
	*/
	span, _ := tracer.StartSpanFromContext(request.Context(), "handler.trigger_check")
	defer span.Finish()

	id := request.PathValue("id")

	getMonitorSpan := tracer.StartSpan("db.get_monitor", tracer.ChildOf(span.Context()))
	monitor, err := db.GetMonitor(id)
	getMonitorSpan.Finish()

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		response.WriteHeader(http.StatusNotFound)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Monitor not found",
		})
		return
	}

	checkSpan := tracer.StartSpan("handler.perform_check", tracer.ChildOf(span.Context()))
	result := PerformCheck(monitor.URL)
	checkSpan.Finish()

	err = db.SaveResult(id, result)
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())

		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "Failed to save result",
		})
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func PerformCheck(targetURL string) models.MonitorResult {
	/*
		This function creates an HTTP client (with timeout), makes a GET request,
		checks duration and returns the result
	*/
	span := tracer.StartSpan("http.check", 
		tracer.ResourceName("GET " + targetURL),
		tracer.SpanType("http"),
	)
	defer span.Finish()

	client := httptrace.WrapClient(&http.Client{
		Timeout: 10 * time.Second,
	})

	startTime := time.Now()

	resp, err := client.Get(targetURL)

	responseTime := time.Since(startTime)

	result := models.MonitorResult{
		Timestamp:    time.Now(),
		ResponseTime: responseTime,
	}

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		result.IsUp = false
		result.Error = err.Error()
		result.StatusCode = 0
	} else {
		defer resp.Body.Close()
		result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300
		result.StatusCode = resp.StatusCode
		result.Error = ""
	}

	return result
}
