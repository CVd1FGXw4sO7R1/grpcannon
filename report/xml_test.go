package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteXML_ValidOutput(t *testing.T) {
	results := makeResults(10, 2)
	r := New(results)

	var buf bytes.Buffer
	if err := WriteXML(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{
		"<report>",
		"<total>12</total>",
		"<successes>10</successes>",
		"<failures>2</failures>",
		"<p50>",
		"<p99>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWriteXML_EmptyResults(t *testing.T) {
	r := New(nil)

	var buf bytes.Buffer
	if err := WriteXML(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<total>0</total>") {
		t.Errorf("expected zero total, got:\n%s", out)
	}
}

func TestWriteXML_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteXML(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<report/>") {
		t.Errorf("expected empty report tag, got: %s", out)
	}
}

func TestWriteXML_SuccessRate(t *testing.T) {
	results := makeResults(8, 2)
	r := New(results)

	var buf bytes.Buffer
	if err := WriteXML(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<success_rate>80") {
		t.Errorf("expected success_rate to start with 80, got:\n%s", out)
	}
}
