package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeP99Results(n int, base time.Time, spread time.Duration) []Result {
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		results[i] = Result{
			StartedAt: base.Add(time.Duration(i) * spread),
			Duration:  time.Duration(10+i) * time.Millisecond,
		}
	}
	return results
}

func TestBuildP99Trend_Empty(t *testing.T) {
	pts := BuildP99Trend(nil, 5)
	if pts != nil {
		t.Errorf("expected nil, got %v", pts)
	}
}

func TestBuildP99Trend_ZeroBuckets(t *testing.T) {
	results := makeP99Results(10, time.Now(), 10*time.Millisecond)
	pts := BuildP99Trend(results, 0)
	if pts != nil {
		t.Errorf("expected nil for zero buckets")
	}
}

func TestBuildP99Trend_BucketCount(t *testing.T) {
	results := makeP99Results(20, time.Now(), 50*time.Millisecond)
	pts := BuildP99Trend(results, 4)
	if len(pts) != 4 {
		t.Errorf("expected 4 buckets, got %d", len(pts))
	}
}

func TestBuildP99Trend_P99Positive(t *testing.T) {
	base := time.Now()
	results := makeP99Results(20, base, 50*time.Millisecond)
	pts := BuildP99Trend(results, 4)
	for _, p := range pts {
		if p.Count > 0 && p.P99Ms < 0 {
			t.Errorf("bucket %d: p99 should be non-negative, got %f", p.Bucket, p.P99Ms)
		}
	}
}

func TestBuildP99Trend_CountsSum(t *testing.T) {
	results := makeP99Results(20, time.Now(), 50*time.Millisecond)
	pts := BuildP99Trend(results, 4)
	total := 0
	for _, p := range pts {
		total += p.Count
	}
	if total != 20 {
		t.Errorf("expected total count 20, got %d", total)
	}
}

func TestWriteP99Trend_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteP99Trend(&buf, nil)
	if !strings.Contains(buf.String(), "no p99 trend") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestWriteP99Trend_ValidOutput(t *testing.T) {
	results := makeP99Results(12, time.Now(), 50*time.Millisecond)
	pts := BuildP99Trend(results, 3)
	var buf bytes.Buffer
	WriteP99Trend(&buf, pts)
	out := buf.String()
	if !strings.Contains(out, "Bucket") {
		t.Errorf("expected header in output, got %q", out)
	}
	if !strings.Contains(out, "P99") {
		t.Errorf("expected P99 column in output")
	}
}
