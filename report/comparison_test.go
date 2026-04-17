package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeReport(total, success int, p50, p95, p99, avg time.Duration) *Report {
	return &Report{
		Total:   total,
		Success: success,
		Failure: total - success,
		P50:     p50,
		P95:     p95,
		P99:     p99,
		Avg:     avg,
	}
}

func TestDelta_Basic(t *testing.T) {
	if got := Delta(100, 110); got != 10.0 {
		t.Fatalf("expected 10.0, got %f", got)
	}
}

func TestDelta_ZeroBaseline(t *testing.T) {
	if got := Delta(0, 50); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestWriteComparison_ValidOutput(t *testing.T) {
	base := makeReport(100, 95, ms(10), ms(20), ms(30), ms(12))
	cand := makeReport(100, 98, ms(8), ms(15), ms(25), ms(10))
	c := &ComparisonReport{
		BaselineLabel:  "v1",
		CandidateLabel: "v2",
		Baseline:       base,
		Candidate:      cand,
	}
	var buf bytes.Buffer
	if err := WriteComparison(&buf, c); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"P50", "P95", "P99", "Avg", "Success", "v1", "v2", "Delta"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestWriteComparison_NilReport(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteComparison(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "missing") {
		t.Error("expected missing message")
	}
}

func TestWriteComparison_NilBaseline(t *testing.T) {
	var buf bytes.Buffer
	c := &ComparisonReport{Candidate: makeReport(10, 10, ms(1), ms(2), ms(3), ms(1))}
	if err := WriteComparison(&buf, c); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "missing") {
		t.Error("expected missing message")
	}
}
