package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeDeadLetterResults() []Result {
	return []Result{
		{Duration: 10 * time.Millisecond, Err: nil},
		{Duration: 50 * time.Millisecond, Err: errors.New("unavailable")},
		{Duration: 20 * time.Millisecond, Err: nil},
		{Duration: 80 * time.Millisecond, Err: errors.New("timeout")},
		{Duration: 5 * time.Millisecond, Err: errors.New("cancelled")},
	}
}

func TestBuildDeadLetterQueue_Empty(t *testing.T) {
	dlq := BuildDeadLetterQueue(nil, 10)
	if len(dlq.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(dlq.Entries))
	}
}

func TestBuildDeadLetterQueue_ZeroN(t *testing.T) {
	dlq := BuildDeadLetterQueue(makeDeadLetterResults(), 0)
	if len(dlq.Entries) != 0 {
		t.Errorf("expected 0 entries for n=0, got %d", len(dlq.Entries))
	}
}

func TestBuildDeadLetterQueue_Total(t *testing.T) {
	dlq := BuildDeadLetterQueue(makeDeadLetterResults(), 10)
	if dlq.Total != 3 {
		t.Errorf("expected total=3, got %d", dlq.Total)
	}
}

func TestBuildDeadLetterQueue_LimitN(t *testing.T) {
	dlq := BuildDeadLetterQueue(makeDeadLetterResults(), 2)
	if len(dlq.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(dlq.Entries))
	}
	if dlq.Total != 3 {
		t.Errorf("expected total=3, got %d", dlq.Total)
	}
}

func TestBuildDeadLetterQueue_SortedByLatencyDesc(t *testing.T) {
	dlq := BuildDeadLetterQueue(makeDeadLetterResults(), 10)
	for i := 1; i < len(dlq.Entries); i++ {
		if dlq.Entries[i].LatencyMs > dlq.Entries[i-1].LatencyMs {
			t.Errorf("entries not sorted descending at index %d", i)
		}
	}
}

func TestWriteDeadLetterQueue_ValidOutput(t *testing.T) {
	dlq := BuildDeadLetterQueue(makeDeadLetterResults(), 10)
	var buf bytes.Buffer
	WriteDeadLetterQueue(&buf, dlq)
	out := buf.String()
	if !strings.Contains(out, "Dead Letter Queue") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "timeout") && !strings.Contains(out, "unavailable") {
		t.Errorf("expected error messages in output, got: %s", out)
	}
}

func TestWriteDeadLetterQueue_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WriteDeadLetterQueue(&buf, nil)
	if !strings.Contains(buf.String(), "nil") {
		t.Errorf("expected nil message, got: %s", buf.String())
	}
}

func TestWriteDeadLetterQueue_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteDeadLetterQueue(&buf, &DeadLetterQueue{})
	if !strings.Contains(buf.String(), "no failures") {
		t.Errorf("expected no failures message, got: %s", buf.String())
	}
}
