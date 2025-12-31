package worker

import (
	"log"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/db"
	"github.com/KerlynD/URL-Monitor/backend/handlers"
	"github.com/KerlynD/URL-Monitor/backend/metrics"
	tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

/*
Function to start a background goroutine that checks all monitors at a given interval
*/
func StartMonitorChecker(interval time.Duration) {
	/*
		This function generates a ticker based off the given interval, then enters a goroutine
		that immediately checks all monitors then loops off each tick checking each monitor.
	*/

	ticker := time.NewTicker(interval)

	go func() {
		checkAllMonitors()

		for range ticker.C {
			checkAllMonitors()
		}
	}()

	log.Println("Monitor checker started with interval", interval)
}

func checkAllMonitors() {
	/*
		This function gets all monitors from the database, then loops through each monitor,
		calling handlers.PerformCheck() to check the monitor. We then save the result to the db.
	*/
	
	span := tracer.StartSpan("worker.check_all_monitors",
		tracer.SpanType("worker"),
		tracer.ResourceName("check_cycle"),
	)
	defer span.Finish()

	monitors, err := db.GetAllMonitors()
	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.message", err.Error())
		log.Printf("Error getting all monitors: %v", err)
		if metrics.Client != nil {
			metrics.Client.Incr("worker.get_monitors.error", nil, 1.0)
		}
		return
	}

	// Track how many monitors are being checked
	if metrics.Client != nil {
		metrics.Client.Gauge("monitors.active", float64(len(monitors)), nil, 1.0)
	}

	for _, monitor := range monitors {
		checkSpan := tracer.StartSpan("worker.check_monitor",
			tracer.ChildOf(span.Context()),
			tracer.Tag("monitor.url", monitor.URL),
			tracer.Tag("monitor.id", monitor.ID),
		)

		// Track check attempt
		if metrics.Client != nil {
			metrics.Client.Incr("checks.performed",
				[]string{"url:" + monitor.URL}, 1.0)
		}

		// Time the check operation
		startTime := time.Now()
		result := handlers.PerformCheck(monitor.URL)
		checkDuration := time.Since(startTime)

		checkSpan.SetTag("check.isUp", result.IsUp)
		checkSpan.SetTag("check.responseTime", result.ResponseTime)
		checkSpan.SetTag("check.statusCode", result.StatusCode)

		// Record check duration
		if metrics.Client != nil {
			metrics.Client.Timing("checks.duration",
				checkDuration,
				[]string{"url:" + monitor.URL}, 1.0)
		}

		// Record the actual response time from the URL
		if metrics.Client != nil {
			metrics.Client.Timing("checks.response_time",
				result.ResponseTime,
				[]string{"url:" + monitor.URL}, 1.0)
		}

		// Track success/failure
		if metrics.Client != nil {
			if result.IsUp {
				metrics.Client.Incr("checks.success",
					[]string{"url:" + monitor.URL}, 1.0)
			} else {
				metrics.Client.Incr("checks.failure",
					[]string{"url:" + monitor.URL}, 1.0)
			}
		}

		saveResultSpan := tracer.StartSpan("db.save_result",
			tracer.ChildOf(checkSpan.Context()),
			tracer.Tag("monitor.id", monitor.ID),
		)
		err = db.SaveResult(monitor.ID, result)
		saveResultSpan.Finish()

		if err != nil {
			saveResultSpan.SetTag("error", true)
			saveResultSpan.SetTag("error.message", err.Error())
			
			log.Printf("Error saving result for monitor %s: %v", monitor.ID, err)

			if metrics.Client != nil {
				metrics.Client.Incr("worker.save_result.error", nil, 1.0)
			}
			continue
		}

		checkSpan.Finish()
		log.Printf("Checked monitor %s, result: %+v", monitor.ID, result)
	}

	// Track cycle completion
	if metrics.Client != nil {
		metrics.Client.Incr("worker.check_cycle.complete", nil, 1.0)
	}
}
