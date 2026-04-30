package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeBoxResults(durations []time.Duration) []Result {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
	}
	return results
}

func TestBuildBoxStats_Empty(t *testing.T) {
	bs := BuildBoxStats(nil)
	if bs == nil {
		t.Fatal("expected non-nil BoxStats")
	}
	if bs.Max != 0 || bs.Min != 0 {
		t.Errorf("expected zero stats, got min=%v max=%v", bs.Min, bs.Max)
	}
}

func TestBuildBoxStats_AllErrors(t *testing.T) {
	results := []Result{
		{Duration: ms(10), Err: errors.New("fail")},
		{Duration: ms(20), Err: errors.New("fail")},
	}
	bs := BuildBoxStats(results)
	if bs.Max != 0 {
		t.Errorf("expected zero max for all-error results, got %v", bs.Max)
	}
}

func TestBuildBoxStats_MinMax(t *testing.T) {
	results := makeBoxResults([]time.Duration{
		ms(10), ms(20), ms(30), ms(40), ms(50),
	})
	bs := BuildBoxStats(results)
	if bs.Min != ms(10) {
		t.Errorf("expected min=10ms, got %v", bs.Min)
	}
	if bs.Max != ms(50) {
		t.Errorf("expected max=50ms, got %v", bs.Max)
	}
}

func TestBuildBoxStats_Median(t *testing.T) {
	results := makeBoxResults([]time.Duration{
		ms(10), ms(20), ms(30), ms(40), ms(50),
	})
	bs := BuildBoxStats(results)
	if bs.Median != ms(30) {
		t.Errorf("expected median=30ms, got %v", bs.Median)
	}
}

func TestBuildBoxStats_IQR(t *testing.T) {
	results := makeBoxResults([]time.Duration{
		ms(10), ms(20), ms(30), ms(40), ms(50),
	})
	bs := BuildBoxStats(results)
	if bs.IQR <= 0 {
		t.Errorf("expected positive IQR, got %v", bs.IQR)
	}
	if bs.IQR != bs.Q3-bs.Q1 {
		t.Errorf("IQR mismatch: Q3-Q1=%v IQR=%v", bs.Q3-bs.Q1, bs.IQR)
	}
}

func TestBuildBoxStats_Mean(t *testing.T) {
	results := makeBoxResults([]time.Duration{
		ms(10), ms(20), ms(30), ms(40), ms(50),
	})
	bs := BuildBoxStats(results)
	if bs.Mean != ms(30) {
		t.Errorf("expected mean=30ms, got %v", bs.Mean)
	}
}

func TestWriteBoxStats_Nil(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteBoxStats(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no box stats") {
		t.Errorf("expected fallback message, got: %s", buf.String())
	}
}

func TestWriteBoxStats_ValidOutput(t *testing.T) {
	results := makeBoxResults([]time.Duration{
		ms(10), ms(20), ms(30), ms(40), ms(50),
	})
	bs := BuildBoxStats(results)
	var buf bytes.Buffer
	if err := WriteBoxStats(&buf, bs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, keyword := range []string{"Min", "Q1", "Median", "Mean", "Q3", "Max", "IQR"} {
		if !strings.Contains(out, keyword) {
			t.Errorf("expected output to contain %q", keyword)
		}
	}
}
