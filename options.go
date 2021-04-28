package datadog

import (
	"fmt"
	"strconv"
	"time"
)

// ReporterOption is function-option used during the construction of a *Reporter
type ReporterOption func(*Reporter) error

// UseFlushInterval configures the flush tick interval to use with `Flush`
func UseFlushInterval(d time.Duration) ReporterOption {
	return func(r *Reporter) error {
		r.interval = d
		return nil
	}
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
