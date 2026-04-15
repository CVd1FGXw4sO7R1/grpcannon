package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeResults(durations []time.Duration, failCount int) []Result {
	results := make([]Result, 0, len(durations)+failCount)
	for _, d := range durations {
		results = append(results, Result{Duration: d})
	}
	for i := 0; i < failCount; i++ {
		results = append(results, Result{Err: errDummy})
	}
	return results
}

var errDummy = &dummyErr{}

type dummyErr struct{}

func (e *dummyErr) Error() string { return "dummy error" }

func TestNew_Counts(t *testing.T) {
	results := makeResults([]time.Duration{10 * time.Millisecond, 20 * time.Millisecond}, 3)
	s := New(results, 100*time.Millisecond)
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Successes != 2 {
		t.Errorf("expected Successes=2, got %d", s.Successes)
	}
	if s.Failures != 3 {
		t.Errorf("expected Failures=3, got %d", s.Failures)
	}
}

func TestPercentile_Empty(t *testing.T) {
	s := &Summary{}
	if s.Percentile(50) != 0 {
		t.Error("expected 0 for empty durations")
	}
}

func TestPercentile_Values(t *testing.T) {
	durations := make([]time.Duration, 100)
	for i := range durations {
		durations[i] = time.Duration(i+1) * time.Millisecond
	}
	s := New(makeResults(durations, 0), time.Second)
	p50 := s.Percentile(50)
	if p50 < 49*time.Millisecond || p50 > 51*time.Millisecond {
		t.Errorf("unexpected p50: %s", p50)
	}
}

func TestPrint_ContainsFields(t *testing.T) {
	results := makeResults([]time.Duration{5 * time.Millisecond, 15 * time.Millisecond}, 1)
	s := New(results, 50*time.Millisecond)
	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()
	for _, want := range []string{"Total", "Successes", "Failures", "p50", "p90", "p99"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q", want)
		}
	}
}

func TestPrintHistogram_NoSuccesses(t *testing.T) {
	s := &Summary{}
	var buf bytes.Buffer
	s.PrintHistogram(&buf)
	if !strings.Contains(buf.String(), "no successful") {
		t.Error("expected no-data message")
	}
}

func TestPrintHistogram_Renders(t *testing.T) {
	durations := []time.Duration{1, 2, 3, 4, 5, 6, 7, 8}
	for i := range durations {
		durations[i] *= time.Millisecond
	}
	s := New(makeResults(durations, 0), 10*time.Millisecond)
	var buf bytes.Buffer
	s.PrintHistogram(&buf)
	if !strings.Contains(buf.String(), "histogram") {
		t.Error("expected histogram header")
	}
}
