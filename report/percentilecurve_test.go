package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeCurveResults() []Result {
	latencies := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		20 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
	}
	out := make([]Result, len(latencies))
	for i, d := range latencies {
		out[i] = Result{Duration: d}
	}
	return out
}

func TestBuildPercentileCurve_Empty(t *testing.T) {
	pts := BuildPercentileCurve(nil, 10)
	if pts != nil {
		t.Errorf("expected nil, got %v", pts)
	}
}

func TestBuildPercentileCurve_ZeroSteps(t *testing.T) {
	pts := BuildPercentileCurve(makeCurveResults(), 0)
	if pts != nil {
		t.Errorf("expected nil for zero steps")
	}
}

func TestBuildPercentileCurve_StepCount(t *testing.T) {
	pts := BuildPercentileCurve(makeCurveResults(), 10)
	if len(pts) != 11 {
		t.Errorf("expected 11 points, got %d", len(pts))
	}
}

func TestBuildPercentileCurve_Monotonic(t *testing.T) {
	pts := BuildPercentileCurve(makeCurveResults(), 20)
	for i := 1; i < len(pts); i++ {
		if pts[i].Latency < pts[i-1].Latency {
			t.Errorf("curve not monotonic at index %d", i)
		}
	}
}

func TestWritePercentileCurve_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileCurve(&buf, nil, 10)
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected no-results message")
	}
}

func TestWritePercentileCurve_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileCurve(&buf, makeCurveResults(), 4)
	out := buf.String()
	if !strings.Contains(out, "Percentile Curve") {
		t.Errorf("missing header")
	}
	if !strings.Contains(out, "100.0") {
		t.Errorf("missing 100th percentile")
	}
}
