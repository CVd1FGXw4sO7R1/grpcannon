package report

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestWriteJSON_ValidOutput(t *testing.T) {
	durations := []time.Duration{
		5 * time.Millisecond,
		10 * time.Millisecond,
		20 * time.Millisecond,
	}
	results := makeResults(durations, 2)
	s := New(results, 100*time.Millisecond)

	var buf bytes.Buffer
	if err := s.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	var js JSONSummary
	if err := json.Unmarshal(buf.Bytes(), &js); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if js.Total != 5 {
		t.Errorf("expected total=5, got %d", js.Total)
	}
	if js.Successes != 3 {
		t.Errorf("expected successes=3, got %d", js.Successes)
	}
	if js.Failures != 2 {
		t.Errorf("expected failures=2, got %d", js.Failures)
	}
	if js.TotalTimeMs != 100.0 {
		t.Errorf("expected total_time_ms=100, got %f", js.TotalTimeMs)
	}
	if js.P50Ms <= 0 {
		t.Errorf("expected positive p50, got %f", js.P50Ms)
	}
}

func TestWriteJSON_EmptyResults(t *testing.T) {
	s := New([]Result{}, 0)
	var buf bytes.Buffer
	if err := s.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}
	var js JSONSummary
	if err := json.Unmarshal(buf.Bytes(), &js); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}
	if js.Total != 0 || js.P50Ms != 0 {
		t.Errorf("unexpected values for empty summary: %+v", js)
	}
}
