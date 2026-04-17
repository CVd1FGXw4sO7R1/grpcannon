package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func makeStatusResults() []Result {
	return []Result{
		{Duration: ms(10), Error: nil},
		{Duration: ms(20), Error: nil},
		{Duration: ms(15), Error: errors.New("codes.Unavailable")},
		{Duration: ms(30), Error: errors.New("codes.DeadlineExceeded")},
		{Duration: ms(25), Error: errors.New("codes.Unavailable")},
	}
}

func TestBuildStatusCodeBreakdown_Counts(t *testing.T) {
	results := makeStatusResults()
	bd := BuildStatusCodeBreakdown(results)

	if bd["OK"] != 2 {
		t.Errorf("expected 2 OK, got %d", bd["OK"])
	}
	if bd["codes.Unavailable"] != 2 {
		t.Errorf("expected 2 Unavailable, got %d", bd["codes.Unavailable"])
	}
	if bd["codes.DeadlineExceeded"] != 1 {
		t.Errorf("expected 1 DeadlineExceeded, got %d", bd["codes.DeadlineExceeded"])
	}
}

func TestBuildStatusCodeBreakdown_AllSuccess(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(20)},
	}
	bd := BuildStatusCodeBreakdown(results)
	if len(bd) != 1 || bd["OK"] != 2 {
		t.Errorf("expected only OK with count 2, got %v", bd)
	}
}

func TestWriteStatusCodeBreakdown_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	WriteStatusCodeBreakdown(&buf, makeStatusResults())
	out := buf.String()

	if !strings.Contains(out, "Status Code Breakdown") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "OK") {
		t.Error("missing OK code")
	}
	if !strings.Contains(out, "codes.Unavailable") {
		t.Error("missing Unavailable code")
	}
}

func TestWriteStatusCodeBreakdown_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteStatusCodeBreakdown(&buf, []Result{})
	if !strings.Contains(buf.String(), "No results") {
		t.Error("expected no results message")
	}
}
