package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// PercentileRank returns the percentage of values less than or equal to the given duration.
func PercentileRank(results []Result, target time.Duration) float64 {
	if len(results) == 0 {
		return 0
	}
	durations := make([]float64, 0, len(results))
	for _, r := range results {
		durations = append(durations, float64(r.Duration))
	}
	sort.Float64s(durations)
	t := float64(target)
	count := 0
	for _, d := range durations {
		if d <= t {
			count++
		}
	}
	return float64(count) / float64(len(durations)) * 100.0
}

// WritePercentileRank writes percentile rank for common SLA thresholds.
func WritePercentileRank(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No results to compute percentile ranks.")
		return
	}
	thresholds := []time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		250 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	fmt.Fprintln(w, "Percentile Ranks (% requests completed within threshold):")
	fmt.Fprintf(w, "%-12s  %s\n", "Threshold", "% Requests")
	fmt.Fprintln(w, "-----------------------------")
	for _, t := range thresholds {
		rank := PercentileRank(results, t)
		fmt.Fprintf(w, "%-12s  %.2f%%\n", t, rank)
	}
}
