package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// FanoutBucket holds aggregated stats for a single concurrency level window.
type FanoutBucket struct {
	Concurrency int
	Count       int
	Successes   int
	Failures    int
	AvgMs       float64
	P99Ms       float64
	RPS         float64
}

// FanoutReport holds the full fan-out analysis across concurrency levels.
type FanoutReport struct {
	Buckets []FanoutBucket
	PeakRPS float64
	Optimal int // concurrency level with best RPS/latency trade-off
}

// BuildFanout groups results by concurrency tag and computes per-level stats.
// Each Result.Concurrency field must be set by the runner.
func BuildFanout(results []Result, levels []int) *FanoutReport {
	if len(results) == 0 || len(levels) == 0 {
		return &FanoutReport{}
	}

	groups := make(map[int][]Result, len(levels))
	for _, r := range results {
		groups[r.Concurrency] = append(groups[r.Concurrency], r)
	}

	buckets := make([]FanoutBucket, 0, len(levels))
	for _, lvl := range levels {
		rs, ok := groups[lvl]
		if !ok || len(rs) == 0 {
			continue
		}
		var successes, failures int
		var durations []time.Duration
		for _, r := range rs {
			if r.Error == nil {
				successes++
			} else {
				failures++
			}
			durations = append(durations, r.Duration)
		}
		sorted := SortedDurationsMs(durations)
		avg := 0.0
		for _, v := range sorted {
			avg += v
		}
		if len(sorted) > 0 {
			avg /= float64(len(sorted))
		}
		p99 := Percentile(sorted, 99)
		span := totalSpan(rs)
		rps := 0.0
		if span > 0 {
			rps = float64(successes) / span.Seconds()
		}
		buckets = append(buckets, FanoutBucket{
			Concurrency: lvl,
			Count:       len(rs),
			Successes:   successes,
			Failures:    failures,
			AvgMs:       avg,
			P99Ms:       p99,
			RPS:         rps,
		})
	}

	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Concurrency < buckets[j].Concurrency
	})

	peak := 0.0
	optimal := 0
	for _, b := range buckets {
		if b.RPS > peak {
			peak = b.RPS
			optimal = b.Concurrency
		}
	}

	return &FanoutReport{Buckets: buckets, PeakRPS: peak, Optimal: optimal}
}

// totalSpan returns the elapsed time from first start to last end in results.
func totalSpan(rs []Result) time.Duration {
	if len(rs) == 0 {
		return 0
	}
	var earliest, latest time.Time
	for i, r := range rs {
		start := r.Timestamp
		end := start.Add(r.Duration)
		if i == 0 || start.Before(earliest) {
			earliest = start
		}
		if end.After(latest) {
			latest = end
		}
	}
	return latest.Sub(earliest)
}

// WriteFanout writes the fan-out report to w.
func WriteFanout(w io.Writer, r *FanoutReport) {
	if r == nil || len(r.Buckets) == 0 {
		fmt.Fprintln(w, "fanout: no data")
		return
	}
	fmt.Fprintf(w, "%-12s %8s %8s %8s %10s %10s %10s\n",
		"Concurrency", "Count", "OK", "Err", "Avg(ms)", "P99(ms)", "RPS")
	for _, b := range r.Buckets {
		fmt.Fprintf(w, "%-12d %8d %8d %8d %10.2f %10.2f %10.2f\n",
			b.Concurrency, b.Count, b.Successes, b.Failures,
			b.AvgMs, b.P99Ms, b.RPS)
	}
	fmt.Fprintf(w, "\npeak RPS %.2f at concurrency %d\n", r.PeakRPS, r.Optimal)
}
