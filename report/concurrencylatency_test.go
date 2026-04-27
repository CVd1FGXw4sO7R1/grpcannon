package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeCLGroups() map[int][]Result {
	return map[int][]Result{
		1: {
			{Duration: 10 * time.Millisecond, Err: nil},
			{Duration: 20 * time.Millisecond, Err: nil},
			{Duration: 30 * time.Millisecond, Err: nil},
		},
		4: {
			{Duration: 50 * time.Millisecond, Err: nil},
			{Duration: 100 * time.Millisecond, Err: nil},
			{Duration: 200 * time.Millisecond, Err: errors.New("timeout")},
		},
	}
}

func TestBuildConcurrencyLatency_Empty(t *testing.T) {
	points := BuildConcurrencyLatency(nil)
	if len(points) != 0 {
		t.Fatalf("expected 0 points, got %d", len(points))
	}
}

func TestBuildConcurrencyLatency_PointCount(t *testing.T) {
	points := BuildConcurrencyLatency(makeCLGroups())
	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}
}

func TestBuildConcurrencyLatency_Sorted(t *testing.T) {
	points := BuildConcurrencyLatency(makeCLGroups())
	if points[0].Concurrency >= points[1].Concurrency {
		t.Errorf("expected ascending concurrency order")
	}
}

func TestBuildConcurrencyLatency_P99Positive(t *testing.T) {
	points := BuildConcurrencyLatency(makeCLGroups())
	for _, p := range points {
		if p.P99Ms <= 0 {
			t.Errorf("concurrency=%d: expected positive P99, got %.2f", p.Concurrency, p.P99Ms)
		}
	}
}

func TestBuildConcurrencyLatency_ErrorRate(t *testing.T) {
	points := BuildConcurrencyLatency(makeCLGroups())
	var c4 ConcurrencyLatencyPoint
	for _, p := range points {
		if p.Concurrency == 4 {
			c4 = p
		}
	}
	// 1 error out of 3 results ≈ 33.33%
	if c4.ErrorRate < 30 || c4.ErrorRate > 40 {
		t.Errorf("expected ~33%% error rate, got %.2f", c4.ErrorRate)
	}
}

func TestWriteConcurrencyLatency_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteConcurrencyLatency(&buf, nil)
	if !strings.Contains(buf.String(), "no concurrency") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteConcurrencyLatency_ValidOutput(t *testing.T) {
	points := BuildConcurrencyLatency(makeCLGroups())
	var buf bytes.Buffer
	WriteConcurrencyLatency(&buf, points)
	out := buf.String()
	if !strings.Contains(out, "Concurrency") {
		t.Errorf("expected header in output")
	}
	if !strings.Contains(out, "P99") {
		t.Errorf("expected P99 column in output")
	}
}
