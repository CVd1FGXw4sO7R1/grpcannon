package report

import (
	"strings"
	"testing"
	"time"
)

func TestWriteFlamegraph_ValidOutput(t *testing.T) {
	results := []Result{
		{Duration: 2 * time.Millisecond},
		{Duration: 7 * time.Millisecond},
		{Duration: 15 * time.Millisecond},
		{Duration: 30 * time.Millisecond},
		{Duration: 75 * time.Millisecond},
		{Duration: 200 * time.Millisecond},
	}
	var sb strings.Builder
	err := WriteFlamegraph(&sb, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "Flamegraph") {
		t.Error("expected header 'Flamegraph'")
	}
	for _, label := range []string{"0-5ms", "5-10ms", "10-25ms", "25-50ms", "50-100ms", "100ms+"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in output", label)
		}
	}
}

func TestWriteFlamegraph_EmptyResults(t *testing.T) {
	var sb strings.Builder
	err := WriteFlamegraph(&sb, []Result{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no results") {
		t.Error("expected 'no results' message")
	}
}

func TestWriteFlamegraph_NilResults(t *testing.T) {
	var sb strings.Builder
	err := WriteFlamegraph(&sb, nil)
	if err == nil {
		t.Error("expected error for nil results")
	}
}

func TestWriteFlamegraph_SuccessRate(t *testing.T) {
	results := make([]Result, 20)
	for i := range results {
		results[i] = Result{Duration: time.Duration(i+1) * time.Millisecond}
	}
	var sb strings.Builder
	err := WriteFlamegraph(&sb, results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "%") {
		t.Error("expected percentage values in output")
	}
}
