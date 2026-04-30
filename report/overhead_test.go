package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeOverheadResults(durations []time.Duration, errCount int) []Result {
	results := make([]Result, 0, len(durations)+errCount)
	for _, d := range durations {
		results = append(results, Result{Duration: d})
	}
	for i := 0; i < errCount; i++ {
		results = append(results, Result{Duration: time.Millisecond, Err: errSentinel})
	}
	return results
}

var errSentinel = fmt.Errorf("sentinel error")

func TestCalcOverhead_Empty(t *testing.T) {
	stats := CalcOverhead(nil)
	if stats.TotalRequests != 0 {
		t.Errorf("expected 0 total, got %d", stats.TotalRequests)
	}
}

func TestCalcOverhead_AllErrors(t *testing.T) {
	results := makeOverheadResults(nil, 5)
	stats := CalcOverhead(results)
	if stats.TotalRequests != 5 {
		t.Errorf("expected 5 total, got %d", stats.TotalRequests)
	}
	if stats.SuccessCount != 0 {
		t.Errorf("expected 0 successes, got %d", stats.SuccessCount)
	}
	if stats.AvgLatencyMs != 0 {
		t.Errorf("expected 0 avg, got %f", stats.AvgLatencyMs)
	}
}

func TestCalcOverhead_KnownValues(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}
	results := makeOverheadResults(durations, 0)
	stats := CalcOverhead(results)

	if stats.TotalRequests != 5 {
		t.Errorf("expected 5 total, got %d", stats.TotalRequests)
	}
	if stats.SuccessCount != 5 {
		t.Errorf("expected 5 successes, got %d", stats.SuccessCount)
	}
	if stats.MinLatencyMs != 10.0 {
		t.Errorf("expected min 10, got %f", stats.MinLatencyMs)
	}
	if stats.MaxLatencyMs != 50.0 {
		t.Errorf("expected max 50, got %f", stats.MaxLatencyMs)
	}
	if stats.AvgLatencyMs != 30.0 {
		t.Errorf("expected avg 30, got %f", stats.AvgLatencyMs)
	}
}

func TestCalcOverhead_OverheadIsAvgMinusP50(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		200 * time.Millisecond, // tail outlier pushes avg up
	}
	results := makeOverheadResults(durations, 0)
	stats := CalcOverhead(results)

	expected := stats.AvgLatencyMs - stats.P50Ms
	if stats.OverheadMs != expected {
		t.Errorf("expected overhead %f, got %f", expected, stats.OverheadMs)
	}
}

func TestWriteOverhead_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteOverhead(&buf, OverheadStats{})
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results', got: %s", buf.String())
	}
}

func TestWriteOverhead_ValidOutput(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}
	results := makeOverheadResults(durations, 1)
	stats := CalcOverhead(results)

	var buf bytes.Buffer
	WriteOverhead(&buf, stats)
	out := buf.String()

	for _, want := range []string{"Overhead Report", "Total Requests", "Successes", "P99", "Tail Overhead"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}
