package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteText_ValidOutput(t *testing.T) {
	results := makeResults(10, 2, []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
		6 * time.Millisecond,
		7 * time.Millisecond,
		8 * time.Millisecond,
	})

	r := New(results)
	var buf bytes.Buffer

	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}

	out := buf.String()

	expected := []string{
		"Total requests:",
		"Successful:",
		"Failed:",
		"Success rate:",
		"Latency",
		"Min:",
		"Mean:",
		"Max:",
		"Percentiles",
		"p50:",
		"p90:",
		"p95:",
		"p99:",
	}

	for _, want := range expected {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWriteText_EmptyResults(t *testing.T) {
	r := New(nil)
	var buf bytes.Buffer

	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Total requests:") {
		t.Errorf("expected output to contain header, got:\n%s", out)
	}
	// success rate line should be omitted when total is 0
	if strings.Contains(out, "Success rate:") {
		t.Errorf("expected no success rate line for empty results, got:\n%s", out)
	}
}

func TestWriteText_SuccessRate(t *testing.T) {
	results := makeResults(5, 5, []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
	})

	r := New(results)
	var buf bytes.Buffer

	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "100.00%") {
		t.Errorf("expected 100.00%%%% success rate, got:\n%s", out)
	}
}

func TestWriteText_PartialFailures(t *testing.T) {
	results := makeResults(10, 3, []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
	})

	r := New(results)
	var buf bytes.Buffer

	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "30.00%") {
		t.Errorf("expected 30.00%%%% success rate, got:\n%s", out)
	}
}

// TestWriteText_FailedCount verifies that the failed count is correctly
// reported as the difference between total requests and successful ones.
func TestWriteText_FailedCount(t *testing.T) {
	results := makeResults(10, 4, []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
	})

	r := New(results)
	var buf bytes.Buffer

	if err := WriteText(r, &buf); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}

	out := buf.String()
	// 10 total - 4 successful = 6 failed
	if !strings.Contains(out, "Failed:") {
		t.Errorf("expected output to contain 'Failed:', got:\n%s", out)
	}
}
