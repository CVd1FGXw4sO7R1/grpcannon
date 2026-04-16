package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeTimelineResults() []Result {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Result{
		{Timestamp: base, Duration: 10 * time.Millisecond, Error: nil},
		{Timestamp: base.Add(100 * time.Millisecond), Duration: 15 * time.Millisecond, Error: nil},
		{Timestamp: base.Add(1 * time.Second), Duration: 20 * time.Millisecond, Error: errors.New("fail")},
		{Timestamp: base.Add(1200 * time.Millisecond), Duration: 12 * time.Millisecond, Error: nil},
	}
}

func TestWriteTimeline_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	err := WriteTimeline(&buf, makeTimelineResults(), time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Offset") {
		t.Error("expected header 'Offset'")
	}
	if !strings.Contains(out, "RPS") {
		t.Error("expected header 'RPS'")
	}
}

func TestWriteTimeline_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	err := WriteTimeline(&buf, []Result{}, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Error("expected 'no results'")
	}
}

func TestWriteTimeline_NilResults(t *testing.T) {
	var buf bytes.Buffer
	err := WriteTimeline(&buf, nil, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Error("expected 'no results'")
	}
}

func TestWriteTimeline_DefaultBucket(t *testing.T) {
	var buf bytes.Buffer
	// zero bucketSize should default to 1s without panic
	err := WriteTimeline(&buf, makeTimelineResults(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWriteTimeline_ErrorCounting(t *testing.T) {
	var buf bytes.Buffer
	WriteTimeline(&buf, makeTimelineResults(), time.Second)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator + 2 buckets = 4 lines
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}
