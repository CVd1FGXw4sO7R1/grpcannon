package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeCompareReport(durations []time.Duration) *Report {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
	}
	return &Report{Results: results}
}

func TestBuildPercentileComparison_NilReports(t *testing.T) {
	out := BuildPercentileComparison(nil, nil, nil, 10)
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestBuildPercentileComparison_DefaultPercentiles(t *testing.T) {
	base := makeCompareReport([]time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond})
	cand := makeCompareReport([]time.Duration{15 * time.Millisecond, 25 * time.Millisecond, 35 * time.Millisecond})
	out := BuildPercentileComparison(base, cand, nil, 10)
	if len(out) != 5 {
		t.Errorf("expected 5 default percentiles, got %d", len(out))
	}
}

func TestBuildPercentileComparison_DeltaPositive(t *testing.T) {
	base := makeCompareReport([]time.Duration{10 * time.Millisecond, 10 * time.Millisecond, 10 * time.Millisecond})
	cand := makeCompareReport([]time.Duration{20 * time.Millisecond, 20 * time.Millisecond, 20 * time.Millisecond})
	out := BuildPercentileComparison(base, cand, []float64{50}, 10)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].DeltaMs <= 0 {
		t.Errorf("expected positive delta, got %.2f", out[0].DeltaMs)
	}
}

func TestBuildPercentileComparison_RegressionFlagged(t *testing.T) {
	base := makeCompareReport([]time.Duration{10 * time.Millisecond, 10 * time.Millisecond})
	cand := makeCompareReport([]time.Duration{25 * time.Millisecond, 25 * time.Millisecond})
	out := BuildPercentileComparison(base, cand, []float64{50}, 10)
	if !out[0].Regressed {
		t.Error("expected regression to be flagged")
	}
}

func TestBuildPercentileComparison_NoRegressionWithinThreshold(t *testing.T) {
	base := makeCompareReport([]time.Duration{10 * time.Millisecond, 10 * time.Millisecond})
	cand := makeCompareReport([]time.Duration{10500 * time.Microsecond, 10500 * time.Microsecond})
	out := BuildPercentileComparison(base, cand, []float64{50}, 10)
	if out[0].Regressed {
		t.Errorf("expected no regression, delta%% = %.2f", out[0].DeltaPct)
	}
}

func TestWritePercentileComparison_Empty(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileComparison(&buf, nil)
	if !strings.Contains(buf.String(), "no percentile") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWritePercentileComparison_ValidOutput(t *testing.T) {
	base := makeCompareReport([]time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond})
	cand := makeCompareReport([]time.Duration{12 * time.Millisecond, 22 * time.Millisecond, 50 * time.Millisecond})
	comparisons := BuildPercentileComparison(base, cand, []float64{50, 99}, 10)
	var buf bytes.Buffer
	WritePercentileComparison(&buf, comparisons)
	out := buf.String()
	if !strings.Contains(out, "Baseline") {
		t.Error("expected header 'Baseline' in output")
	}
	if !strings.Contains(out, "Candidate") {
		t.Error("expected header 'Candidate' in output")
	}
	if !strings.Contains(out, "p50") {
		t.Error("expected p50 row in output")
	}
}
