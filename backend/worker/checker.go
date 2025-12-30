package worker

import (
	"log"
	"time"

	"github.com/KerlynD/URL-Monitor/backend/db"
	"github.com/KerlynD/URL-Monitor/backend/handlers"
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

	monitors, err := db.GetAllMonitors()
	if err != nil {
		log.Printf("Error getting all monitors: %v", err)
		return
	}

	for _, monitor := range monitors {
		result := handlers.PerformCheck(monitor.URL)

		err = db.SaveResult(monitor.ID, result)
		if err != nil {
			log.Printf("Error saving result for monitor %s: %v", monitor.ID, err)
			continue
		}

		log.Printf("Checked monitor %s, result: %+v", monitor.ID, result)
	}
}