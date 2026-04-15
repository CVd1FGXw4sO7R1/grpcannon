package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteTable_ValidOutput(t *testing.T) {
	results := makeResults(10, 2)
	r := New(results)

	var buf bytes.Buffer
	if err := WriteTable(r, &buf); err != nil {
		t.Fatalf("WriteTable returned error: %v", err)
	}

	out := buf.String()
	for _, expected := range []string{
		"Total Requests",
		"Successful",
		"Failed",
		"Success Rate",
		"Min Latency",
		"Mean Latency",
		"Max Latency",
		"p50 Latency",
		"p90 Latency",
		"p95 Latency",
		"p99 Latency",
	} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, out)
		}
	}
}

func TestWriteTable_EmptyResults(t *testing.T) {
	r := New(nil)

	var buf bytes.Buffer
	if err := WriteTable(r, &buf); err != nil {
		t.Fatalf("WriteTable returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "N/A") {
		t.Errorf("expected 'N/A' for success rate on empty results, got:\n%s", out)
	}
}

func TestWriteTable_SuccessRate(t *testing.T) {
	results := makeResults(8, 2)
	r := New(results)

	var buf bytes.Buffer
	if err := WriteTable(r, &buf); err != nil {
		t.Fatalf("WriteTable returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "80.00%") {
		t.Errorf("expected 80.00%% success rate in output, got:\n%s", out)
	}
}

func TestWriteTable_CountsMatch(t *testing.T) {
	results := makeResults(5, 3)
	r := New(results)

	if r.Total != 8 {
		t.Errorf("expected Total=8, got %d", r.Total)
	}
	if r.Successful != 5 {
		t.Errorf("expected Successful=5, got %d", r.Successful)
	}
	if r.Failed != 3 {
		t.Errorf("expected Failed=3, got %d", r.Failed)
	}

	_ = time.Second // ensure time import used via makeResults
}
