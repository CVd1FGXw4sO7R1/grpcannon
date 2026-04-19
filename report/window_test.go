package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeWindowResults(n int, dur time.Duration, errEvery int, spacing time.Duration) []Result {
	base := time.Now()
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		r := Result{Duration: dur, StartedAt: base.Add(time.Duration(i) * spacing)}
		if errEvery > 0 && (i+1)%errEvery == 0 {
			r.Err = errors.New("fail")
		}
		results[i] = r
	}
	return results
}

func TestCalcWindows_Empty(t *testing.T) {
	out := CalcWindows(nil, time.Second)
	if out != nil {
		t.Errorf("expected nil, got %v", out)
	}
}

func TestCalcWindows_ZeroWindow(t *testing.T) {
	results := makeWindowResults(5, 10*time.Millisecond, 0, 100*time.Millisecond)
	out := CalcWindows(results, 0)
	if out != nil {
		t.Errorf("expected nil for zero window")
	}
}

func TestCalcWindows_BucketCount(t *testing.T) {
	// 10 results spaced 100ms apart = 1s total; 500ms window => ~2 buckets
	results := makeWindowResults(10, 5*time.Millisecond, 0, 100*time.Millisecond)
	win := CalcWindows(results, 500*time.Millisecond)
	if len(win) < 2 {
		t.Errorf("expected at least 2 windows, got %d", len(win))
	}
}

func TestCalcWindows_TotalCount(t *testing.T) {
	results := makeWindowResults(8, 5*time.Millisecond, 0, 50*time.Millisecond)
	win := CalcWindows(results, time.Second)
	if len(win) != 1 {
		t.Fatalf("expected 1 window, got %d", len(win))
	}
	if win[0].Total != 8 {
		t.Errorf("expected total 8, got %d", win[0].Total)
	}
}

func TestCalcWindows_ErrorCounts(t *testing.T) {
	results := makeWindowResults(10, 5*time.Millisecond, 2, 50*time.Millisecond)
	win := CalcWindows(results, time.Second)
	var failures int
	for _, w := range win {
		failures += w.Failures
	}
	if failures != 5 {
		t.Errorf("expected 5 failures, got %d", failures)
	}
}

func TestWriteWindows_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteWindows(&buf, nil)
	if !strings.Contains(buf.String(), "no window") {
		t.Errorf("expected 'no window' message")
	}
}

func TestWriteWindows_ValidOutput(t *testing.T) {
	results := makeWindowResults(4, 10*time.Millisecond, 0, 50*time.Millisecond)
	win := CalcWindows(results, time.Second)
	var buf bytes.Buffer
	WriteWindows(&buf, win)
	out := buf.String()
	if !strings.Contains(out, "window") {
		t.Errorf("expected header in output")
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected window row in output")
	}
}
