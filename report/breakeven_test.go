package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeStepResults() []StepResult {
	return []StepResult{
		{Concurrency: 1, RPS: 50, P99Ms: 10, Duration: time.Second},
		{Concurrency: 5, RPS: 200, P99Ms: 20, Duration: time.Second},
		{Concurrency: 10, RPS: 350, P99Ms: 80, Duration: time.Second},
		{Concurrency: 20, RPS: 360, P99Ms: 300, Duration: time.Second},
	}
}

func TestFindBreakeven_Empty(t *testing.T) {
	points, best := FindBreakeven(nil)
	if len(points) != 0 {
		t.Errorf("expected 0 points, got %d", len(points))
	}
	if best != 0 {
		t.Errorf("expected best=0, got %d", best)
	}
}

func TestFindBreakeven_PointCount(t *testing.T) {
	steps := makeStepResults()
	points, _ := FindBreakeven(steps)
	if len(points) != len(steps) {
		t.Errorf("expected %d points, got %d", len(steps), len(points))
	}
}

func TestFindBreakeven_ScorePositive(t *testing.T) {
	points, _ := FindBreakeven(makeStepResults())
	for _, p := range points {
		if p.Score <= 0 {
			t.Errorf("expected positive score, got %f for concurrency %d", p.Score, p.Concurrency)
		}
	}
}

func TestFindBreakeven_OptimalConcurrency(t *testing.T) {
	// concurrency=5 has score 200/20=10, concurrency=10 has 350/80=4.375
	// concurrency=1 has 50/10=5, concurrency=20 has 360/300=1.2
	// so concurrency=5 should win
	_, best := FindBreakeven(makeStepResults())
	if best != 5 {
		t.Errorf("expected optimal concurrency=5, got %d", best)
	}
}

func TestFindBreakeven_ZeroP99Skipped(t *testing.T) {
	steps := []StepResult{
		{Concurrency: 1, RPS: 100, P99Ms: 0},
		{Concurrency: 2, RPS: 200, P99Ms: 10},
	}
	points, best := FindBreakeven(steps)
	if len(points) != 1 {
		t.Errorf("expected 1 point (zero P99 skipped), got %d", len(points))
	}
	if best != 2 {
		t.Errorf("expected best=2, got %d", best)
	}
}

func TestWriteBreakeven_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	WriteBreakeven(&buf, makeStepResults())
	out := buf.String()
	if !strings.Contains(out, "Breakeven") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "optimal") {
		t.Error("expected optimal marker in output")
	}
	if !strings.Contains(out, "Optimal concurrency") {
		t.Error("expected summary line in output")
	}
}

func TestWriteBreakeven_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteBreakeven(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Error("expected 'no data' for empty input")
	}
}
