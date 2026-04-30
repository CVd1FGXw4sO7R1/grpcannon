package report

import (
	"fmt"
	"io"
	"time"
)

// PacingPoint represents a single pacing measurement window.
type PacingPoint struct {
	Window    int
	Sent      int
	Target    int
	ActualRPS float64
	TargetRPS float64
	DriftPct  float64
}

// PacingReport holds all pacing analysis windows.
type PacingReport struct {
	Points    []PacingPoint
	AvgDrift  float64
	MaxDrift  float64
	OnTarget  bool
}

// CalcPacing analyses how closely the actual request rate matched the target
// rate across equal-width time windows.
func CalcPacing(results []Result, targetRPS float64, windows int) *PacingReport {
	if len(results) == 0 || windows <= 0 || targetRPS <= 0 {
		return &PacingReport{}
	}

	sorted := make([]Result, len(results))
	copy(sorted, results)

	var earliest, latest time.Time
	for i, r := range sorted {
		if i == 0 || r.Start.Before(earliest) {
			earliest = r.Start
		}
		end := r.Start.Add(r.Duration)
		if end.After(latest) {
			latest = end
		}
	}

	total := latest.Sub(earliest)
	if total <= 0 {
		return &PacingReport{}
	}

	windowDur := total / time.Duration(windows)
	if windowDur <= 0 {
		return &PacingReport{}
	}

	counts := make([]int, windows)
	for _, r := range sorted {
		idx := int(r.Start.Sub(earliest) / windowDur)
		if idx >= windows {
			idx = windows - 1
		}
		counts[idx]++
	}

	windowSec := windowDur.Seconds()
	points := make([]PacingPoint, windows)
	var totalDrift, maxDrift float64

	for i, cnt := range counts {
		actual := float64(cnt) / windowSec
		drift := 0.0
		if targetRPS > 0 {
			drift = ((actual - targetRPS) / targetRPS) * 100.0
		}
		if drift < 0 {
			drift = -drift
		}
		if drift > maxDrift {
			maxDrift = drift
		}
		totalDrift += drift
		points[i] = PacingPoint{
			Window:    i + 1,
			Sent:      cnt,
			Target:    int(targetRPS * windowSec),
			ActualRPS: actual,
			TargetRPS: targetRPS,
			DriftPct:  drift,
		}
	}

	avgDrift := totalDrift / float64(windows)
	return &PacingReport{
		Points:   points,
		AvgDrift: avgDrift,
		MaxDrift: maxDrift,
		OnTarget: avgDrift < 5.0,
	}
}

// WritePacing writes a pacing report to w.
func WritePacing(w io.Writer, r *PacingReport) {
	if r == nil || len(r.Points) == 0 {
		fmt.Fprintln(w, "no pacing data")
		return
	}
	fmt.Fprintf(w, "%-8s %-8s %-8s %-10s %-10s %-10s\n",
		"Window", "Sent", "Target", "ActualRPS", "TargetRPS", "Drift%")
	for _, p := range r.Points {
		fmt.Fprintf(w, "%-8d %-8d %-8d %-10.2f %-10.2f %-10.2f\n",
			p.Window, p.Sent, p.Target, p.ActualRPS, p.TargetRPS, p.DriftPct)
	}
	onTarget := "NO"
	if r.OnTarget {
		onTarget = "YES"
	}
	fmt.Fprintf(w, "\nAvg Drift: %.2f%%  Max Drift: %.2f%%  On-Target: %s\n",
		r.AvgDrift, r.MaxDrift, onTarget)
}
