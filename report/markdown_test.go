package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteMarkdown_ValidOutput(t *testing.T) {
	results := makeResults(10, 2, 10*time.Millisecond, 100*time.Millisecond)
	r := New(results)

	var buf bytes.Buffer
	err := WriteMarkdown(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	expected := []string{
		"# gRPCannon Load Test Report",
		"## Summary",
		"## Latency Percentiles",
		"| Total Requests | 10 |",
		"| Successful | 8 |",
		"| Failed | 2 |",
		"| p50 |",
		"| p90 |",
		"| p95 |",
		"| p99 |",
	}

	for _, e := range expected {
		if !strings.Contains(out, e) {
			t.Errorf("expected output to contain %q\ngot:\n%s", e, out)
		}
	}
}

func TestWriteMarkdown_EmptyResults(t *testing.T) {
	r := New(nil)

	var buf bytes.Buffer
	err := WriteMarkdown(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "| Success Rate | N/A |") {
		t.Errorf("expected N/A success rate for empty results, got:\n%s", out)
	}
}

func TestWriteMarkdown_NilReport(t *testing.T) {
	var buf bytes.Buffer
	err := WriteMarkdown(&buf, nil)
	if err == nil {
		t.Fatal("expected error for nil report, got nil")
	}
}

func TestWriteMarkdown_SuccessRate(t *testing.T) {
	results := makeResults(4, 1, 5*time.Millisecond, 50*time.Millisecond)
	r := New(results)

	var buf bytes.Buffer
	_ = WriteMarkdown(&buf, r)

	out := buf.String()
	if !strings.Contains(out, "75.00%") {
		t.Errorf("expected 75.00%% success rate in output, got:\n%s", out)
	}
}
