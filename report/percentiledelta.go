package report

import (
	"fmt"
	"io"
)

// PercentileDelta holds the difference in a percentile latency between two reports.
type PercentileDelta struct {
	Percentile float64
	BaselineMs float64
	CurrentMs  float64
	DeltaMs    float64
	DeltaPct   float64
}

// BuildPercentileDeltas compares percentile latencies between a baseline and current report.
func BuildPercentileDeltas(baseline, current *Report, percentiles []float64) []PercentileDelta {
	if baseline == nil || current == nil {
		return nil
	}
	if len(percentiles) == 0 {
		percentiles = []float64{50, 75, 90, 95, 99}
	}

	baseDurs := SortedDurationsMs(baseline.Results)
	curDurs := SortedDurationsMs(current.Results)

	deltas := make([]PercentileDelta, 0, len(percentiles))
	for _, p := range percentiles {
		baseVal := Percentile(baseDurs, p)
		curVal := Percentile(curDurs, p)
		deltaMs := curVal - baseVal
		var deltaPct float64
		if baseVal != 0 {
			deltaPct = (deltaMs / baseVal) * 100
		}
		deltas = append(deltas, PercentileDelta{
			Percentile: p,
			BaselineMs: baseVal,
			CurrentMs:  curVal,
			DeltaMs:    deltaMs,
			DeltaPct:   deltaPct,
		})
	}
	return deltas
}

// WritePercentileDeltas writes a human-readable percentile delta table to w.
func WritePercentileDeltas(w io.Writer, deltas []PercentileDelta) {
	if len(deltas) == 0 {
		fmt.Fprintln(w, "no percentile delta data available")
		return
	}
	fmt.Fprintf(w, "%-12s %12s %12s %12s %10s\n", "Percentile", "Baseline(ms)", "Current(ms)", "Delta(ms)", "Delta(%)")
	fmt.Fprintf(w, "%s\n", "------------------------------------------------------------")
	for _, d := range deltas {
		sign := ""
		if d.DeltaMs > 0 {
			sign = "+"
		}
		fmt.Fprintf(w, "p%-11.0f %12.2f %12.2f %s%11.2f %s%9.1f%%\n",
			d.Percentile, d.BaselineMs, d.CurrentMs, sign, d.DeltaMs, sign, d.DeltaPct)
	}
}
