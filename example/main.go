package main

import (
	"log"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/syntaqx/go-metrics-datadog"
)

func main() {
	rep, err := datadog.NewReporter(
		nil,              // Metrics registry, or nil for default
		"127.0.0.1:8125", // DogStatsD UDP address
		time.Second*10,   // Update interval
	)
	if err != nil {
		log.Fatal(err)
	}

	// configure a prefix, and send the EC2 availability zone as a tag with
	// every metric.
	rep.Client.Namespace = "test."
	rep.Client.Tags = append(reg.Client.Tags, "us-east-1a")

	go reg.Flush()

	cn := metrics.NewCounter()
	metrics.Register("first.count", cn)
	cn.Inc(1)
}
