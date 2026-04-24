package report

import (
	"fmt"
	"io"
)

// PercentileComparison holds the comparison of a single percentile between two reports.
type PercentileComparison struct {
	Percentile float64
	BaselineMs float64
	CandidateMs float64
	DeltaMs    float64
	DeltaPct   float64
	Regressed  bool
}

// BuildPercentileComparison compares latency percentiles between a baseline and candidate report.
// threshold is the maximum allowed percentage increase before marking as regressed.
func BuildPercentileComparison(baseline, candidate *Report, percentiles []float64, thresholdPct float64) []PercentileComparison {
	if baseline == nil || candidate == nil {
		return nil
	}
	if len(percentiles) == 0 {
		percentiles = []float64{50, 75, 90, 95, 99}
	}

	baseDurs := SortedDurationsMs(baseline.Results)
	candDurs := SortedDurationsMs(candidate.Results)

	results := make([]PercentileComparison, 0, len(percentiles))
	for _, p := range percentiles {
		baseVal := Percentile(baseDurs, p)
		candVal := Percentile(candDurs, p)
		delta := candVal - baseVal
		var deltaPct float64
		if baseVal > 0 {
			deltaPct = (delta / baseVal) * 100
		}
		results = append(results, PercentileComparison{
			Percentile:  p,
			BaselineMs:  baseVal,
			CandidateMs: candVal,
			DeltaMs:     delta,
			DeltaPct:    deltaPct,
			Regressed:   deltaPct > thresholdPct,
		})
	}
	return results
}

// WritePercentileComparison writes a formatted percentile comparison table to w.
func WritePercentileComparison(w io.Writer, comparisons []PercentileComparison) {
	if len(comparisons) == 0 {
		fmt.Fprintln(w, "no percentile comparison data")
		return
	}
	fmt.Fprintf(w, "%-8s  %10s  %10s  %10s  %8s  %s\n",
		"P", "Baseline", "Candidate", "Delta", "Delta%", "Status")
	fmt.Fprintf(w, "%s\n", "--------  ----------  ----------  ----------  --------  ------")
	for _, c := range comparisons {
		status := "OK"
		if c.Regressed {
			status = "REGRESSED"
		}
		fmt.Fprintf(w, "p%-7.0f  %9.2fms  %9.2fms  %+9.2fms  %+7.1f%%  %s\n",
			c.Percentile, c.BaselineMs, c.CandidateMs, c.DeltaMs, c.DeltaPct, status)
	}
}
