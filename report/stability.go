package report

import (
	"fmt"
	"io"
	"math"
	"time"
)

// StabilityWindow holds metrics for a single time window.
type StabilityWindow struct {
	Start     time.Time
	End       time.Time
	P99Ms     float64
	ErrorRate float64
	RPS       float64
}

// StabilityReport summarizes how stable the system was over the run.
type StabilityReport struct {
	Windows        []StabilityWindow
	MaxP99JitterMs float64
	MaxErrorDelta  float64
	Stable         bool
	Reason         string
}

// BuildStability partitions results into windows and evaluates stability.
// A run is considered stable when p99 jitter < p99JitterThresholdMs and
// error-rate delta across windows < errorDeltaThreshold.
func BuildStability(results []Result, windows int, p99JitterThresholdMs, errorDeltaThreshold float64) *StabilityReport {
	if len(results) == 0 || windows <= 0 {
		return &StabilityReport{Stable: true, Reason: "no data"}
	}

	buckets := Bucketize(results, windows)
	win := make([]StabilityWindow, 0, len(buckets))

	for _, b := range buckets {
		if len(b.Results) == 0 {
			continue
		}
		total := len(b.Results)
		failed := 0
		for _, r := range b.Results {
			if r.Error != nil {
				failed++
			}
		}
		p99 := Percentile(b.Results, 99)
		win = append(win, StabilityWindow{
			Start:     b.Start,
			End:       b.End,
			P99Ms:     float64(p99) / float64(time.Millisecond),
			ErrorRate: float64(failed) / float64(total),
			RPS:       b.RPS,
		})
	}

	if len(win) == 0 {
		return &StabilityReport{Stable: true, Reason: "no data"}
	}

	minP99, maxP99 := win[0].P99Ms, win[0].P99Ms
	minErr, maxErr := win[0].ErrorRate, win[0].ErrorRate
	for _, w := range win[1:] {
		minP99 = math.Min(minP99, w.P99Ms)
		maxP99 = math.Max(maxP99, w.P99Ms)
		minErr = math.Min(minErr, w.ErrorRate)
		maxErr = math.Max(maxErr, w.ErrorRate)
	}

	jitter := maxP99 - minP99
	errDelta := maxErr - minErr

	stable := true
	reason := "ok"
	if jitter > p99JitterThresholdMs {
		stable = false
		reason = fmt.Sprintf("p99 jitter %.2fms exceeds threshold %.2fms", jitter, p99JitterThresholdMs)
	} else if errDelta > errorDeltaThreshold {
		stable = false
		reason = fmt.Sprintf("error-rate delta %.4f exceeds threshold %.4f", errDelta, errorDeltaThreshold)
	}

	return &StabilityReport{
		Windows:        win,
		MaxP99JitterMs: jitter,
		MaxErrorDelta:  errDelta,
		Stable:         stable,
		Reason:         reason,
	}
}

// WriteStability writes a human-readable stability summary to w.
func WriteStability(w io.Writer, sr *StabilityReport) {
	if sr == nil {
		fmt.Fprintln(w, "stability: no report")
		return
	}
	status := "STABLE"
	if !sr.Stable {
		status = "UNSTABLE"
	}
	fmt.Fprintf(w, "Stability: %s (%s)\n", status, sr.Reason)
	fmt.Fprintf(w, "  Max P99 Jitter : %.2f ms\n", sr.MaxP99JitterMs)
	fmt.Fprintf(w, "  Max Error Delta: %.4f\n", sr.MaxErrorDelta)
	fmt.Fprintf(w, "  Windows        : %d\n", len(sr.Windows))
	for i, win := range sr.Windows {
		fmt.Fprintf(w, "  [%2d] p99=%.2fms err=%.2f%% rps=%.1f\n",
			i+1, win.P99Ms, win.ErrorRate*100, win.RPS)
	}
}
