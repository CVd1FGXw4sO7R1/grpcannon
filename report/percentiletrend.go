package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// PercentileTrendPoint holds percentile latencies for a time bucket.
type PercentileTrendPoint struct {
	Bucket  int
	Start   time.Time
	P50Ms   float64
	P90Ms   float64
	P95Ms   float64
	P99Ms   float64
	Count   int
}

// BuildPercentileTrend splits results into buckets and computes p50/p90/p95/p99
// for each bucket, returning a time-ordered slice of PercentileTrendPoint.
func BuildPercentileTrend(results []Result, buckets int) []PercentileTrendPoint {
	if len(results) == 0 || buckets <= 0 {
		return nil
	}

	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].StartedAt.Before(sorted[j].StartedAt)
	})

	start := sorted[0].StartedAt
	end := sorted[len(sorted)-1].StartedAt
	total := end.Sub(start)
	if total == 0 {
		total = time.Millisecond
	}
	bucketDur := total / time.Duration(buckets)
	if bucketDur == 0 {
		bucketDur = time.Millisecond
	}

	groups := make([][]time.Duration, buckets)
	starts := make([]time.Time, buckets)
	for i := 0; i < buckets; i++ {
		starts[i] = start.Add(time.Duration(i) * bucketDur)
	}

	for _, r := range sorted {
		if r.Err != nil {
			continue
		}
		idx := int(r.StartedAt.Sub(start) / bucketDur)
		if idx >= buckets {
			idx = buckets - 1
		}
		groups[idx] = append(groups[idx], r.Duration)
	}

	points := make([]PercentileTrendPoint, 0, buckets)
	for i, g := range groups {
		if len(g) == 0 {
			continue
		}
		ms := make([]float64, len(g))
		for j, d := range g {
			ms[j] = float64(d.Milliseconds())
		}
		sort.Float64s(ms)
		points = append(points, PercentileTrendPoint{
			Bucket: i,
			Start:  starts[i],
			P50Ms:  percentileFloat(ms, 50),
			P90Ms:  percentileFloat(ms, 90),
			P95Ms:  percentileFloat(ms, 95),
			P99Ms:  percentileFloat(ms, 99),
			Count:  len(g),
		})
	}
	return points
}

// WritePercentileTrend writes a formatted percentile trend table to w.
func WritePercentileTrend(w io.Writer, points []PercentileTrendPoint) {
	if len(points) == 0 {
		fmt.Fprintln(w, "no data for percentile trend")
		return
	}
	fmt.Fprintf(w, "%-8s  %-8s  %-8s  %-8s  %-8s  %s\n", "Bucket", "P50(ms)", "P90(ms)", "P95(ms)", "P99(ms)", "Count")
	for _, p := range points {
		fmt.Fprintf(w, "%-8d  %-8.2f  %-8.2f  %-8.2f  %-8.2f  %d\n",
			p.Bucket, p.P50Ms, p.P90Ms, p.P95Ms, p.P99Ms, p.Count)
	}
}

func percentileFloat(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p / 100.0)
	return sorted[idx]
}
