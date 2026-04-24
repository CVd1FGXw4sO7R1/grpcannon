package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeRateResults(n int, window time.Duration, errEvery int) []Result {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		start := base.Add(time.Duration(i) * (window / time.Duration(n)))
		var err error
		if errEvery > 0 && i%errEvery == 0 {
			err = fmt.Errorf("rpc error")
		}
		results[i] = Result{
			StartedAt: start,
			EndedAt:   start.Add(5 * time.Millisecond),
			Duration:  5 * time.Millisecond,
			Err:       err,
		}
	}
	return results
}

func TestCalcRateControl_Empty(t *testing.T) {
	out := CalcRateControl(nil, time.Second)
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestCalcRateControl_ZeroWindow(t *testing.T) {
	results := makeRateResults(10, time.Second, 0)
	out := CalcRateControl(results, 0)
	if out != nil {
		t.Fatalf("expected nil for zero window, got %v", out)
	}
}

func TestCalcRateControl_WindowCount(t *testing.T) {
	results := makeRateResults(20, 2*time.Second, 0)
	win := 500 * time.Millisecond
	out := CalcRateControl(results, win)
	if len(out) == 0 {
		t.Fatal("expected at least one window")
	}
}

func TestCalcRateControl_TotalCountSums(t *testing.T) {
	results := makeRateResults(20, 2*time.Second, 0)
	out := CalcRateControl(results, 500*time.Millisecond)
	total := 0
	for _, w := range out {
		total += w.Total
	}
	if total != len(results) {
		t.Fatalf("expected total %d across windows, got %d", len(results), total)
	}
}

func TestCalcRateControl_SuccessesWithErrors(t *testing.T) {
	results := makeRateResults(10, time.Second, 2)
	out := CalcRateControl(results, time.Second)
	for _, w := range out {
		if w.Successes > w.Total {
			t.Fatalf("successes %d > total %d", w.Successes, w.Total)
		}
	}
}

func TestCalcRateControl_RPSPositive(t *testing.T) {
	results := makeRateResults(10, time.Second, 0)
	out := CalcRateControl(results, 200*time.Millisecond)
	for _, w := range out {
		if w.Total > 0 && w.RPS <= 0 {
			t.Fatalf("expected positive RPS for non-empty window")
		}
	}
}

func TestWriteRateControl_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteRateControl(&buf, nil)
	if !strings.Contains(buf.String(), "no rate-control data") {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestWriteRateControl_ValidOutput(t *testing.T) {
	results := makeRateResults(10, time.Second, 0)
	out := CalcRateControl(results, 200*time.Millisecond)
	var buf bytes.Buffer
	WriteRateControl(&buf, out)
	if !strings.Contains(buf.String(), "Window Start") {
		t.Fatalf("expected header in output: %s", buf.String())
	}
}
