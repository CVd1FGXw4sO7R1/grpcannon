package report

import (
	"fmt"
	"io"
	"sort"
)

// LatencyRatioPoint holds the ratio of two percentile latencies (e.g. P99/P50)
// for a given label or bucket index.
type LatencyRatioPoint struct {
	Label    string
	Numer    float64 // numerator percentile value in ms
	Denom    float64 // denominator percentile value in ms
	Ratio    float64 // Numer / Denom, 0 if Denom == 0
}

// BuildLatencyRatio computes the ratio of two percentiles (numerP and denomP,
// expressed as 0–100) across fixed-size time buckets of results.
// buckets <= 0 defaults to 10.
func BuildLatencyRatio(results []Result, buckets int, numerP, denomP float64) []LatencyRatioPoint {
	if len(results) == 0 {
		return nil
	}
	if buckets <= 0 {
		buckets = 10
	}

	size := (len(results) + buckets - 1) / buckets
	points := make([]LatencyRatioPoint, 0, buckets)

	for i := 0; i < len(results); i += size {
		end := i + size
		if end > len(results) {
			end = len(results)
		}
		chunk := results[i:end]

		var durations []float64
		for _, r := range chunk {
			if r.Error == nil {
				durations = append(durations, float64(r.Duration.Milliseconds()))
			}
		}
		sort.Float64s(durations)

		numer := percentileFloat64(durations, numerP)
		denom := percentileFloat64(durations, denomP)

		ratio := 0.0
		if denom > 0 {
			ratio = numer / denom
		}

		points = append(points, LatencyRatioPoint{
			Label: fmt.Sprintf("bucket_%d", len(points)+1),
			Numer: numer,
			Denom: denom,
			Ratio: ratio,
		})
	}
	return points
}

func percentileFloat64(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(p/100.0*float64(len(sorted)-1)+0.5)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

// WriteLatencyRatio writes a human-readable latency ratio table to w.
func WriteLatencyRatio(w io.Writer, points []LatencyRatioPoint, numerP, denomP float64) {
	if len(points) == 0 {
		fmt.Fprintln(w, "no data")
		return
	}
	fmt.Fprintf(w, "Latency Ratio P%.0f/P%.0f\n", numerP, denomP)
	fmt.Fprintf(w, "%-14s %10s %10s %8s\n", "Bucket", fmt.Sprintf("P%.0f(ms)", numerP), fmt.Sprintf("P%.0f(ms)", denomP), "Ratio")
	for _, pt := range points {
		fmt.Fprintf(w, "%-14s %10.2f %10.2f %8.2f\n", pt.Label, pt.Numer, pt.Denom, pt.Ratio)
	}
}
