package report

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// PercentileCurvePoint holds a percentile and its corresponding latency.
type PercentileCurvePoint struct {
	Percentile float64
	Latency    time.Duration
}

// BuildPercentileCurve returns latency values at evenly-spaced percentile
// steps from 0 to 100.
func BuildPercentileCurve(results []Result, steps int) []PercentileCurvePoint {
	if len(results) == 0 || steps <= 0 {
		return nil
	}

	durations := make([]float64, 0, len(results))
	for _, r := range results {
		durations = append(durations, float64(r.Duration.Milliseconds()))
	}
	sort.Float64s(durations)

	points := make([]PercentileCurvePoint, 0, steps+1)
	for i := 0; i <= steps; i++ {
		p := float64(i) * 100.0 / float64(steps)
		v := Percentile(durations, p)
		points = append(points, PercentileCurvePoint{
			Percentile: p,
			Latency:    time.Duration(v) * time.Millisecond,
		})
	}
	return points
}

// WritePercentileCurve writes a percentile-latency curve to w.
func WritePercentileCurve(w io.Writer, results []Result, steps int) {
	points := BuildPercentileCurve(results, steps)
	if len(points) == 0 {
		fmt.Fprintln(w, "No results to display.")
		return
	}
	fmt.Fprintln(w, "Percentile Curve:")
	fmt.Fprintf(w, "  %-10s %s\n", "Percentile", "Latency")
	for _, pt := range points {
		fmt.Fprintf(w, "  %-10.1f %s\n", pt.Percentile, roundDuration(pt.Latency))
	}
}
