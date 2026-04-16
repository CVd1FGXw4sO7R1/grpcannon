package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestWriteInflux_ValidOutput(t *testing.T) {
	results := makeResults(10, 2, 100*time.Millisecond)
	r := New(results)
	var buf bytes.Buffer
	err := WriteInflux(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "grpcannon_summary") {
		t.Error("expected grpcannon_summary line")
	}
	if !strings.Contains(out, "grpcannon_latency") {
		t.Error("expected grpcannon_latency line")
	}
	if !strings.Contains(out, "total=10i") {
		t.Errorf("expected total=10i in output, got: %s", out)
	}
	if !strings.Contains(out, "failures=2i") {
		t.Errorf("expected failures=2i in output, got: %s", out)
	}
}

func TestWriteInflux_EmptyResults(t *testing.T) {
	r := New(nil)
	var buf bytes.Buffer
	err := WriteInflux(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "total=0i") {
		t.Errorf("expected total=0i, got: %s", out)
	}
	if strings.Contains(out, "grpcannon_latency") {
		t.Error("expected no latency line for empty results")
	}
}

func TestWriteInflux_NilReport(t *testing.T) {
	var buf bytes.Buffer
	err := WriteInflux(&buf, nil)
	if err == nil {
		t.Error("expected error for nil report")
	}
}

func TestWriteInflux_SuccessRate(t *testing.T) {
	results := makeResults(5, 0, 50*time.Millisecond)
	r := New(results)
	var buf bytes.Buffer
	err := WriteInflux(&buf, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "success=5i") {
		t.Errorf("expected success=5i, got: %s", out)
	}
}
