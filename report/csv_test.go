package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteCSV_ValidOutput(t *testing.T) {
	results := makeResults(10, 2)
	r := New(results, 2*time.Second)

	var buf bytes.Buffer
	if err := WriteCSV(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + data), got %d", len(lines))
	}

	expectedHeaders := []string{"total", "successes", "failures", "success_rate", "rps", "min_ms", "mean_ms", "p50_ms", "p95_ms", "p99_ms", "max_ms"}
	headerLine := lines[0]
	for _, h := range expectedHeaders {
		if !strings.Contains(headerLine, h) {
			t.Errorf("header missing field %q", h)
		}
	}

	dataFields := strings.Split(lines[1], ",")
	if len(dataFields) != len(expectedHeaders) {
		t.Errorf("expected %d data fields, got %d", len(expectedHeaders), len(dataFields))
	}

	// total should be 10
	if dataFields[0] != "10" {
		t.Errorf("expected total=10, got %q", dataFields[0])
	}
	// failures should be 2
	if dataFields[2] != "2" {
		t.Errorf("expected failures=2, got %q", dataFields[2])
	}
}

func TestWriteCSV_EmptyResults(t *testing.T) {
	r := New(nil, time.Second)

	var buf bytes.Buffer
	if err := WriteCSV(r, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	dataFields := strings.Split(lines[1], ",")
	if dataFields[0] != "0" {
		t.Errorf("expected total=0, got %q", dataFields[0])
	}
	if dataFields[3] != "0.00" {
		t.Errorf("expected success_rate=0.00, got %q", dataFields[3])
	}
}
