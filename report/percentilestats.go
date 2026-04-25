package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// PercentileStats holds a snapshot of common percentile values for a result set.
type PercentileStats struct {
	P50  time.Duration
	P75  time.Duration
	P90  time.Duration
	P95  time.Duration
	P99  time.Duration
	P999 time.Duration
	Min  time.Duration
	Max  time.Duration
}

// BuildPercentileStats computes a full set of percentile statistics from a
// slice of Results. Only successful results are included in the calculation.
// Returns a zero-value PercentileStats when there are no successful results.
func BuildPercentileStats(results []Result) PercentileStats {
	var durations []float64
	for _, r := range results {
		if r.IsSuccess() {
			durations = append(durations, float64(r.Duration.Milliseconds()))
		}
	}
	if len(durations) == 0 {
		return PercentileStats{}
	}
	sort.Float64s(durations)

	toD := func(ms float64) time.Duration {
		return time.Duration(ms * float64(time.Millisecond))
	}

	return PercentileStats{
		P50:  toD(Percentile(durations, 50)),
		P75:  toD(Percentile(durations, 75)),
		P90:  toD(Percentile(durations, 90)),
		P95:  toD(Percentile(durations, 95)),
		P99:  toD(Percentile(durations, 99)),
		P999: toD(Percentile(durations, 99.9)),
		Min:  toD(durations[0]),
		Max:  toD(durations[len(durations)-1]),
	}
}

// WritePercentileStats writes a formatted percentile statistics table to w.
func WritePercentileStats(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no results")
		return
	}
	stats := BuildPercentileStats(results)
	if stats.P50 == 0 && stats.Max == 0 {
		fmt.Fprintln(w, "no successful results")
		return
	}
	fmt.Fprintln(w, "Percentile Statistics (successful requests only)")
	fmt.Fprintln(w, "------------------------------------------------")
	fmt.Fprintf(w, "  Min   : %v\n", roundDuration(stats.Min))
	fmt.Fprintf(w, "  p50   : %v\n", roundDuration(stats.P50))
	fmt.Fprintf(w, "  p75   : %v\n", roundDuration(stats.P75))
	fmt.Fprintf(w, "  p90   : %v\n", roundDuration(stats.P90))
	fmt.Fprintf(w, "  p95   : %v\n", roundDuration(stats.P95))
	fmt.Fprintf(w, "  p99   : %v\n", roundDuration(stats.P99))
	fmt.Fprintf(w, "  p99.9 : %v\n", roundDuration(stats.P999))
	fmt.Fprintf(w, "  Max   : %v\n", roundDuration(stats.Max))
}
