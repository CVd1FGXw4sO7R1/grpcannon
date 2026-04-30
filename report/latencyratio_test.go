package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeRatioResults(n int, baseDur time.Duration) []Result {
	res := make([]Result, n)
	for i := 0; i < n; i++ {
		res[i] = Result{Duration: baseDur + time.Duration(i)*time.Millisecond}
	}
	return res
}

func TestBuildLatencyRatio_Empty(t *testing.T) {
	pts := BuildLatencyRatio(nil, 5, 99, 50)
	if pts != nil {
		t.Errorf("expected nil, got %v", pts)
	}
}

func TestBuildLatencyRatio_ZeroBuckets(t *testing.T) {
	results := makeRatioResults(20, 10*time.Millisecond)
	pts := BuildLatencyRatio(results, 0, 99, 50)
	if len(pts) == 0 {
		t.Error("expected non-empty points with zero buckets defaulting to 10")
	}
}

func TestBuildLatencyRatio_BucketCount(t *testing.T) {
	results := makeRatioResults(50, 5*time.Millisecond)
	pts := BuildLatencyRatio(results, 5, 99, 50)
	if len(pts) != 5 {
		t.Errorf("expected 5 buckets, got %d", len(pts))
	}
}

func TestBuildLatencyRatio_RatioPositive(t *testing.T) {
	results := makeRatioResults(100, 1*time.Millisecond)
	pts := BuildLatencyRatio(results, 10, 99, 50)
	for _, pt := range pts {
		if pt.Ratio < 1.0 {
			t.Errorf("expected P99/P50 ratio >= 1, got %.4f for %s", pt.Ratio, pt.Label)
		}
	}
}

func TestBuildLatencyRatio_ZeroDenominator(t *testing.T) {
	// All results have same duration => P50 == P99, ratio should be 1
	// But if denom is zero (zero-duration errors), ratio should be 0
	results := []Result{
		{Duration: 0, Error: errors.New("fail")},
		{Duration: 0, Error: errors.New("fail")},
	}
	pts := BuildLatencyRatio(results, 1, 99, 50)
	if len(pts) != 1 {
		t.Fatalf("expected 1 point, got %d", len(pts))
	}
	if pts[0].Ratio != 0 {
		t.Errorf("expected ratio 0 when denom is 0, got %.4f", pts[0].Ratio)
	}
}

func TestBuildLatencyRatio_SkipsErrors(t *testing.T) {
	results := []Result{
		{Duration: 100 * time.Millisecond},
		{Duration: 200 * time.Millisecond, Error: errors.New("fail")},
		{Duration: 300 * time.Millisecond},
	}
	pts := BuildLatencyRatio(results, 1, 99, 50)
	if len(pts) != 1 {
		t.Fatalf("expected 1 point, got %d", len(pts))
	}
	// Only 100ms and 300ms are counted; P50=100, P99=300 => ratio=3
	if pts[0].Denom == 0 {
		t.Error("denom should not be zero with valid results present")
	}
}

func TestWriteLatencyRatio_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteLatencyRatio(&buf, nil, 99, 50)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data', got: %s", buf.String())
	}
}

func TestWriteLatencyRatio_ValidOutput(t *testing.T) {
	results := makeRatioResults(30, 10*time.Millisecond)
	pts := BuildLatencyRatio(results, 3, 99, 50)
	var buf bytes.Buffer
	WriteLatencyRatio(&buf, pts, 99, 50)
	out := buf.String()
	if !strings.Contains(out, "Ratio") {
		t.Errorf("expected header 'Ratio' in output, got: %s", out)
	}
	if !strings.Contains(out, "bucket_1") {
		t.Errorf("expected 'bucket_1' in output, got: %s", out)
	}
}
