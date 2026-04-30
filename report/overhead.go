package report

import (
	"fmt"
	"io"
	"time"
)

// OverheadStats holds timing overhead metrics derived from a set of results.
type OverheadStats struct {
	TotalRequests  int
	SuccessCount   int
	MinLatencyMs   float64
	MaxLatencyMs   float64
	AvgLatencyMs   float64
	OverheadMs     float64 // avg - p50, a proxy for tail overhead
	P50Ms          float64
	P99Ms          float64
}

// CalcOverhead computes overhead statistics from a slice of Results.
func CalcOverhead(results []Result) OverheadStats {
	if len(results) == 0 {
		return OverheadStats{}
	}

	var successes []time.Duration
	var totalMs float64

	for _, r := range results {
		if r.Err == nil {
			successes = append(successes, r.Duration)
			totalMs += float64(r.Duration.Milliseconds())
		}
	}

	stats := OverheadStats{
		TotalRequests: len(results),
		SuccessCount:  len(successes),
	}

	if len(successes) == 0 {
		return stats
	}

	sorted := SortedDurationsMs(successes)

	var minV, maxV float64
	minV = sorted[0]
	maxV = sorted[len(sorted)-1]

	stats.MinLatencyMs = minV
	stats.MaxLatencyMs = maxV
	stats.AvgLatencyMs = totalMs / float64(len(successes))
	stats.P50Ms = Percentile(sorted, 50)
	stats.P99Ms = Percentile(sorted, 99)
	stats.OverheadMs = stats.AvgLatencyMs - stats.P50Ms

	return stats
}

// WriteOverhead writes a human-readable overhead report to w.
func WriteOverhead(w io.Writer, stats OverheadStats) {
	if stats.TotalRequests == 0 {
		fmt.Fprintln(w, "Overhead: no results")
		return
	}
	fmt.Fprintf(w, "Overhead Report\n")
	fmt.Fprintf(w, "  Total Requests : %d\n", stats.TotalRequests)
	fmt.Fprintf(w, "  Successes      : %d\n", stats.SuccessCount)
	fmt.Fprintf(w, "  Min Latency    : %.2f ms\n", stats.MinLatencyMs)
	fmt.Fprintf(w, "  Max Latency    : %.2f ms\n", stats.MaxLatencyMs)
	fmt.Fprintf(w, "  Avg Latency    : %.2f ms\n", stats.AvgLatencyMs)
	fmt.Fprintf(w, "  P50 Latency    : %.2f ms\n", stats.P50Ms)
	fmt.Fprintf(w, "  P99 Latency    : %.2f ms\n", stats.P99Ms)
	fmt.Fprintf(w, "  Tail Overhead  : %.2f ms\n", stats.OverheadMs)
}
