package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Outlier holds a single outlier result with its index and duration.
type Outlier struct {
	Index    int
	Duration time.Duration
	Error    error
}

// FindOutliers returns results whose latency exceeds mean + threshold*stddev.
// A minimum of 2 results is required; returns nil otherwise.
func FindOutliers(results []Result, threshold float64) []Outlier {
	if len(results) < 2 {
		return nil
	}
	if threshold <= 0 {
		threshold = 2.0
	}

	var durations []float64
	for _, r := range results {
		if r.Error == nil {
			durations = append(durations, float64(r.Duration.Milliseconds()))
		}
	}
	if len(durations) < 2 {
		return nil
	}

	var sum float64
	for _, d := range durations {
		sum += d
	}
	mean := sum / float64(len(durations))

	var variance float64
	for _, d := range durations {
		diff := d - mean
		variance += diff * diff
	}
	std := 0.0
	if n := float64(len(durations)); n > 1 {
		std = sqrtF(variance / (n - 1))
	}

	limit := mean + threshold*std

	var out []Outlier
	for i, r := range results {
		if r.Error == nil && float64(r.Duration.Milliseconds()) > limit {
			out = append(out, Outlier{Index: i, Duration: r.Duration, Error: r.Error})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Duration > out[j].Duration
	})
	return out
}

// WriteOutliers writes outlier results to w.
func WriteOutliers(w io.Writer, results []Result, threshold float64) {
	out := FindOutliers(results, threshold)
	if len(out) == 0 {
		fmt.Fprintln(w, "No outliers detected.")
		return
	}
	fmt.Fprintf(w, "Outliers (threshold=%.1fσ): %d\n", threshold, len(out))
	fmt.Fprintf(w, "  %-8s  %s\n", "Index", "Latency")
	for _, o := range out {
		fmt.Fprintf(w, "  %-8d  %s\n", o.Index, roundDuration(o.Duration))
	}
}

func sqrtF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x / 2
	for i := 0; i < 50; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
