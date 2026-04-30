package report

import (
	"fmt"
	"io"
	"sort"
)

// FencePoint represents a latency fence at a given percentile.
type FencePoint struct {
	Percentile float64
	LatencyMs  float64
	Breached   bool
}

// PercentileFenceReport holds fence evaluation results.
type PercentileFenceReport struct {
	Fences      []FencePoint
	TotalFences int
	Breached    int
	Passed      int
}

// BuildPercentileFence evaluates whether latency percentiles exceed given fence thresholds.
// thresholds maps percentile (e.g. 50, 90, 99) to max allowed latency in ms.
func BuildPercentileFence(r *Report, thresholds map[float64]float64) *PercentileFenceReport {
	if r == nil || len(r.Results) == 0 || len(thresholds) == 0 {
		return &PercentileFenceReport{}
	}

	var durations []float64
	for _, res := range r.Results {
		if res.Error == nil {
			durations = append(durations, float64(res.Duration.Milliseconds()))
		}
	}
	if len(durations) == 0 {
		return &PercentileFenceReport{}
	}
	sort.Float64s(durations)

	// Collect and sort percentile keys for deterministic output.
	keys := make([]float64, 0, len(thresholds))
	for p := range thresholds {
		keys = append(keys, p)
	}
	sort.Float64s(keys)

	var fences []FencePoint
	breached := 0
	for _, p := range keys {
		max := thresholds[p]
		latency := percentileFloat64(durations, p)
		b := latency > max
		if b {
			breached++
		}
		fences = append(fences, FencePoint{
			Percentile: p,
			LatencyMs:  latency,
			Breached:   b,
		})
	}

	return &PercentileFenceReport{
		Fences:      fences,
		TotalFences: len(fences),
		Breached:    breached,
		Passed:      len(fences) - breached,
	}
}

// WritePercentileFence writes the fence report to w.
func WritePercentileFence(w io.Writer, fr *PercentileFenceReport) {
	if fr == nil || len(fr.Fences) == 0 {
		fmt.Fprintln(w, "percentile fence: no data")
		return
	}
	fmt.Fprintf(w, "Percentile Fence Report (%d fences, %d breached, %d passed)\n",
		fr.TotalFences, fr.Breached, fr.Passed)
	fmt.Fprintf(w, "%-12s %-14s %-8s\n", "Percentile", "Latency(ms)", "Status")
	for _, f := range fr.Fences {
		status := "OK"
		if f.Breached {
			status = "BREACHED"
		}
		fmt.Fprintf(w, "P%-11.0f %-14.2f %s\n", f.Percentile, f.LatencyMs, status)
	}
}
