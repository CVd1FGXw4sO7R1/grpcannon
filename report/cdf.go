package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// CDFPoint represents a single point on the cumulative distribution function.
type CDFPoint struct {
	LatencyMs float64
	Cumulative float64 // fraction [0,1]
}

// BuildCDF constructs a CDF from the successful results using the given number
// of evenly-spaced sample points across the latency range.
func BuildCDF(results []Result, steps int) []CDFPoint {
	if len(results) == 0 || steps <= 0 {
		return nil
	}

	var durations []time.Duration
	for _, r := range results {
		if r.Err == nil {
			durations = append(durations, r.Duration)
		}
	}
	if len(durations) == 0 {
		return nil
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	minMs := float64(durations[0]) / float64(time.Millisecond)
	maxMs := float64(durations[len(durations)-1]) / float64(time.Millisecond)
	span := maxMs - minMs
	if span == 0 {
		return []CDFPoint{{LatencyMs: minMs, Cumulative: 1.0}}
	}

	total := float64(len(durations))
	points := make([]CDFPoint, 0, steps)
	for i := 0; i <= steps; i++ {
		latMs := minMs + span*float64(i)/float64(steps)
		count := 0
		for _, d := range durations {
			if float64(d)/float64(time.Millisecond) <= latMs {
				count++
			}
		}
		points = append(points, CDFPoint{
			LatencyMs:  latMs,
			Cumulative: float64(count) / total,
		})
	}
	return points
}

// WriteCDF writes a CDF table to w.
func WriteCDF(w io.Writer, results []Result, steps int) {
	points := BuildCDF(results, steps)
	if len(points) == 0 {
		fmt.Fprintln(w, "no data for CDF")
		return
	}
	fmt.Fprintf(w, "%-14s  %s\n", "Latency (ms)", "Cumulative")
	fmt.Fprintf(w, "%-14s  %s\n", "------------", "----------")
	for _, p := range points {
		fmt.Fprintf(w, "%-14.3f  %.4f\n", p.LatencyMs, p.Cumulative)
	}
}
