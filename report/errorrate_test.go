package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeErrorRateResults(n int, errEvery int, spacing time.Duration) []Result {
	base := time.Now()
	out := make([]Result, n)
	for i := range out {
		out[i] = Result{
			Timestamp: base.Add(time.Duration(i) * spacing),
			Duration:  ms(10),
		}
		if errEvery > 0 && i%errEvery == 0 {
			out[i].Err = errors.New("injected")
		}
	}
	return out
}

func TestCalcErrorRate_Empty(t *testing.T) {
	got := CalcErrorRate(nil, 5)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCalcErrorRate_ZeroWindows(t *testing.T) {
	results := makeErrorRateResults(20, 0, 10*time.Millisecond)
	got := CalcErrorRate(results, 0)
	// defaults to 10
	if len(got) == 0 {
		t.Fatal("expected non-empty windows")
	}
}

func TestCalcErrorRate_WindowCount(t *testing.T) {
	results := makeErrorRateResults(50, 0, 10*time.Millisecond)
	got := CalcErrorRate(results, 5)
	if len(got) != 5 {
		t.Fatalf("expected 5 windows, got %d", len(got))
	}
}

func TestCalcErrorRate_TotalCountSums(t *testing.T) {
	n := 40
	results := makeErrorRateResults(n, 0, 10*time.Millisecond)
	windows := CalcErrorRate(results, 4)
	total := 0
	for _, w := range windows {
		total += w.Total
	}
	if total != n {
		t.Fatalf("expected total %d across windows, got %d", n, total)
	}
}

func TestCalcErrorRate_ErrorCounts(t *testing.T) {
	// every 2nd result is an error → 50 % error rate
	results := makeErrorRateResults(20, 2, 10*time.Millisecond)
	windows := CalcErrorRate(results, 2)
	for _, w := range windows {
		if w.ErrorRate < 0 || w.ErrorRate > 1 {
			t.Fatalf("error rate out of range: %f", w.ErrorRate)
		}
		if w.Errors > w.Total {
			t.Fatalf("errors %d > total %d", w.Errors, w.Total)
		}
	}
}

func TestCalcErrorRate_AllSuccess(t *testing.T) {
	results := makeErrorRateResults(30, 0, 5*time.Millisecond)
	windows := CalcErrorRate(results, 3)
	for _, w := range windows {
		if w.Errors != 0 {
			t.Fatalf("expected 0 errors, got %d", w.Errors)
		}
		if w.ErrorRate != 0 {
			t.Fatalf("expected 0 error rate, got %f", w.ErrorRate)
		}
	}
}

func TestWriteErrorRate_ValidOutput(t *testing.T) {
	results := makeErrorRateResults(20, 2, 10*time.Millisecond)
	windows := CalcErrorRate(results, 4)
	var buf bytes.Buffer
	WriteErrorRate(&buf, windows)
	out := buf.String()
	if !strings.Contains(out, "window_start") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "error_rate") {
		t.Errorf("expected error_rate column, got:\n%s", out)
	}
}

func TestWriteErrorRate_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteErrorRate(&buf, nil)
	if !strings.Contains(buf.String(), "no error rate data") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
