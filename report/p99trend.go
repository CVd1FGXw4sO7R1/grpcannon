package report

import (
	"fmt"
	"io"
	"time"
)

// P99TrendPoint holds the p99 latency for a time bucket.
type P99TrendPoint struct {
	Bucket    int
	Start     time.Time
	P99Ms     float64
	Count     int
}

// BuildP99Trend divides results into buckets and computes p99 latency per bucket.
func BuildP99Trend(results []Result, buckets int) []P99TrendPoint {
	if len(results) == 0 || buckets <= 0 {
		return nil
	}

	// find time range
	start := results[0].StartedAt
	end := results[0].StartedAt
	for _, r := range results {
		if r.StartedAt.Before(start) {
			start = r.StartedAt
		}
		if r.StartedAt.After(end) {
			end = r.StartedAt
		}
	}

	total := end.Sub(start)
	if total <= 0 {
		total = time.Millisecond
	}
	bucketDur := total / time.Duration(buckets)
	if bucketDur <= 0 {
		bucketDur = time.Millisecond
	}

	groups := make([][]Result, buckets)
	for _, r := range results {
		idx := int(r.StartedAt.Sub(start) / bucketDur)
		if idx >= buckets {
			idx = buckets - 1
		}
		groups[idx] = append(groups[idx], r)
	}

	points := make([]P99TrendPoint, 0, buckets)
	for i, g := range groups {
		var durations []Result
		for _, r := range g {
			if r.Error == nil {
				durations = append(durations, r)
			}
		}
		p99 := 0.0
		if len(durations) > 0 {
			ds := make([]Result, len(durations))
			copy(ds, durations)
			p99 = Percentile(SortedDurationsMs(resultDurations(ds)), 99)
		}
		points = append(points, P99TrendPoint{
			Bucket: i,
			Start:  start.Add(time.Duration(i) * bucketDur),
			P99Ms:  p99,
			Count:  len(g),
		})
	}
	return points
}

func resultDurations(results []Result) []Result {
	return results
}

// WriteP99Trend writes the p99 latency trend to w.
func WriteP99Trend(w io.Writer, points []P99TrendPoint) {
	if len(points) == 0 {
		fmt.Fprintln(w, "no p99 trend data")
		return
	}
	fmt.Fprintf(w, "%-8s %-12s %10s %8s\n", "Bucket", "Start", "P99(ms)", "Count")
	for _, p := range points {
		fmt.Fprintf(w, "%-8d %-12s %10.2f %8d\n", p.Bucket, p.Start.Format("15:04:05"), p.P99Ms, p.Count)
	}
}
