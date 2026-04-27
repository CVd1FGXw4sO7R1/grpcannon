package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// ConcurrencyLatencyPoint holds aggregated latency stats for a concurrency level.
type ConcurrencyLatencyPoint struct {
	Concurrency int
	Count       int
	P50Ms       float64
	P95Ms       float64
	P99Ms       float64
	AvgMs       float64
	ErrorRate   float64
}

// BuildConcurrencyLatency groups results by concurrency level and computes
// latency percentiles for each group. Results are sorted by concurrency.
func BuildConcurrencyLatency(groups map[int][]Result) []ConcurrencyLatencyPoint {
	if len(groups) == 0 {
		return nil
	}

	points := make([]ConcurrencyLatencyPoint, 0, len(groups))

	for concurrency, results := range groups {
		if len(results) == 0 {
			continue
		}

		var durations []time.Duration
		var totalMs float64
		errCount := 0

		for _, r := range results {
			if r.Err != nil {
				errCount++
				continue
			}
			durations = append(durations, r.Duration)
			totalMs += float64(r.Duration.Milliseconds())
		}

		var p50, p95, p99, avg float64
		if len(durations) > 0 {
			p50 = Percentile(durations, 50)
			p95 = Percentile(durations, 95)
			p99 = Percentile(durations, 99)
			avg = totalMs / float64(len(durations))
		}

		points = append(points, ConcurrencyLatencyPoint{
			Concurrency: concurrency,
			Count:       len(results),
			P50Ms:       p50,
			P95Ms:       p95,
			P99Ms:       p99,
			AvgMs:       avg,
			ErrorRate:   float64(errCount) / float64(len(results)) * 100,
		})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Concurrency < points[j].Concurrency
	})

	return points
}

// WriteConcurrencyLatency writes a concurrency-vs-latency table to w.
func WriteConcurrencyLatency(w io.Writer, points []ConcurrencyLatencyPoint) {
	if len(points) == 0 {
		fmt.Fprintln(w, "no concurrency-latency data")
		return
	}

	fmt.Fprintf(w, "%-12s  %6s  %8s  %8s  %8s  %8s  %8s\n",
		"Concurrency", "Count", "Avg(ms)", "P50(ms)", "P95(ms)", "P99(ms)", "ErrRate")
	fmt.Fprintln(w, "--------------------------------------------------------------------------")

	for _, p := range points {
		fmt.Fprintf(w, "%-12d  %6d  %8.2f  %8.2f  %8.2f  %8.2f  %7.2f%%\n",
			p.Concurrency, p.Count, p.AvgMs, p.P50Ms, p.P95Ms, p.P99Ms, p.ErrorRate)
	}
}
