package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteDotPlot_ValidOutput(t *testing.T) {
	results := makeResults(20, 2)
	r := New(results)
	var buf bytes.Buffer
	if err := WriteDotPlot(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Latency Distribution") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "Min:") {
		t.Errorf("expected Min in output, got: %s", out)
	}
}

func TestWriteDotPlot_EmptyResults(t *testing.T) {
	r := New(nil)
	var buf bytes.Buffer
	if err := WriteDotPlot(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteDotPlot_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteDotPlot(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteDotPlot_SuccessRate(t *testing.T) {
	results := makeResults(10, 0)
	r := New(results)
	var buf bytes.Buffer
	if err := WriteDotPlot(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.SuccessCount != 10 {
		t.Errorf("expected 10 successes, got %d", r.SuccessCount)
	}
}
