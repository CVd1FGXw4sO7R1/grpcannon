package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeThroughputResults(total, failures int) []Result {
	results := make([]Result, total)
	for i := 0; i < failures; i++ {
		results[i] = Result{Err: errors.New("fail")}
	}
	return results
}

func TestCalcThroughput_Empty(t *testing.T) {
	stats := CalcThroughput(nil, time.Second)
	if stats.TotalRequests != 0 || stats.RPS != 0 {
		t.Errorf("expected zero stats, got %+v", stats)
	}
}

func TestCalcThroughput_ZeroDuration(t *testing.T) {
	stats := CalcThroughput(makeThroughputResults(10, 0), 0)
	if stats.RPS != 0 {
		t.Errorf("expected zero RPS for zero duration")
	}
}

func TestCalcThroughput_AllSuccess(t *testing.T) {
	results := makeThroughputResults(100, 0)
	stats := CalcThroughput(results, 10*time.Second)
	if stats.TotalRequests != 100 {
		t.Errorf("expected 100 total, got %d", stats.TotalRequests)
	}
	if stats.Failed != 0 {
		t.Errorf("expected 0 failed, got %d", stats.Failed)
	}
	if stats.RPS != 10.0 {
		t.Errorf("expected RPS=10, got %.2f", stats.RPS)
	}
	if stats.SuccessRPS != 10.0 {
		t.Errorf("expected SuccessRPS=10, got %.2f", stats.SuccessRPS)
	}
}

func TestCalcThroughput_WithFailures(t *testing.T) {
	results := makeThroughputResults(100, 25)
	stats := CalcThroughput(results, 5*time.Second)
	if stats.Successful != 75 {
		t.Errorf("expected 75 successful, got %d", stats.Successful)
	}
	if stats.Failed != 25 {
		t.Errorf("expected 25 failed, got %d", stats.Failed)
	}
	if stats.RPS != 20.0 {
		t.Errorf("expected RPS=20, got %.2f", stats.RPS)
	}
}

func TestWriteThroughput_ValidOutput(t *testing.T) {
	results := makeThroughputResults(50, 5)
	stats := CalcThroughput(results, 5*time.Second)
	var buf bytes.Buffer
	if err := WriteThroughput(&buf, stats); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Throughput", "RPS", "50", "45", "5"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\n%s", want, out)
		}
	}
}

func TestWriteThroughput_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteThroughput(&buf, ThroughputStats{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no results") {
		t.Errorf("expected 'no results' message, got: %s", buf.String())
	}
}
