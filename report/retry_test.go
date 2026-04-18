package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func makeRetryResults() []Result {
	return []Result{
		{Duration: ms(10), Err: nil},
		{Duration: ms(20), Err: nil},
		{Duration: ms(5), Err: errors.New("unavailable")},
	}
}

func TestBuildRetryBreakdown_Empty(t *testing.T) {
	buckets := BuildRetryBreakdown(nil)
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(buckets))
	}
}

func TestBuildRetryBreakdown_Counts(t *testing.T) {
	results := makeRetryResults()
	buckets := BuildRetryBreakdown(results)
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != len(results) {
		t.Fatalf("expected total %d, got %d", len(results), total)
	}
}

func TestBuildRetryBreakdown_SuccessAttempt1(t *testing.T) {
	results := []Result{
		{Duration: ms(10), Err: nil},
		{Duration: ms(15), Err: nil},
	}
	buckets := BuildRetryBreakdown(results)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Attempt != 1 || buckets[0].Count != 2 {
		t.Fatalf("unexpected bucket: %+v", buckets[0])
	}
}

func TestWriteRetryBreakdown_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	WriteRetryBreakdown(&buf, makeRetryResults())
	out := buf.String()
	if !strings.Contains(out, "Retry Breakdown") {
		t.Error("expected header")
	}
	if !strings.Contains(out, "error") {
		t.Error("expected error row")
	}
}

func TestWriteRetryBreakdown_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteRetryBreakdown(&buf, nil)
	if !strings.Contains(buf.String(), "No results") {
		t.Error("expected no results message")
	}
}
