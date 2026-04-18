package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeJitterResults(durations []time.Duration, withErr bool) []Result {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
		if withErr && i == len(durations)-1 {
			results[i].Error = errors.New("rpc error")
		}
	}
	return results
}

func TestCalcJitter_Empty(t *testing.T) {
	if CalcJitter(nil) != nil {
		t.Fatal("expected nil for empty input")
	}
}

func TestCalcJitter_Single(t *testing.T) {
	results := makeJitterResults([]time.Duration{ms(10)}, false)
	if CalcJitter(results) != nil {
		t.Fatal("expected nil for single result")
	}
}

func TestCalcJitter_AllErrors(t *testing.T) {
	results := []Result{
		{Duration: ms(10), Error: errors.New("e")},
		{Duration: ms(20), Error: errors.New("e")},
	}
	if CalcJitter(results) != nil {
		t.Fatal("expected nil when all results are errors")
	}
}

func TestCalcJitter_KnownValues(t *testing.T) {
	// deltas: 10ms, 10ms, 10ms => mean=10ms, stddev=0
	results := makeJitterResults([]time.Duration{
		ms(10), ms(20), ms(30), ms(40),
	}, false)
	j := CalcJitter(results)
	if j == nil {
		t.Fatal("expected non-nil jitter")
	}
	if j.Min != ms(10) {
		t.Errorf("Min: got %v, want %v", j.Min, ms(10))
	}
	if j.Max != ms(10) {
		t.Errorf("Max: got %v, want %v", j.Max, ms(10))
	}
	if j.Mean != ms(10) {
		t.Errorf("Mean: got %v, want %v", j.Mean, ms(10))
	}
	if j.StdDev != 0 {
		t.Errorf("StdDev: got %v, want 0", j.StdDev)
	}
}

func TestWriteJitter_ValidOutput(t *testing.T) {
	results := makeJitterResults([]time.Duration{ms(5), ms(15), ms(25)}, false)
	var buf bytes.Buffer
	WriteJitter(&buf, results)
	out := buf.String()
	for _, s := range []string{"Jitter", "Min", "Max", "Mean", "StdDev"} {
		if !strings.Contains(out, s) {
			t.Errorf("expected %q in output", s)
		}
	}
}

func TestWriteJitter_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteJitter(&buf, nil)
	if !strings.Contains(buf.String(), "insufficient") {
		t.Error("expected insufficient data message")
	}
}
