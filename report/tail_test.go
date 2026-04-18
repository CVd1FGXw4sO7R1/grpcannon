package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeTailResults(durations []time.Duration) []Result {
	out := make([]Result, len(durations))
	for i, d := range durations {
		out[i] = Result{Duration: d}
	}
	return out
}

func TestCalcTailLatency_Empty(t *testing.T) {
	tl := CalcTailLatency(nil, 10)
	if tl.Count != 0 {
		t.Fatalf("expected 0 count, got %d", tl.Count)
	}
}

func TestCalcTailLatency_ZeroN(t *testing.T) {
	res := makeTailResults([]time.Duration{ms(10), ms(20)})
	tl := CalcTailLatency(res, 0)
	if tl.Count != 0 {
		t.Fatalf("expected 0 count for n=0")
	}
}

func TestCalcTailLatency_WindowTruncates(t *testing.T) {
	var durations []time.Duration
	for i := 1; i <= 20; i++ {
		durations = append(durations, ms(i))
	}
	res := makeTailResults(durations)
	tl := CalcTailLatency(res, 5)
	if tl.Count != 5 {
		t.Fatalf("expected window count 5, got %d", tl.Count)
	}
	if tl.Max != ms(20) {
		t.Fatalf("expected max 20ms, got %s", tl.Max)
	}
}

func TestCalcTailLatency_AllErrors(t *testing.T) {
	res := []Result{
		{Duration: ms(10), Error: errors.New("err")},
		{Duration: ms(20), Error: errors.New("err")},
	}
	tl := CalcTailLatency(res, 10)
	if tl.P99 != 0 {
		t.Fatalf("expected zero p99 when all errors")
	}
}

func TestCalcTailLatency_KnownPercentiles(t *testing.T) {
	durations := []time.Duration{ms(1), ms(2), ms(3), ms(4), ms(5),
		ms(6), ms(7), ms(8), ms(9), ms(10)}
	res := makeTailResults(durations)
	tl := CalcTailLatency(res, 100)
	if tl.P50 == 0 {
		t.Fatal("expected non-zero p50")
	}
	if tl.P99 < tl.P90 {
		t.Fatal("p99 should be >= p90")
	}
}

func TestWriteTailLatency_ValidOutput(t *testing.T) {
	res := makeTailResults([]time.Duration{ms(5), ms(10), ms(15)})
	tl := CalcTailLatency(res, 10)
	var buf bytes.Buffer
	WriteTailLatency(&buf, tl)
	out := buf.String()
	if !strings.Contains(out, "p50") {
		t.Fatal("expected p50 in output")
	}
}

func TestWriteTailLatency_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteTailLatency(&buf, TailLatency{})
	if !strings.Contains(buf.String(), "no results") {
		t.Fatal("expected 'no results' message")
	}
}
