package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makePercentileDeltaReport(durations []time.Duration) *Report {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
	}
	return &Report{Results: results}
}

func TestBuildPercentileDeltas_NilReports(t *testing.T) {
	deltas := BuildPercentileDeltas(nil, nil, nil)
	if deltas != nil {
		t.Errorf("expected nil, got %v", deltas)
	}
}

func TestBuildPercentileDeltas_DefaultPercentiles(t *testing.T) {
	base := makePercentileDeltaReport([]time.Duration{
		10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond,
		40 * time.Millisecond, 50 * time.Millisecond,
	})
	cur := makePercentileDeltaReport([]time.Duration{
		20 * time.Millisecond, 30 * time.Millisecond, 40 * time.Millisecond,
		50 * time.Millisecond, 60 * time.Millisecond,
	})
	deltas := BuildPercentileDeltas(base, cur, nil)
	if len(deltas) != 5 {
		t.Fatalf("expected 5 deltas, got %d", len(deltas))
	}
}

func TestBuildPercentileDeltas_DeltaPositive(t *testing.T) {
	base := makePercentileDeltaReport([]time.Duration{
		10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond,
	})
	cur := makePercentileDeltaReport([]time.Duration{
		20 * time.Millisecond, 40 * time.Millisecond, 60 * time.Millisecond,
	})
	deltas := BuildPercentileDeltas(base, cur, []float64{50})
	if deltas[0].DeltaMs <= 0 {
		t.Errorf("expected positive delta, got %f", deltas[0].DeltaMs)
	}
}

func TestBuildPercentileDeltas_DeltaPctZeroBaseline(t *testing.T) {
	base := makePercentileDeltaReport([]time.Duration{})
	cur := makePercentileDeltaReport([]time.Duration{10 * time.Millisecond})
	deltas := BuildPercentileDeltas(base, cur, []float64{50})
	if deltas[0].DeltaPct != 0 {
		t.Errorf("expected 0 pct when baseline is zero, got %f", deltas[0].DeltaPct)
	}
}

func TestWritePercentileDeltas_Empty(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileDeltas(&buf, nil)
	if !strings.Contains(buf.String(), "no percentile delta") {
		t.Errorf("expected no-data message, got: %s", buf.String())
	}
}

func TestWritePercentileDeltas_ValidOutput(t *testing.T) {
	base := makePercentileDeltaReport([]time.Duration{
		10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond,
	})
	cur := makePercentileDeltaReport([]time.Duration{
		15 * time.Millisecond, 25 * time.Millisecond, 35 * time.Millisecond,
	})
	deltas := BuildPercentileDeltas(base, cur, []float64{50, 99})
	var buf bytes.Buffer
	WritePercentileDeltas(&buf, deltas)
	out := buf.String()
	if !strings.Contains(out, "Percentile") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "p50") {
		t.Errorf("expected p50 row in output, got: %s", out)
	}
}
