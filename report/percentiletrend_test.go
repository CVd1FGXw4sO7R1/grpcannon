package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makePTrendResults(n int, base time.Duration, step time.Duration) []Result {
	results := make([]Result, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		results[i] = Result{
			StartedAt: now.Add(time.Duration(i) * time.Millisecond * 10),
			Duration:  base + time.Duration(i)*step,
		}
	}
	return results
}

func TestBuildPercentileTrend_Empty(t *testing.T) {
	points := BuildPercentileTrend(nil, 5)
	if len(points) != 0 {
		t.Fatalf("expected empty, got %d", len(points))
	}
}

func TestBuildPercentileTrend_ZeroBuckets(t *testing.T) {
	results := makePTrendResults(10, 10*time.Millisecond, time.Millisecond)
	points := BuildPercentileTrend(results, 0)
	if len(points) != 0 {
		t.Fatalf("expected empty for zero buckets, got %d", len(points))
	}
}

func TestBuildPercentileTrend_BucketCount(t *testing.T) {
	results := makePTrendResults(20, 5*time.Millisecond, time.Millisecond)
	points := BuildPercentileTrend(results, 4)
	if len(points) == 0 {
		t.Fatal("expected non-empty points")
	}
	if len(points) > 4 {
		t.Fatalf("expected at most 4 buckets, got %d", len(points))
	}
}

func TestBuildPercentileTrend_P99Positive(t *testing.T) {
	results := makePTrendResults(20, 5*time.Millisecond, time.Millisecond)
	points := BuildPercentileTrend(results, 4)
	for _, p := range points {
		if p.P99Ms <= 0 {
			t.Errorf("bucket %d: expected P99 > 0, got %.2f", p.Bucket, p.P99Ms)
		}
	}
}

func TestBuildPercentileTrend_SkipsErrors(t *testing.T) {
	results := makePTrendResults(10, 5*time.Millisecond, time.Millisecond)
	for i := range results {
		if i%2 == 0 {
			results[i].Err = errTest
		}
	}
	points := BuildPercentileTrend(results, 3)
	for _, p := range points {
		if p.P50Ms < 0 {
			t.Errorf("unexpected negative P50")
		}
	}
}

func TestBuildPercentileTrend_P50LteP99(t *testing.T) {
	results := makePTrendResults(30, 2*time.Millisecond, 500*time.Microsecond)
	points := BuildPercentileTrend(results, 5)
	for _, p := range points {
		if p.P50Ms > p.P99Ms+0.001 {
			t.Errorf("bucket %d: P50 %.2f > P99 %.2f", p.Bucket, p.P50Ms, p.P99Ms)
		}
	}
}

func TestWritePercentileTrend_ValidOutput(t *testing.T) {
	results := makePTrendResults(20, 5*time.Millisecond, time.Millisecond)
	points := BuildPercentileTrend(results, 4)
	var buf bytes.Buffer
	WritePercentileTrend(&buf, points)
	out := buf.String()
	if !strings.Contains(out, "P50") {
		t.Error("expected P50 header")
	}
	if !strings.Contains(out, "P99") {
		t.Error("expected P99 header")
	}
}

func TestWritePercentileTrend_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileTrend(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Error("expected 'no data' message")
	}
}
