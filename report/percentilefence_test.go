package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeFenceResults(n int, dur time.Duration) []Result {
	results := make([]Result, n)
	for i := range results {
		results[i] = Result{Duration: dur}
	}
	return results
}

func TestBuildPercentileFence_Empty(t *testing.T) {
	fr := BuildPercentileFence(nil, map[float64]float64{99: 100})
	if fr.TotalFences != 0 {
		t.Errorf("expected 0 fences, got %d", fr.TotalFences)
	}
}

func TestBuildPercentileFence_NoThresholds(t *testing.T) {
	r := &Report{Results: makeFenceResults(10, 50*time.Millisecond)}
	fr := BuildPercentileFence(r, map[float64]float64{})
	if fr.TotalFences != 0 {
		t.Errorf("expected 0 fences, got %d", fr.TotalFences)
	}
}

func TestBuildPercentileFence_AllPass(t *testing.T) {
	r := &Report{Results: makeFenceResults(100, 10*time.Millisecond)}
	thresholds := map[float64]float64{50: 100, 90: 200, 99: 500}
	fr := BuildPercentileFence(r, thresholds)
	if fr.Breached != 0 {
		t.Errorf("expected 0 breached, got %d", fr.Breached)
	}
	if fr.Passed != 3 {
		t.Errorf("expected 3 passed, got %d", fr.Passed)
	}
}

func TestBuildPercentileFence_Breached(t *testing.T) {
	// Mix of fast and slow results to push P99 high.
	results := makeFenceResults(99, 10*time.Millisecond)
	results = append(results, Result{Duration: 500 * time.Millisecond})
	r := &Report{Results: results}
	thresholds := map[float64]float64{99: 50}
	fr := BuildPercentileFence(r, thresholds)
	if fr.Breached != 1 {
		t.Errorf("expected 1 breached, got %d", fr.Breached)
	}
}

func TestBuildPercentileFence_SkipsErrors(t *testing.T) {
	err := fmt.Errorf("rpc error")
	results := []Result{
		{Duration: 10 * time.Millisecond},
		{Duration: 10 * time.Millisecond, Error: err},
		{Duration: 10 * time.Millisecond},
	}
	r := &Report{Results: results}
	thresholds := map[float64]float64{50: 100}
	fr := BuildPercentileFence(r, thresholds)
	if fr.TotalFences != 1 {
		t.Errorf("expected 1 fence, got %d", fr.TotalFences)
	}
}

func TestWritePercentileFence_ValidOutput(t *testing.T) {
	r := &Report{Results: makeFenceResults(100, 20*time.Millisecond)}
	thresholds := map[float64]float64{50: 100, 99: 500}
	fr := BuildPercentileFence(r, thresholds)
	var buf bytes.Buffer
	WritePercentileFence(&buf, fr)
	out := buf.String()
	if !strings.Contains(out, "Percentile Fence Report") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK status in output")
	}
}

func TestWritePercentileFence_EmptyReport(t *testing.T) {
	var buf bytes.Buffer
	WritePercentileFence(&buf, &PercentileFenceReport{})
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' message")
	}
}
