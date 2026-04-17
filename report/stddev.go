package report

import (
	"fmt"
	"io"
	"math"
	"time"
)

// StdDev calculates the standard deviation of result durations in milliseconds.
func StdDev(results []Result) float64 {
	if len(results) == 0 {
		return 0
	}

	var sum float64
	for _, r := range results {
		sum += float64(r.Duration / time.Millisecond)
	}
	mean := sum / float64(len(results))

	var variance float64
	for _, r := range results {
		diff := float64(r.Duration/time.Millisecond) - mean
		variance += diff * diff
	}
	variance /= float64(len(results))
	return math.Sqrt(variance)
}

// WriteStdDev writes standard deviation stats to w.
func WriteStdDev(w io.Writer, results []Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "StdDev: no results")
		return err
	}

	var sum float64
	for _, r := range results {
		sum += float64(r.Duration / time.Millisecond)
	}
	mean := sum / float64(len(results))
	std := StdDev(results)

	_, err := fmt.Fprintf(w, "StdDev Report\n-------------\nMean:   %.2f ms\nStdDev: %.2f ms\nCV:     %.2f%%\n",
		mean, std, safeCV(std, mean))
	return err
}

func safeCV(std, mean float64) float64 {
	if mean == 0 {
		return 0
	}
	return (std / mean) * 100
}
