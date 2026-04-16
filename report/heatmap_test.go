package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeHeatmapResults(n int, d time.Duration) []Result {
	base := time.Now()
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		results[i] = Result{
			Start:    base.Add(time.Duration(i) * 10 * time.Millisecond),
			Duration: d,
			Err:      nil,
		}
	}
	return results
}

func TestWriteHeatmap_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	results := makeHeatmapResults(20, 8*time.Millisecond)
	err := WriteHeatmap(&buf, results, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Heatmap") {
		t.Error("expected 'Heatmap' in output")
	}
	if !strings.Contains(out, "<10ms") {
		t.Error("expected latency band label in output")
	}
}

func TestWriteHeatmap_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	err := WriteHeatmap(&buf, []Result{}, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Error("expected 'no results'")
	}
}

func TestWriteHeatmap_NilResults(t *testing.T) {
	var buf bytes.Buffer
	err := WriteHeatmap(&buf, nil, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Error("expected 'no results'")
	}
}

func TestWriteHeatmap_DefaultBuckets(t *testing.T) {
	var buf bytes.Buffer
	results := makeHeatmapResults(10, 3*time.Millisecond)
	err := WriteHeatmap(&buf, results, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestWriteHeatmap_HighLatency(t *testing.T) {
	var buf bytes.Buffer
	results := makeHeatmapResults(5, 2*time.Second)
	err := WriteHeatmap(&buf, results, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), ">=1s") {
		t.Error("expected '>=1s' band in output")
	}
}
