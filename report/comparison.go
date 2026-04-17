package report

import (
	"fmt"
	"io"
)

// ComparisonReport holds two named reports for side-by-side comparison.
type ComparisonReport struct {
	BaselineLabel string
	CandidateLabel string
	Baseline  *Report
	Candidate *Report
}

// Delta returns the percentage change from baseline to candidate for a given value.
func Delta(baseline, candidate float64) float64 {
	if baseline == 0 {
		return 0
	}
	return ((candidate - baseline) / baseline) * 100
}

// WriteComparison writes a human-readable comparison table to w.
func WriteComparison(w io.Writer, c *ComparisonReport) error {
	if c == nil || c.Baseline == nil || c.Candidate == nil {
		_, err := fmt.Fprintln(w, "comparison: missing baseline or candidate")
		return err
	}

	fmt.Fprintf(w, "%-20s %12s %12s %10s\n", "Metric", c.BaselineLabel, c.CandidateLabel, "Delta")
	fmt.Fprintf(w, "%s\n", repeat("-", 58))

	writeRow := func(label string, b, cand float64, unit string) {
		d := Delta(b, cand)
		fmt.Fprintf(w, "%-20s %11.2f%s %11.2f%s %+9.2f%%\n", label, b, unit, cand, unit, d)
	}

	writeRow("P50 Latency", msFloat(c.Baseline.P50), msFloat(c.Candidate.P50), "ms")
	writeRow("P95 Latency", msFloat(c.Baseline.P95), msFloat(c.Candidate.P95), "ms")
	writeRow("P99 Latency", msFloat(c.Baseline.P99), msFloat(c.Candidate.P99), "ms")
	writeRow("Avg Latency", msFloat(c.Baseline.Avg), msFloat(c.Candidate.Avg), "ms")

	bSucc := float64(c.Baseline.Success) / float64(max1(c.Baseline.Total)) * 100
	cSucc := float64(c.Candidate.Success) / float64(max1(c.Candidate.Total)) * 100
	writeRow("Success Rate", bSucc, cSucc, "%")

	return nil
}

func max1(n int) int {
	if n < 1 {
		return 1
	}
	return n
}
