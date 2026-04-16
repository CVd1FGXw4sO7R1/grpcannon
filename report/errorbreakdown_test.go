package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func makeErrorResults() []Result {
	return []Result{
		{Duration: ms(10), Error: nil},
		{Duration: ms(12), Error: errors.New("connection refused")},
		{Duration: ms(8), Error: errors.New("timeout")},
		{Duration: ms(9), Error: errors.New("connection refused")},
		{Duration: ms(11), Error: errors.New("timeout")},
		{Duration: ms(7), Error: errors.New("timeout")},
	}
}

func TestBuildErrorBreakdown_Counts(t *testing.T) {
	results := makeErrorResults()
	bd := BuildErrorBreakdown(results)
	if bd.Total != 5 {
		t.Errorf("expected total 5, got %d", bd.Total)
	}
	if bd.Errors["connection refused"] != 2 {
		t.Errorf("expected 2 connection refused, got %d", bd.Errors["connection refused"])
	}
	if bd.Errors["timeout"] != 3 {
		t.Errorf("expected 3 timeout, got %d", bd.Errors["timeout"])
	}
}

func TestBuildErrorBreakdown_NoErrors(t *testing.T) {
	results := []Result{{Duration: ms(10)}, {Duration: ms(20)}}
	bd := BuildErrorBreakdown(results)
	if bd.Total != 0 {
		t.Errorf("expected 0 total errors, got %d", bd.Total)
	}
}

func TestWriteErrorBreakdown_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteErrorBreakdown(&buf, makeErrorResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "connection refused") {
		t.Error("expected 'connection refused' in output")
	}
	if !strings.Contains(out, "timeout") {
		t.Error("expected 'timeout' in output")
	}
	if !strings.Contains(out, "5 total errors") {
		t.Error("expected total errors count in output")
	}
}

func TestWriteErrorBreakdown_NoErrors(t *testing.T) {
	var buf bytes.Buffer
	results := []Result{{Duration: ms(5)}}
	if err := WriteErrorBreakdown(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No errors") {
		t.Error("expected 'No errors' message")
	}
}

func TestWriteErrorBreakdown_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteErrorBreakdown(&buf, []Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No errors") {
		t.Error("expected 'No errors' for empty results")
	}
}
