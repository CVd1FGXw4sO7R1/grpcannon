package report

import (
	"fmt"
	"io"
)

// RegressionThresholds defines acceptable degradation limits vs a baseline.
type RegressionThresholds struct {
	MaxP50DeltaMs  float64
	MaxP99DeltaMs  float64
	MaxErrorDelta  float64 // percentage points
}

// RegressionResult holds the outcome of a regression check.
type RegressionResult struct {
	Field    string
	Baseline float64
	Current  float64
	Delta    float64
	Passed   bool
}

// EvaluateRegression compares current report against a baseline report.
func EvaluateRegression(baseline, current *Report, t RegressionThresholds) []RegressionResult {
	if baseline == nil || current == nil {
		return nil
	}

	results := []RegressionResult{}

	p50Base := Percentile(SortedDurationsMs(baseline.Results), 50)
	p50Cur := Percentile(SortedDurationsMs(current.Results), 50)
	d50 := p50Cur - p50Base
	results = append(results, RegressionResult{
		Field:    "P50 latency (ms)",
		Baseline: p50Base,
		Current:  p50Cur,
		Delta:    d50,
		Passed:   d50 <= t.MaxP50DeltaMs,
	})

	p99Base := Percentile(SortedDurationsMs(baseline.Results), 99)
	p99Cur := Percentile(SortedDurationsMs(current.Results), 99)
	d99 := p99Cur - p99Base
	results = append(results, RegressionResult{
		Field:    "P99 latency (ms)",
		Baseline: p99Base,
		Current:  p99Cur,
		Delta:    d99,
		Passed:   d99 <= t.MaxP99DeltaMs,
	})

	errBase := baseline.ErrorRate * 100
	errCur := current.ErrorRate * 100
	dErr := errCur - errBase
	results = append(results, RegressionResult{
		Field:    "Error rate (%)",
		Baseline: errBase,
		Current:  errCur,
		Delta:    dErr,
		Passed:   dErr <= t.MaxErrorDelta,
	})

	return results
}

// WriteRegression writes regression results to w.
func WriteRegression(w io.Writer, results []RegressionResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No regression data.")
		return
	}
	fmt.Fprintf(w, "%-22s %10s %10s %10s %s\n", "Field", "Baseline", "Current", "Delta", "Status")
	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(w, "%-22s %10.2f %10.2f %+10.2f %s\n", r.Field, r.Baseline, r.Current, r.Delta, status)
	}
}
