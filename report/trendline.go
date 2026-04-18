package report

import (
	"fmt"
	"io"
	"time"
)

// TrendPoint represents RPS and p99 latency at a point in time.
type TrendPoint struct {
	At    time.Time
	RPS   float64
	P99Ms float64
}

// BuildTrendline buckets results into windows and computes RPS + p99 per window.
func BuildTrendline(results []Result, buckets int) []TrendPoint {
	if len(results) == 0 || buckets <= 0 {
		return nil
	}

	start := results[0].StartedAt
	end := results[len(results)-1].StartedAt
	total := end.Sub(start)
	if total <= 0 {
		return nil
	}

	width := total / time.Duration(buckets)
	type bucket struct {
		durations []time.Duration
		count     int
	}
	bins := make([]bucket, buckets)

	for _, r := range results {
		idx := int(r.StartedAt.Sub(start) / width)
		if idx >= buckets {
			idx = buckets - 1
		}
		bins[idx].count++
		if r.Err == nil {
			bins[idx].durations = append(bins[idx].durations, r.Duration)
		}
	}

	points := make([]TrendPoint, 0, buckets)
	widthSec := width.Seconds()
	for i, b := range bins {
		rps := 0.0
		if widthSec > 0 {
			rps = float64(b.count) / widthSec
		}
		p99 := Percentile(SortedDurationsMs(b.durations), 99)
		points = append(points, TrendPoint{
			At:    start.Add(time.Duration(i) * width),
			RPS:   rps,
			P99Ms: p99,
		})
	}
	return points
}

// WriteTrendline writes a trendline table to w.
func WriteTrendline(w io.Writer, results []Result, buckets int) {
	points := BuildTrendline(results, buckets)
	if len(points) == 0 {
		fmt.Fprintln(w, "no trendline data")
		return
	}
	fmt.Fprintf(w, "%-12s  %8s  %10s\n", "offset", "rps", "p99_ms")
	base := points[0].At
	for _, p := range points {
		offset := p.At.Sub(base).Round(time.Millisecond)
		fmt.Fprintf(w, "%-12s  %8.2f  %10.2f\n", offset, p.RPS, p.P99Ms)
	}
}
