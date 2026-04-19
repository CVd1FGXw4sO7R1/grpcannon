package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeScoringResults(n int, errEvery int, dur time.Duration) []Result {
	results := make([]Result, n)
	for i := range results {
		var err error
		if errEvery > 0 && i%errEvery == 0 {
			err = fmt.Errorf("fail")
		}
		results[i] = Result{Duration: dur, Err: err, At: time.Now().Add(time.Duration(i) * time.Millisecond)}
	}
	return results
}

func TestCalcScore_NilReport(t *testing.T) {
	sc := CalcScore(nil)
	if sc.Total != 0 {
		t.Errorf("expected 0, got %f", sc.Total)
	}
}

func TestCalcScore_EmptyResults(t *testing.T) {
	sc := CalcScore(&Report{})
	if sc.Total != 0 {
		t.Errorf("expected 0, got %f", sc.Total)
	}
}

func TestCalcScore_AllSuccessFast(t *testing.T) {
	results := makeScoringResults(100, 0, 5*time.Millisecond)
	r := &Report{Results: results}
	sc := CalcScore(r)
	if sc.Total <= 80 {
		t.Errorf("expected high score, got %.1f", sc.Total)
	}
	if sc.ErrorScore != 100.0 {
		t.Errorf("expected perfect error score, got %.1f", sc.ErrorScore)
	}
}

func TestCalcScore_HighErrorRate(t *testing.T) {
	results := makeScoringResults(100, 2, 5*time.Millisecond) // 50% errors
	r := &Report{Results: results}
	sc := CalcScore(r)
	if sc.ErrorScore >= 60 {
		t.Errorf("expected low error score, got %.1f", sc.ErrorScore)
	}
}

func TestCalcScore_SlowLatency(t *testing.T) {
	results := makeScoringResults(50, 0, 3*time.Second)
	r := &Report{Results: results}
	sc := CalcScore(r)
	if sc.LatencyScore > 5 {
		t.Errorf("expected near-zero latency score, got %.1f", sc.LatencyScore)
	}
}

func TestWriteScore_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WriteScore(&buf, nil)
	if !strings.Contains(buf.String(), "N/A") {
		t.Errorf("expected N/A message")
	}
}

func TestWriteScore_ValidOutput(t *testing.T) {
	results := makeScoringResults(100, 0, 5*time.Millisecond)
	r := &Report{Results: results}
	var buf bytes.Buffer
	WriteScore(&buf, r)
	out := buf.String()
	for _, want := range []string{"Performance Score", "Latency", "Error Rate", "Throughput"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}
