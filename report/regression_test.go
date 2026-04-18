package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeRegressionReport(p50, p99 time.Duration, errRate float64) *Report {
	r := &Report{ErrorRate: errRate}
	count := 20
	for i := 0; i < count; i++ {
		d := p50
		if i == count-1 {
			d = p99
		}
		r.Results = append(r.Results, Result{Duration: d})
	}
	return r
}

func TestEvaluateRegression_NilReports(t *testing.T) {
	results := EvaluateRegression(nil, nil, RegressionThresholds{})
	if results != nil {
		t.Error("expected nil for nil reports")
	}
}

func TestEvaluateRegression_AllPass(t *testing.T) {
	base := makeRegressionReport(10*time.Millisecond, 20*time.Millisecond, 0.01)
	cur := makeRegressionReport(11*time.Millisecond, 21*time.Millisecond, 0.01)
	thresh := RegressionThresholds{MaxP50DeltaMs: 5, MaxP99DeltaMs: 5, MaxErrorDelta: 1}
	results := EvaluateRegression(base, cur, thresh)
	for _, r := range results {
		if !r.Passed {
			t.Errorf("expected %s to pass", r.Field)
		}
	}
}

func TestEvaluateRegression_P99Fail(t *testing.T) {
	base := makeRegressionReport(10*time.Millisecond, 20*time.Millisecond, 0.0)
	cur := makeRegressionReport(10*time.Millisecond, 100*time.Millisecond, 0.0)
	thresh := RegressionThresholds{MaxP50DeltaMs: 5, MaxP99DeltaMs: 5, MaxErrorDelta: 1}
	results := EvaluateRegression(base, cur, thresh)
	var p99Result *RegressionResult
	for i := range results {
		if results[i].Field == "P99 latency (ms)" {
			p99Result = &results[i]
		}
	}
	if p99Result == nil {
		t.Fatal("P99 result missing")
	}
	if p99Result.Passed {
		t.Error("expected P99 to fail")
	}
}

func TestEvaluateRegression_ResultCount(t *testing.T) {
	base := makeRegressionReport(10*time.Millisecond, 20*time.Millisecond, 0.0)
	cur := makeRegressionReport(10*time.Millisecond, 20*time.Millisecond, 0.0)
	results := EvaluateRegression(base, cur, RegressionThresholds{MaxP50DeltaMs: 10, MaxP99DeltaMs: 10, MaxErrorDelta: 5})
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestWriteRegression_ValidOutput(t *testing.T) {
	base := makeRegressionReport(10*time.Millisecond, 20*time.Millisecond, 0.01)
	cur := makeRegressionReport(12*time.Millisecond, 25*time.Millisecond, 0.02)
	results := EvaluateRegression(base, cur, RegressionThresholds{MaxP50DeltaMs: 5, MaxP99DeltaMs: 5, MaxErrorDelta: 5})
	var buf bytes.Buffer
	WriteRegression(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "P99 latency") {
		t.Error("expected P99 latency in output")
	}
	if !strings.Contains(out, "PASS") && !strings.Contains(out, "FAIL") {
		t.Error("expected PASS or FAIL in output")
	}
}

func TestWriteRegression_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteRegression(&buf, nil)
	if !strings.Contains(buf.String(), "No regression") {
		t.Error("expected no regression message")
	}
}
