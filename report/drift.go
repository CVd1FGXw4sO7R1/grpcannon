package report

import (
	"fmt"
	"io"
	"time"
)

// DriftPoint represents latency drift over a rolling window.
type DriftPoint struct {
	Window   int
	AvgMs    float64
	DeltaMs  float64
	DriftPct float64
}

// CalcDrift computes per-window average latency and drift relative to the first window.
func CalcDrift(results []Result, buckets int) []DriftPoint {
	if len(results) == 0 || buckets <= 0 {
		return nil
	}

	size := len(results) / buckets
	if size == 0 {
		size = 1
	}

	var points []DriftPoint
	var baseline float64

	for i := 0; i < buckets; i++ {
		start := i * size
		end := start + size
		if end > len(results) {
			end = len(results)
		}
		if start >= len(results) {
			break
		}

		slice := results[start:end]
		var sum time.Duration
		var count int
		for _, r := range slice {
			if r.Err == nil {
				sum += r.Duration
				count++
			}
		}
		if count == 0 {
			continue
		}

		avg := float64(sum.Milliseconds()) / float64(count)
		if len(points) == 0 {
			baseline = avg
		}

		delta := avg - baseline
		var pct float64
		if baseline != 0 {
			pct = (delta / baseline) * 100
		}

		points = append(points, DriftPoint{
			Window:   i + 1,
			AvgMs:    avg,
			DeltaMs:  delta,
			DriftPct: pct,
		})
	}
	return points
}

// WriteDrift writes latency drift analysis to w.
func WriteDrift(w io.Writer, results []Result, buckets int) {
	points := CalcDrift(results, buckets)
	if len(points) == 0 {
		fmt.Fprintln(w, "Drift: no data")
		return
	}
	fmt.Fprintln(w, "Latency Drift Analysis:")
	fmt.Fprintf(w, "  %-8s %-10s %-10s %s\n", "Window", "Avg(ms)", "Delta(ms)", "Drift%")
	for _, p := range points {
		fmt.Fprintf(w, "  %-8d %-10.2f %-10.2f %.1f%%\n", p.Window, p.AvgMs, p.DeltaMs, p.DriftPct)
	}
}
