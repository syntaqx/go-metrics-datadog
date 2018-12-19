package main

import (
	"log"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/syntaqx/go-metrics-datadog"
)

func main() {
	reporter, err := datadog.NewReporter(
		nil,              // Metrics registry, or nil for default
		"127.0.0.1:8125", // DogStatsD UDP address
		time.Second*10,   // Update interval
		datadog.UsePercentiles([]float64{0.25, 0.99}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// configure a prefix, and send the EC2 availability zone as a tag with
	// every metric.
	reporter.Client.Namespace = "test."
	reporter.Client.Tags = append(reporter.Client.Tags, "us-east-1a")

	go reporter.Flush()

	cn := metrics.NewCounter()
	metrics.Register("first.count", cn)
	cn.Inc(1)
}
