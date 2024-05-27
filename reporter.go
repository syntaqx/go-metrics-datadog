package datadog

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	defaultFlushInterval = time.Second * 10
)

// Expect the tags in the pattern
// namespace.metricName[tag1:value1,tag2:value2,etc....]
var tagPattern = regexp.MustCompile("([\\w\\.]+)\\[([\\w\\W]+)\\]")

// Reporter wraps a metrics registry with a given statsd client.
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
func NewReporter(r metrics.Registry, addr string, options ...ReporterOption) (*Reporter, error) {
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
		interval: defaultFlushInterval,
		ss:       make(map[string]int64),
	}
	for _, option := range options {
		if err := option(reporter); err != nil {
			return nil, err
		}
	}
	return reporter, nil
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
	reportCount := func(name string, tags []string, value int64) {
		metric := getMetric(name, tags)
		last := r.ss[metric]
		r.Client.Count(name, value-last, tags, 1)
		r.ss[metric] = value
	}
	r.Registry.Each(func(metricName string, i interface{}) {
		name, tags := r.splitNameAndTags(metricName)

		switch metric := i.(type) {
		case metrics.Counter:
			reportCount(name, tags, metric.Count())

		case metrics.Gauge:
			r.Client.Gauge(name, float64(metric.Value()), tags, 1)

		case metrics.GaugeFloat64:
			r.Client.Gauge(name, metric.Value(), tags, 1)

		case metrics.Histogram:
			ms := metric.Snapshot()

			reportCount(name+".count", tags, ms.Count())
			r.Client.Gauge(name+".max", float64(ms.Max()), tags, 1)
			r.Client.Gauge(name+".min", float64(ms.Min()), tags, 1)
			r.Client.Gauge(name+".mean", ms.Mean(), tags, 1)
			r.Client.Gauge(name+".stddev", ms.StdDev(), tags, 1)
			r.Client.Gauge(name+".sum", float64(ms.Sum()), tags, 1)
			r.Client.Gauge(name+".var", ms.Variance(), tags, 1)

			if len(r.percentiles) > 0 {
				values := ms.Percentiles(r.percentiles)
				for i, p := range r.p {
					r.Client.Gauge(name+p, values[i], tags, 1)
				}
			}

		case metrics.Meter:
			ms := metric.Snapshot()

			reportCount(name+".count", tags, ms.Count())
			r.Client.Gauge(name+".rate1", ms.Rate1(), tags, 1)
			r.Client.Gauge(name+".rate5", ms.Rate5(), tags, 1)
			r.Client.Gauge(name+".rate15", ms.Rate15(), tags, 1)
			r.Client.Gauge(name+".mean", ms.RateMean(), tags, 1)

		case metrics.Timer:
			ms := metric.Snapshot()

			reportCount(name+".count", tags, ms.Count())
			r.Client.Gauge(name+".max", time.Duration(ms.Max()).Seconds()*1000, tags, 1)
			r.Client.Gauge(name+".min", time.Duration(ms.Min()).Seconds()*1000, tags, 1)
			r.Client.Gauge(name+".mean", time.Duration(ms.Mean()).Seconds()*1000, tags, 1)
			r.Client.Gauge(name+".stddev", time.Duration(ms.StdDev()).Seconds()*1000, tags, 1)
			r.Client.Gauge(name+".sum", float64(ms.Sum()), tags, 1)

			if len(r.percentiles) > 0 {
				values := ms.Percentiles(r.percentiles)
				for i, p := range r.p {
					r.Client.Gauge(name+p, time.Duration(values[i]).Seconds()*1000, tags, 1)
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

func (r *Reporter) splitNameAndTags(metric string) (string, []string) {
	if res := tagPattern.FindStringSubmatch(metric); len(res) == 3 {
		if r.tags == nil {
			return res[1], append(strings.Split(res[2], ","))
		} else {
			return res[1], append(strings.Split(res[2], ","), r.tags...)
		}
	}
	return metric, r.tags
}

func getMetric(name string, tags []string) string {
	var labels string
	for _, tag := range tags {
		labels += tag
	}
	return fmt.Sprintf("%s[%s]", name, labels)
}
