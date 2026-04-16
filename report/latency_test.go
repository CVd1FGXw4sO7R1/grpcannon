package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeLatencyResults() []Result {
	base := time.Millisecond
	return []Result{
		{Duration: 1 * base},
		{Duration: 2 * base},
		{Duration: 3 * base},
		{Duration: 4 * base},
		{Duration: 5 * base},
		{Duration: 10 * base},
	}
}

func TestLatencyBands_Empty(t *testing.T) {
	result := LatencyBands(nil, 5)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestLatencyBands_ZeroBands(t *testing.T) {
	result := LatencyBands(makeLatencyResults(), 0)
	if result != nil {
		t.Errorf("expected nil for zero bands")
	}
}

func TestLatencyBands_CountsSum(t *testing.T) {
	results := makeLatencyResults()
	bands := LatencyBands(results, 3)
	total := 0
	for _, b := range bands {
		total += b.Count
	}
	if total != len(results) {
		t.Errorf("expected total %d, got %d", len(results), total)
	}
}

func TestLatencyBands_BandCount(t *testing.T) {
	bands := LatencyBands(makeLatencyResults(), 4)
	if len(bands) != 4 {
		t.Errorf("expected 4 bands, got %d", len(bands))
	}
}

func TestWriteLatencyBands_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	err := WriteLatencyBands(&buf, makeLatencyResults(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Low") || !strings.Contains(out, "High") || !strings.Contains(out, "Count") {
		t.Errorf("missing header in output: %s", out)
	}
}

func TestWriteLatencyBands_NilResults(t *testing.T) {
	var buf bytes.Buffer
	err := WriteLatencyBands(&buf, nil, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results', got: %s", buf.String())
	}
}

func TestWriteLatencyBands_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	err := WriteLatencyBands(&buf, []Result{}, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results', got: %s", buf.String())
	}
}
