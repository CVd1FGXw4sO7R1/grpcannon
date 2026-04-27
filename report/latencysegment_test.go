package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeSegmentResults(n int, dur time.Duration) []Result {
	base := time.Now()
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		results[i] = Result{
			StartedAt: base.Add(time.Duration(i) * 10 * time.Millisecond),
			Duration:  dur,
		}
	}
	return results
}

func TestBuildLatencySegments_Empty(t *testing.T) {
	segs := BuildLatencySegments(nil, 4)
	if segs != nil {
		t.Fatalf("expected nil, got %v", segs)
	}
}

func TestBuildLatencySegments_ZeroSegments(t *testing.T) {
	results := makeSegmentResults(10, 5*time.Millisecond)
	segs := BuildLatencySegments(results, 0)
	if segs != nil {
		t.Fatalf("expected nil for zero segments")
	}
}

func TestBuildLatencySegments_SegmentCount(t *testing.T) {
	results := makeSegmentResults(20, 5*time.Millisecond)
	segs := BuildLatencySegments(results, 4)
	if len(segs) != 4 {
		t.Fatalf("expected 4 segments, got %d", len(segs))
	}
}

func TestBuildLatencySegments_TotalCountSums(t *testing.T) {
	results := makeSegmentResults(20, 5*time.Millisecond)
	segs := BuildLatencySegments(results, 4)
	total := 0
	for _, s := range segs {
		total += s.Count
	}
	if total != 20 {
		t.Fatalf("expected total count 20, got %d", total)
	}
}

func TestBuildLatencySegments_WithErrors(t *testing.T) {
	results := makeSegmentResults(10, 5*time.Millisecond)
	results[0].Error = errors.New("fail")
	results[1].Error = errors.New("fail")
	segs := BuildLatencySegments(results, 2)
	errorTotal := 0
	for _, s := range segs {
		errorTotal += s.Errors
	}
	if errorTotal != 2 {
		t.Fatalf("expected 2 errors, got %d", errorTotal)
	}
}

func TestBuildLatencySegments_P99Positive(t *testing.T) {
	results := makeSegmentResults(20, 10*time.Millisecond)
	segs := BuildLatencySegments(results, 2)
	for _, s := range segs {
		if s.Count > s.Errors && s.P99 <= 0 {
			t.Fatalf("expected P99 > 0 for segment %s", s.Label)
		}
	}
}

func TestWriteLatencySegments_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteLatencySegments(&buf, nil)
	if !strings.Contains(buf.String(), "no latency segment data") {
		t.Fatalf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteLatencySegments_ValidOutput(t *testing.T) {
	results := makeSegmentResults(12, 8*time.Millisecond)
	segs := BuildLatencySegments(results, 3)
	var buf bytes.Buffer
	WriteLatencySegments(&buf, segs)
	out := buf.String()
	if !strings.Contains(out, "seg01") {
		t.Fatalf("expected seg01 in output, got: %s", out)
	}
	if !strings.Contains(out, "Segment") {
		t.Fatalf("expected header in output")
	}
}
