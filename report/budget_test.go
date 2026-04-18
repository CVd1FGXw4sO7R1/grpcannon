package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeBudgetResults(durations []time.Duration) *Report {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
	}
	return &Report{Results: results}
}

func TestEvaluateBudget_NilReport(t *testing.T) {
	br := EvaluateBudget(nil, 100*time.Millisecond)
	if br.Overall {
		t.Error("expected Overall=false for nil report")
	}
}

func TestEvaluateBudget_AllPass(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}
	r := makeBudgetResults(durations)
	br := EvaluateBudget(r, 100*time.Millisecond)
	if !br.Overall {
		t.Error("expected all to pass within 100ms budget")
	}
	if !br.P50Pass || !br.P95Pass || !br.P99Pass {
		t.Error("expected all percentiles to pass")
	}
}

func TestEvaluateBudget_P99Fail(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		10 * time.Millisecond,
		10 * time.Millisecond,
		10 * time.Millisecond,
		500 * time.Millisecond,
	}
	r := makeBudgetResults(durations)
	br := EvaluateBudget(r, 50*time.Millisecond)
	if br.Overall {
		t.Error("expected Overall=false when p99 exceeds budget")
	}
	if !br.P50Pass {
		t.Error("expected P50 to pass")
	}
	if br.P99Pass {
		t.Error("expected P99 to fail")
	}
}

func TestWriteBudget_ValidOutput(t *testing.T) {
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}
	r := makeBudgetResults(durations)
	br := EvaluateBudget(r, 100*time.Millisecond)
	var buf bytes.Buffer
	WriteBudget(&buf, br)
	out := buf.String()
	if !strings.Contains(out, "Latency Budget") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "PASS") {
		t.Error("expected PASS in output")
	}
}

func TestWriteBudget_NilResult(t *testing.T) {
	var buf bytes.Buffer
	WriteBudget(&buf, nil)
	if !strings.Contains(buf.String(), "no budget result") {
		t.Error("expected nil message")
	}
}
