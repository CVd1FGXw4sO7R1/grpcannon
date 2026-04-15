package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWriteCSV_ValidOutput(t *testing.T) {
	results := makeResults()
	r := New(results)

	var buf bytes.Buffer
	if err := WriteCSV(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "status,duration_ms,error") {
		t.Errorf("expected CSV header, got:\n%s", output)
	}

	if !strings.Contains(output, "OK") {
		t.Errorf("expected OK status in output, got:\n%s", output)
	}
}

func TestWriteCSV_EmptyResults(t *testing.T) {
	r := New(nil)

	var buf bytes.Buffer
	if err := WriteCSV(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 1 {
		t.Errorf("expected only header line for empty results, got %d lines", len(lines))
	}
}

func TestWriteCSV_WithErrors(t *testing.T) {
	results := []Result{
		{Status: "OK", Duration: 10 * time.Millisecond, Error: nil},
		{Status: "ERROR", Duration: 5 * time.Millisecond, Error: errors.New("deadline exceeded")},
	}
	r := New(results)

	var buf bytes.Buffer
	if err := WriteCSV(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "deadline exceeded") {
		t.Errorf("expected error message in CSV output, got:\n%s", output)
	}

	if !strings.Contains(output, "ERROR") {
		t.Errorf("expected ERROR status in CSV output, got:\n%s", output)
	}
}
