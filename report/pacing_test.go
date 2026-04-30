package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makePacingResults(n int, rps float64) []Result {
	results := make([]Result, n)
	interval := time.Duration(float64(time.Second) / rps)
	base := time.Now()
	for i := range results {
		results[i] = Result{
			Start:    base.Add(time.Duration(i) * interval),
			Duration: 5 * time.Millisecond,
		}
	}
	return results
}

func TestCalcPacing_Empty(t *testing.T) {
	r := CalcPacing(nil, 100, 5)
	if len(r.Points) != 0 {
		t.Errorf("expected 0 points, got %d", len(r.Points))
	}
}

func TestCalcPacing_ZeroWindows(t *testing.T) {
	results := makePacingResults(50, 10)
	r := CalcPacing(results, 10, 0)
	if len(r.Points) != 0 {
		t.Errorf("expected 0 points for zero windows")
	}
}

func TestCalcPacing_ZeroTarget(t *testing.T) {
	results := makePacingResults(50, 10)
	r := CalcPacing(results, 0, 5)
	if len(r.Points) != 0 {
		t.Errorf("expected 0 points for zero target RPS")
	}
}

func TestCalcPacing_WindowCount(t *testing.T) {
	results := makePacingResults(100, 20)
	r := CalcPacing(results, 20, 5)
	if len(r.Points) != 5 {
		t.Errorf("expected 5 windows, got %d", len(r.Points))
	}
}

func TestCalcPacing_TotalSent(t *testing.T) {
	n := 100
	results := makePacingResults(n, 20)
	r := CalcPacing(results, 20, 4)
	total := 0
	for _, p := range r.Points {
		total += p.Sent
	}
	if total != n {
		t.Errorf("expected total sent %d, got %d", n, total)
	}
}

func TestCalcPacing_DriftNonNegative(t *testing.T) {
	results := makePacingResults(60, 10)
	r := CalcPacing(results, 10, 3)
	for _, p := range r.Points {
		if p.DriftPct < 0 {
			t.Errorf("drift should be non-negative, got %.2f", p.DriftPct)
		}
	}
}

func TestCalcPacing_OnTargetUniform(t *testing.T) {
	results := makePacingResults(100, 10)
	r := CalcPacing(results, 10, 5)
	if !r.OnTarget {
		t.Errorf("expected on-target for uniform load, avg drift=%.2f", r.AvgDrift)
	}
}

func TestWritePacing_Empty(t *testing.T) {
	var buf bytes.Buffer
	WritePacing(&buf, &PacingReport{})
	if !strings.Contains(buf.String(), "no pacing data") {
		t.Errorf("expected 'no pacing data' message")
	}
}

func TestWritePacing_ValidOutput(t *testing.T) {
	results := makePacingResults(50, 10)
	r := CalcPacing(results, 10, 3)
	var buf bytes.Buffer
	WritePacing(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Window") {
		t.Errorf("expected header in output")
	}
	if !strings.Contains(out, "Drift%") {
		t.Errorf("expected Drift%% column")
	}
	if !strings.Contains(out, "On-Target") {
		t.Errorf("expected On-Target summary")
	}
}

func TestWritePacing_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WritePacing(&buf, nil)
	if !strings.Contains(buf.String(), "no pacing data") {
		t.Errorf("expected no pacing data for nil report")
	}
}
