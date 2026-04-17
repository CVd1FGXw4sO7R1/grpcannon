package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeRPSResults() []Result {
	now := time.Now()
	return []Result{
		{Start: now, Duration: ms(10), Error: nil},
		{Start: now.Add(100 * time.Millisecond), Duration: ms(10), Error: nil},
		{Start: now.Add(1 * time.Second), Duration: ms(10), Error: errors.New("fail")},
		{Start: now.Add(1*time.Second + 200*time.Millisecond), Duration: ms(10), Error: nil},
		{Start: now.Add(2 * time.Second), Duration: ms(10), Error: nil},
	}
}

func TestCalcRPS_Empty(t *testing.T) {
	points := CalcRPS(nil)
	if points != nil {
		t.Errorf("expected nil, got %v", points)
	}
}

func TestCalcRPS_BucketCount(t *testing.T) {
	results := makeRPSResults()
	points := CalcRPS(results)
	if len(points) != 3 {
		t.Errorf("expected 3 buckets, got %d", len(points))
	}
}

func TestCalcRPS_Values(t *testing.T) {
	results := makeRPSResults()
	points := CalcRPS(results)
	if points[0].RPS != 2 {
		t.Errorf("expected 2 RPS in second 0, got %.0f", points[0].RPS)
	}
	if points[1].RPS != 2 {
		t.Errorf("expected 2 RPS in second 1, got %.0f", points[1].RPS)
	}
	if points[1].Errors != 1 {
		t.Errorf("expected 1 error in second 1, got %d", points[1].Errors)
	}
}

func TestPeakRPS(t *testing.T) {
	points := []RPSPoint{{RPS: 3}, {RPS: 7}, {RPS: 2}}
	if p := peakRPS(points); p != 7 {
		t.Errorf("expected peak 7, got %.0f", p)
	}
}

func TestAvgRPS(t *testing.T) {
	points := []RPSPoint{{RPS: 4}, {RPS: 8}}
	if a := avgRPS(points); a != 6 {
		t.Errorf("expected avg 6, got %.0f", a)
	}
}

func TestAvgRPS_Empty(t *testing.T) {
	if a := avgRPS(nil); a != 0 {
		t.Errorf("expected 0, got %f", a)
	}
}

func TestWriteRPS_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRPS(&buf, makeRPSResults()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Second") {
		t.Error("expected header in output")
	}
	if !strings.Contains(out, "RPS") {
		t.Error("expected RPS column")
	}
}

func TestWriteRPS_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteRPS(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No results") {
		t.Error("expected empty message")
	}
}
