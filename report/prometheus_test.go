package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWritePrometheus_ValidOutput(t *testing.T) {
	results := makeResults(10, 2)
	r := New(results, time.Second)

	var buf bytes.Buffer
	err := WritePrometheus(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	expectedFragments := []string{
		"grpcannon_requests_total 10",
		"grpcannon_requests_success 8",
		"grpcannon_requests_failed 2",
		"grpcannon_duration_seconds",
		"grpcannon_rps",
		"grpcannon_latency_seconds{quantile=\"0.9\"}",
		"grpcannon_latency_seconds{quantile=\"0.95\"}",
		"grpcannon_latency_seconds{quantile=\"0.99\"}",
		"grpcannon_latency_seconds_fastest",
		"grpcannon_latency_seconds_slowest",
		"grpcannon_latency_seconds_mean",
		"# HELP grpcannon_requests_total",
		"# TYPE grpcannon_requests_total counter",
	}

	for _, frag := range expectedFragments {
		if !strings.Contains(output, frag) {
			t.Errorf("expected output to contain %q, got:\n%s", frag, output)
		}
	}
}

func TestWritePrometheus_EmptyResults(t *testing.T) {
	r := New(nil, time.Second)

	var buf bytes.Buffer
	err := WritePrometheus(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "grpcannon_requests_total 0") {
		t.Errorf("expected zero total, got:\n%s", output)
	}
}

func TestWritePrometheus_NilReport(t *testing.T) {
	var buf bytes.Buffer
	err := WritePrometheus(&buf, nil)
	if err == nil {
		t.Fatal("expected error for nil report, got nil")
	}
}

func TestWritePrometheus_SuccessRate(t *testing.T) {
	results := makeResults(5, 5)
	r := New(results, 2*time.Second)

	var buf bytes.Buffer
	_ = WritePrometheus(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "grpcannon_requests_success 0") {
		t.Errorf("expected 0 successes when all fail, got:\n%s", output)
	}
	if !strings.Contains(output, "grpcannon_requests_failed 5") {
		t.Errorf("expected 5 failures, got:\n%s", output)
	}
}
