package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeStdDevResults(durations ...time.Duration) []Result {
	var results []Result
	for _, d := range durations {
		results = append(results, Result{Duration: d})
	}
	return results
}

func TestStdDev_Empty(t *testing.T) {
	if got := StdDev(nil); got != 0 {
		t.Errorf("expected 0, got %f", got)
	}
}

func TestStdDev_SingleElement(t *testing.T) {
	results := makeStdDevResults(50 * time.Millisecond)
	if got := StdDev(results); got != 0 {
		t.Errorf("expected 0 for single element, got %f", got)
	}
}

func TestStdDev_KnownValues(t *testing.T) {
	// durations: 10, 20, 30 ms => mean=20, variance=((100+0+100)/3)=66.67, std~8.165
	results := makeStdDevResults(10*time.Millisecond, 20*time.Millisecond, 30*time.Millisecond)
	got := StdDev(results)
	if got < 8.1 || got > 8.2 {
		t.Errorf("expected ~8.165, got %f", got)
	}
}

func TestStdDev_Uniform(t *testing.T) {
	results := makeStdDevResults(100*time.Millisecond, 100*time.Millisecond, 100*time.Millisecond)
	if got := StdDev(results); got != 0 {
		t.Errorf("expected 0 for uniform durations, got %f", got)
	}
}

func TestWriteStdDev_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteStdDev(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results' message, got: %s", buf.String())
	}
}

func TestWriteStdDev_ValidOutput(t *testing.T) {
	results := makeStdDevResults(10*time.Millisecond, 20*time.Millisecond, 30*time.Millisecond)
	var buf bytes.Buffer
	if err := WriteStdDev(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"StdDev Report", "Mean:", "StdDev:", "CV:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}
