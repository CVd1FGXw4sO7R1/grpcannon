package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeStatsResults(durations []time.Duration, withError bool) []Result {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
		if withError && i == len(durations)-1 {
			results[i].Err = errors.New("rpc error")
		}
	}
	return results
}

func TestBuildPercentileStats_Empty(t *testing.T) {
	stats := BuildPercentileStats(nil)
	if stats.P50 != 0 || stats.Max != 0 {
		t.Errorf("expected zero stats for empty input, got %+v", stats)
	}
}

func TestBuildPercentileStats_AllErrors(t *testing.T) {
	results := []Result{
		{Duration: 10 * time.Millisecond, Err: errors.New("fail")},
		{Duration: 20 * time.Millisecond, Err: errors.New("fail")},
	}
	stats := BuildPercentileStats(results)
	if stats.P99 != 0 {
		t.Errorf("expected zero P99 for all-error input, got %v", stats.P99)
	}
}

func TestBuildPercentileStats_MinMax(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		200 * time.Millisecond,
	}
	results := makeStatsResults(durations, false)
	stats := BuildPercentileStats(results)

	if stats.Min != 10*time.Millisecond {
		t.Errorf("expected Min=10ms, got %v", stats.Min)
	}
	if stats.Max != 200*time.Millisecond {
		t.Errorf("expected Max=200ms, got %v", stats.Max)
	}
}

func TestBuildPercentileStats_P99Positive(t *testing.T) {
	var durations []time.Duration
	for i := 1; i <= 100; i++ {
		durations = append(durations, time.Duration(i)*time.Millisecond)
	}
	results := makeStatsResults(durations, false)
	stats := BuildPercentileStats(results)

	if stats.P99 <= 0 {
		t.Errorf("expected positive P99, got %v", stats.P99)
	}
	if stats.P99 < stats.P50 {
		t.Errorf("expected P99 >= P50, got P99=%v P50=%v", stats.P99, stats.P50)
	}
}

func TestBuildPercentileStats_SkipsErrors(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		1000 * time.Millisecond, // this one will be an error
	}
	results := makeStatsResults(durations, true)
	stats := BuildPercentileStats(results)

	// Max should be from successful results only (20ms), not the errored 1000ms
	if stats.Max > 20*time.Millisecond {
		t.Errorf("expected Max <= 20ms (errors excluded), got %v", stats.Max)
	}
}

func TestWritePercentileStats_ValidOutput(t *testing.T) {
	var durations []time.Duration
	for i := 1; i <= 100; i++ {
		durations = append(durations, time.Duration(i)*time.Millisecond)
	}
	results := makeStatsResults(durations, false)

	var buf bytes.Buffer
	WritePercentileStats(&buf, results)
	out := buf.String()

	for _, want := range []string{"p50", "p75", "p90", "p95", "p99", "p99.9", "Min", "Max"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWritePercentileStats_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileStats(&buf, nil)
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results' message, got: %s", buf.String())
	}
}

func TestWritePercentileStats_AllErrors(t *testing.T) {
	results := []Result{
		{Duration: 5 * time.Millisecond, Err: errors.New("fail")},
	}
	var buf bytes.Buffer
	WritePercentileStats(&buf, results)
	if !strings.Contains(buf.String(), "no successful") {
		t.Errorf("expected 'no successful results' message, got: %s", buf.String())
	}
}
