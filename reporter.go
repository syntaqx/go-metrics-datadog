package datadog

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	metrics "github.com/rcrowley/go-metrics"
)

// ReporterOption is function-option used during the construction of a *Reporter
type ReporterOption func(*Reporter) error

// Reporter represents a metrics registry, and the statsd client the metrics
// will be flushed to
type Reporter struct {
	// Registry matrices that need to be reported to the Client
	Registry metrics.Registry

	// Client is the configured statsd instance
	Client *statsd.Client

	// Time interval between two consecutive Flush calls to store the matrix
	// value to the Client.
	interval time.Duration

	// Reporter type configuration settings
	tags []string
	ss   map[string]int64

	// Optional parameters
	percentiles []float64
	p           []string
}

// NewReporter creates a new Reporter with a pre-configured statsd client.
func NewReporter(r metrics.Registry, addr string, d time.Duration, options ...ReporterOption) (*Reporter, error) {
	if r == nil {
		r = metrics.DefaultRegistry
	}

	client, err := statsd.New(addr)
	if err != nil {
		return nil, err
	}
	reporter := &Reporter{
		Client:   client,
		Registry: r,
		interval: d,
		ss:       make(map[string]int64),
	}
	for _, option := range options {
		if err := option(reporter); err != nil {
			return nil, err
		}
	}
	return reporter, nil
}

// UsePercentiles builds a *Reporter that reports the specified percentiles
// for Histograms and TimedMetrics
func UsePercentiles(percentiles []float64) ReporterOption {
	return func(r *Reporter) error {
		if len(percentiles) == 0 {
			return fmt.Errorf("Must specify at least 1 percentile")
		}
		var err error
		r.percentiles = percentiles
		r.p, err = getPercentileNames(percentiles)
		return err
	}
}

func getPercentileNames(percentiles []float64) ([]string, error) {
	names := make([]string, len(percentiles))
	for i, percentile := range percentiles {
		if percentile <= 0 || percentile >= 1 {
			return nil, fmt.Errorf("Percentile must lie in interval (0,1)")
		}
		names[i] = ".p" + strconv.FormatFloat(percentile, 'f', -1, 64)[2:]
	}
	return names, nil
}

// Flush is a blocking exporter function which reports metrics in the registry
// to the statsd client, flushing every d duration
func (r *Reporter) Flush() {
	defer func() {
		if rec := recover(); rec != nil {
			handlePanic(rec)
		}
	}()

	for range time.Tick(r.interval) {
		if err := r.FlushOnce(); err != nil {
			log.Println(err)
		}
	}
}

// FlushOnce submits a snapshot submission of the registry to DataDog. This can
// be used in a loop similarly to FlushWithInterval for custom error handling or
// data submission variations.
func (r *Reporter) FlushOnce() error {
	r.Registry.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			v := metric.Count()
			l := r.ss[name]
			r.Client.Count(name, v-l, r.tags, 1)
			r.ss[name] = v

		case metrics.Gauge:
			r.Client.Gauge(name, float64(metric.Value()), r.tags, 1)

		case metrics.GaugeFloat64:
			r.Client.Gauge(name, metric.Value(), r.tags, 1)

		case metrics.Histogram:
			ms := metric.Snapshot()

			r.Client.Gauge(name+".count", float64(ms.Count()), r.tags, 1)
			r.Client.Gauge(name+".max", float64(ms.Max()), r.tags, 1)
			r.Client.Gauge(name+".min", float64(ms.Min()), r.tags, 1)
			r.Client.Gauge(name+".mean", ms.Mean(), r.tags, 1)
			r.Client.Gauge(name+".stddev", ms.StdDev(), r.tags, 1)
			r.Client.Gauge(name+".var", ms.Variance(), r.tags, 1)

			if len(r.percentiles) > 0 {
				values := ms.Percentiles(r.percentiles)
				for i, p := range r.p {
					r.Client.Gauge(name+p, values[i], r.tags, 1)
				}
			}

		case metrics.Meter:
			ms := metric.Snapshot()

			r.Client.Gauge(name+".count", float64(ms.Count()), r.tags, 1)
			r.Client.Gauge(name+".rate1", ms.Rate1(), r.tags, 1)
			r.Client.Gauge(name+".rate5", ms.Rate5(), r.tags, 1)
			r.Client.Gauge(name+".rate15", ms.Rate15(), r.tags, 1)
			r.Client.Gauge(name+".mean", ms.RateMean(), r.tags, 1)

		case metrics.Timer:
			ms := metric.Snapshot()

			r.Client.Gauge(name+".count", float64(ms.Count()), r.tags, 1)
			r.Client.Gauge(name+".max", time.Duration(ms.Max()).Seconds()*1000, r.tags, 1)
			r.Client.Gauge(name+".min", time.Duration(ms.Min()).Seconds()*1000, r.tags, 1)
			r.Client.Gauge(name+".mean", time.Duration(ms.Mean()).Seconds()*1000, r.tags, 1)
			r.Client.Gauge(name+".stddev", time.Duration(ms.StdDev()).Seconds()*1000, r.tags, 1)

			if len(r.percentiles) > 0 {
				values := ms.Percentiles(r.percentiles)
				for i, p := range r.p {
					r.Client.Gauge(name+p, time.Duration(values[i]).Seconds()*1000, r.tags, 1)
				}
			}
		}
	})

	return nil
}

func handlePanic(rec interface{}) {
	callers := ""
	for i := 2; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = callers + fmt.Sprintf("%v:%v\n", file, line)
	}
	log.Printf("Recovered from panic: %#v \n%v", rec, callers)
}
