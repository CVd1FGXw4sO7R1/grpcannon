package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeStabilityResults(n int, errEvery int) []Result {
	base := time.Now()
	out := make([]Result, n)
	for i := range out {
		out[i] = Result{
			Start:    base.Add(time.Duration(i) * 10 * time.Millisecond),
			Duration: 20 * time.Millisecond,
		}
		if errEvery > 0 && i%errEvery == 0 {
			out[i].Error = errSentinel
		}
	}
	return out
}

func TestBuildStability_Empty(t *testing.T) {
	sr := BuildStability(nil, 5, 50, 0.05)
	if !sr.Stable {
		t.Error("expected stable for empty input")
	}
}

func TestBuildStability_ZeroWindows(t *testing.T) {
	res := makeStabilityResults(20, 0)
	sr := BuildStability(res, 0, 50, 0.05)
	if !sr.Stable {
		t.Error("expected stable for zero windows")
	}
}

func TestBuildStability_WindowCount(t *testing.T) {
	res := makeStabilityResults(40, 0)
	sr := BuildStability(res, 4, 200, 0.1)
	if len(sr.Windows) == 0 {
		t.Error("expected at least one window")
	}
	if len(sr.Windows) > 4 {
		t.Errorf("expected at most 4 windows, got %d", len(sr.Windows))
	}
}

func TestBuildStability_StableUniformLatency(t *testing.T) {
	res := makeStabilityResults(60, 0)
	sr := BuildStability(res, 6, 100, 0.05)
	if !sr.Stable {
		t.Errorf("expected stable, got reason: %s", sr.Reason)
	}
}

func TestBuildStability_UnstableHighJitter(t *testing.T) {
	base := time.Now()
	res := make([]Result, 40)
	for i := range res {
		d := 10 * time.Millisecond
		if i >= 20 {
			d = 300 * time.Millisecond
		}
		res[i] = Result{Start: base.Add(time.Duration(i) * 20 * time.Millisecond), Duration: d}
	}
	sr := BuildStability(res, 2, 50, 0.1)
	if sr.Stable {
		t.Error("expected unstable due to high jitter")
	}
	if !strings.Contains(sr.Reason, "jitter") {
		t.Errorf("expected jitter in reason, got: %s", sr.Reason)
	}
}

func TestBuildStability_MaxP99JitterPositive(t *testing.T) {
	base := time.Now()
	res := make([]Result, 30)
	for i := range res {
		d := time.Duration(i+1) * 5 * time.Millisecond
		res[i] = Result{Start: base.Add(time.Duration(i) * 15 * time.Millisecond), Duration: d}
	}
	sr := BuildStability(res, 3, 10000, 1.0)
	if sr.MaxP99JitterMs < 0 {
		t.Error("max p99 jitter should be non-negative")
	}
}

func TestWriteStability_ValidOutput(t *testing.T) {
	res := makeStabilityResults(20, 0)
	sr := BuildStability(res, 2, 200, 0.1)
	var buf bytes.Buffer
	WriteStability(&buf, sr)
	out := buf.String()
	if !strings.Contains(out, "Stability:") {
		t.Error("expected 'Stability:' in output")
	}
	if !strings.Contains(out, "Max P99 Jitter") {
		t.Error("expected jitter line in output")
	}
}

func TestWriteStability_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WriteStability(&buf, nil)
	if !strings.Contains(buf.String(), "no report") {
		t.Error("expected 'no report' for nil input")
	}
}
