package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeWarmupResults(ds []time.Duration) []Result {
	out := make([]Result, len(ds))
	for i, d := range ds {
		out[i] = Result{Duration: d}
	}
	return out
}

func TestCalcWarmup_Empty(t *testing.T) {
	s := CalcWarmup(nil, 5)
	if s.WarmupCount != 0 {
		t.Fatal("expected zero warmup count")
	}
}

func TestCalcWarmup_ZeroN(t *testing.T) {
	res := makeWarmupResults([]time.Duration{ms(10), ms(20)})
	s := CalcWarmup(res, 0)
	if s.WarmupCount != 0 {
		t.Fatal("expected zero warmup count for n=0")
	}
}

func TestCalcWarmup_NExceedsLen(t *testing.T) {
	res := makeWarmupResults([]time.Duration{ms(10), ms(20)})
	s := CalcWarmup(res, 100)
	if s.WarmupCount != 2 {
		t.Fatalf("expected warmup count 2, got %d", s.WarmupCount)
	}
	if s.SteadyCount != 0 {
		t.Fatalf("expected steady count 0, got %d", s.SteadyCount)
	}
}

func TestCalcWarmup_ImprovementPositive(t *testing.T) {
	// warmup: slow, steady: fast
	res := makeWarmupResults([]time.Duration{ms(100), ms(100), ms(10), ms(10)})
	s := CalcWarmup(res, 2)
	if s.ImprovementPct <= 0 {
		t.Fatalf("expected positive improvement, got %.2f", s.ImprovementPct)
	}
}

func TestCalcWarmup_SplitCounts(t *testing.T) {
	res := makeWarmupResults([]time.Duration{ms(1), ms(2), ms(3), ms(4), ms(5)})
	s := CalcWarmup(res, 2)
	if s.WarmupCount != 2 || s.SteadyCount != 3 {
		t.Fatalf("unexpected split: warmup=%d steady=%d", s.WarmupCount, s.SteadyCount)
	}
}

func TestWriteWarmup_ValidOutput(t *testing.T) {
	res := makeWarmupResults([]time.Duration{ms(50), ms(50), ms(20), ms(20)})
	s := CalcWarmup(res, 2)
	var buf bytes.Buffer
	WriteWarmup(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "warmup") {
		t.Fatal("expected 'warmup' in output")
	}
	if !strings.Contains(out, "improvement") {
		t.Fatal("expected 'improvement' in output")
	}
}

func TestWriteWarmup_NoData(t *testing.T) {
	var buf bytes.Buffer
	WriteWarmup(&buf, WarmupSummary{})
	if !strings.Contains(buf.String(), "no data") {
		t.Fatal("expected 'no data' message")
	}
}
