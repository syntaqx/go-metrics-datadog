package main

import (
	"log"
	"os"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	datadog "github.com/syntaqx/go-metrics-datadog"
)

func main() {
	// https://docs.datadoghq.com/developers/dogstatsd
	statsdAddr, ok := os.LookupEnv("DD_AGENT_HOST")
	if !ok {
		statsdAddr = "127.0.0.1:8125"
	}

	ddOpts := []datadog.ReporterOption{
		datadog.UseFlushInterval(time.Second * 10),
		datadog.UsePercentiles([]float64{0.25, 0.99}),
	}

	reporter, err := datadog.NewReporter(nil, statsdAddr, ddOpts...)
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
