package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteSparkline_ValidOutput(t *testing.T) {
	results := makeResults(10, 0, 5*time.Millisecond, 50*time.Millisecond)
	r := New(results)
	var buf bytes.Buffer
	if err := WriteSparkline(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Latency distribution") {
		t.Errorf("expected header line, got: %s", out)
	}
	// sparkline characters should be present
	found := false
	for _, ch := range sparkChars {
		if strings.ContainsRune(out, ch) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected spark characters in output, got: %s", out)
	}
}

func TestWriteSparkline_EmptyResults(t *testing.T) {
	r := New(nil)
	var buf bytes.Buffer
	if err := WriteSparkline(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected no-results message")
	}
}

func TestWriteSparkline_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteSparkline(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected no-results message")
	}
}

func TestWriteSparkline_AllErrors(t *testing.T) {
	results := makeResults(0, 5, time.Millisecond, 10*time.Millisecond)
	r := New(results)
	var buf bytes.Buffer
	if err := WriteSparkline(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No successful results") {
		t.Errorf("expected no-successful-results message, got: %s", buf.String())
	}
}
