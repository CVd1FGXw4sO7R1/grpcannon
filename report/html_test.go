package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteHTML_ValidOutput(t *testing.T) {
	results := makeResults(10, 2)
	r := New(results)

	var buf bytes.Buffer
	if err := WriteHTML(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{
		"<!DOCTYPE html>",
		"grpcannon Load Test Report",
		"Total Requests",
		"Successful",
		"Failed",
		"Success Rate",
		"P95 Latency",
		"P99 Latency",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestWriteHTML_EmptyResults(t *testing.T) {
	r := New(nil)

	var buf bytes.Buffer
	if err := WriteHTML(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<!DOCTYPE html>") {
		t.Error("expected valid HTML even for empty results")
	}
}

func TestWriteHTML_NilReport(t *testing.T) {
	var buf bytes.Buffer
	err := WriteHTML(&buf, nil)
	if err == nil {
		t.Error("expected error for nil report")
	}
}

func TestWriteHTML_SuccessRate(t *testing.T) {
	results := makeResults(5, 5)
	r := New(results)

	var buf bytes.Buffer
	if err := WriteHTML(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "50.00%") {
		t.Errorf("expected 50.00%% success rate in output, got:\n%s", out)
	}
}
