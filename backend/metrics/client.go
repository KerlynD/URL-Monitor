package metrics

import (
	"log"

	"github.com/DataDog/datadog-go/v5/statsd"
)

var (
	// Need client to be global so it can be used in other packages
	Client *statsd.Client
)

// Init the client
func InitMetrics(addr string) error {
	/*
		This function initializes the metrics client with the given address.
	*/
	var err error

	Client, err = statsd.New(addr,
		statsd.WithNamespace("url_monitor."),
		statsd.WithTags([]string{
			"service:url-monitor",
			"env:dev",
		}),
	)

	if err != nil {
		return err
	}

	log.Println("Metrics client initialized")
	return nil
}

func CloseMetrics() {
	if Client != nil {
		Client.Close()
	}
}
