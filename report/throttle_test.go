package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeThrottleResults(n int, window time.Duration, errEvery int) []Result {
	base := time.Now()
	out := make([]Result, n)
	for i := 0; i < n; i++ {
		var err error
		if errEvery > 0 && i%errEvery == 0 {
			err = errors.New("rpc error")
		}
		out[i] = Result{
			Timestamp: base.Add(time.Duration(i) * (window / time.Duration(n/2+1))),
			Duration:  10 * time.Millisecond,
			Error:     err,
		}
	}
	return out
}

func TestCalcThrottle_Empty(t *testing.T) {
	got := CalcThrottle(nil, time.Second)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d windows", len(got))
	}
}

func TestCalcThrottle_ZeroWindow(t *testing.T) {
	results := makeThrottleResults(10, time.Second, 0)
	got := CalcThrottle(results, 0)
	if len(got) != 0 {
		t.Fatalf("expected empty for zero window, got %d", len(got))
	}
}

func TestCalcThrottle_TotalCount(t *testing.T) {
	results := makeThrottleResults(20, 2*time.Second, 0)
	windows := CalcThrottle(results, time.Second)
	total := 0
	for _, w := range windows {
		total += w.Total
	}
	if total != 20 {
		t.Fatalf("expected total 20, got %d", total)
	}
}

func TestCalcThrottle_SuccessCount(t *testing.T) {
	results := makeThrottleResults(10, time.Second, 2) // every 2nd is an error => 5 errors
	windows := CalcThrottle(results, 500*time.Millisecond)
	successes, failures := 0, 0
	for _, w := range windows {
		successes += w.Successes
		failures += w.Failures
	}
	if successes+failures != 10 {
		t.Fatalf("successes+failures should equal 10, got %d", successes+failures)
	}
	if failures == 0 {
		t.Fatal("expected some failures")
	}
}

func TestCalcThrottle_RPSPositive(t *testing.T) {
	results := makeThrottleResults(12, 2*time.Second, 0)
	windows := CalcThrottle(results, time.Second)
	for _, w := range windows {
		if w.RPS <= 0 {
			t.Fatalf("expected positive RPS, got %f", w.RPS)
		}
	}
}

func TestWriteThrottle_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteThrottle(&buf, nil)
	if !strings.Contains(buf.String(), "no throttle") {
		t.Fatalf("expected no-data message, got: %s", buf.String())
	}
}

func TestWriteThrottle_ValidOutput(t *testing.T) {
	results := makeThrottleResults(20, 2*time.Second, 0)
	windows := CalcThrottle(results, time.Second)
	var buf bytes.Buffer
	WriteThrottle(&buf, windows)
	out := buf.String()
	if !strings.Contains(out, "window") {
		t.Fatalf("expected header, got: %s", out)
	}
	if !strings.Contains(out, "rps") {
		t.Fatalf("expected rps column, got: %s", out)
	}
}
