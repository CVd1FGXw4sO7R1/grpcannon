package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeDriftResults(n int, base time.Duration, step time.Duration) []Result {
	results := make([]Result, n)
	for i := range results {
		results[i] = Result{Duration: base + time.Duration(i)*step}
	}
	return results
}

func TestCalcDrift_Empty(t *testing.T) {
	if CalcDrift(nil, 5) != nil {
		t.Fatal("expected nil for empty input")
	}
}

func TestCalcDrift_ZeroBuckets(t *testing.T) {
	results := makeDriftResults(10, 10*time.Millisecond, time.Millisecond)
	if CalcDrift(results, 0) != nil {
		t.Fatal("expected nil for zero buckets")
	}
}

func TestCalcDrift_BucketCount(t *testing.T) {
	results := makeDriftResults(20, 10*time.Millisecond, time.Millisecond)
	points := CalcDrift(results, 4)
	if len(points) != 4 {
		t.Fatalf("expected 4 points, got %d", len(points))
	}
}

func TestCalcDrift_BaselineZeroDelta(t *testing.T) {
	results := makeDriftResults(10, 10*time.Millisecond, 0)
	points := CalcDrift(results, 2)
	for _, p := range points {
		if p.DeltaMs != 0 {
			t.Fatalf("expected zero delta for uniform latency, got %.2f", p.DeltaMs)
		}
	}
}

func TestCalcDrift_DriftIncreases(t *testing.T) {
	// 20 results: first 10 at 10ms, last 10 at 20ms
	var results []Result
	for i := 0; i < 10; i++ {
		results = append(results, Result{Duration: 10 * time.Millisecond})
	}
	for i := 0; i < 10; i++ {
		results = append(results, Result{Duration: 20 * time.Millisecond})
	}
	points := CalcDrift(results, 2)
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
	if points[1].DeltaMs <= 0 {
		t.Fatalf("expected positive drift in second window, got %.2f", points[1].DeltaMs)
	}
}

func TestWriteDrift_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteDrift(&buf, nil, 5)
	if !strings.Contains(buf.String(), "no data") {
		t.Fatal("expected 'no data' message")
	}
}

func TestWriteDrift_ValidOutput(t *testing.T) {
	results := makeDriftResults(20, 10*time.Millisecond, time.Millisecond)
	var buf bytes.Buffer
	WriteDrift(&buf, results, 4)
	out := buf.String()
	if !strings.Contains(out, "Drift") {
		t.Fatal("expected header in output")
	}
	if !strings.Contains(out, "Window") {
		t.Fatal("expected column header")
	}
}
