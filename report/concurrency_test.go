package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeSeries(workers ...int) ConcurrencySeries {
	base := time.Now()
	cs := make(ConcurrencySeries, len(workers))
	for i, w := range workers {
		cs[i] = ConcurrencySnapshot{At: base.Add(time.Duration(i) * 100 * time.Millisecond), Workers: w}
	}
	return cs
}

func TestConcurrencySeries_Peak(t *testing.T) {
	cs := makeSeries(2, 5, 3, 1)
	if got := cs.Peak(); got != 5 {
		t.Errorf("expected peak 5, got %d", got)
	}
}

func TestConcurrencySeries_Peak_Empty(t *testing.T) {
	var cs ConcurrencySeries
	if got := cs.Peak(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestConcurrencySeries_Average(t *testing.T) {
	cs := makeSeries(2, 4)
	if got := cs.Average(); got != 3.0 {
		t.Errorf("expected 3.0, got %f", got)
	}
}

func TestConcurrencySeries_Average_Empty(t *testing.T) {
	var cs ConcurrencySeries
	if got := cs.Average(); got != 0 {
		t.Errorf("expected 0, got %f", got)
	}
}

func TestWriteConcurrency_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	cs := makeSeries(1, 3, 2)
	if err := WriteConcurrency(&buf, cs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Concurrency Over Time") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "peak: 3") {
		t.Error("expected peak annotation")
	}
}

func TestWriteConcurrency_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteConcurrency(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No concurrency data") {
		t.Error("expected empty message")
	}
}
