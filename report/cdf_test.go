package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeCDFResults(n int, base time.Duration) []Result {
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		results[i] = Result{Duration: base + time.Duration(i)*time.Millisecond}
	}
	return results
}

func TestBuildCDF_Empty(t *testing.T) {
	pts := BuildCDF(nil, 10)
	if pts != nil {
		t.Errorf("expected nil, got %v", pts)
	}
}

func TestBuildCDF_ZeroSteps(t *testing.T) {
	results := makeCDFResults(5, 10*time.Millisecond)
	pts := BuildCDF(results, 0)
	if pts != nil {
		t.Errorf("expected nil for zero steps, got %v", pts)
	}
}

func TestBuildCDF_AllErrors(t *testing.T) {
	results := []Result{
		{Duration: 10 * time.Millisecond, Err: fmt.Errorf("err")},
		{Duration: 20 * time.Millisecond, Err: fmt.Errorf("err")},
	}
	pts := BuildCDF(results, 5)
	if pts != nil {
		t.Errorf("expected nil when all errors, got %v", pts)
	}
}

func TestBuildCDF_StepCount(t *testing.T) {
	results := makeCDFResults(20, 5*time.Millisecond)
	const steps = 8
	pts := BuildCDF(results, steps)
	// steps+1 points expected
	if len(pts) != steps+1 {
		t.Errorf("expected %d points, got %d", steps+1, len(pts))
	}
}

func TestBuildCDF_Monotonic(t *testing.T) {
	results := makeCDFResults(30, 1*time.Millisecond)
	pts := BuildCDF(results, 10)
	for i := 1; i < len(pts); i++ {
		if pts[i].Cumulative < pts[i-1].Cumulative {
			t.Errorf("CDF not monotonic at index %d: %.4f < %.4f", i, pts[i].Cumulative, pts[i-1].Cumulative)
		}
	}
}

func TestBuildCDF_LastPointIsOne(t *testing.T) {
	results := makeCDFResults(10, 2*time.Millisecond)
	pts := BuildCDF(results, 5)
	last := pts[len(pts)-1]
	if last.Cumulative != 1.0 {
		t.Errorf("expected last cumulative=1.0, got %.4f", last.Cumulative)
	}
}

func TestWriteCDF_ValidOutput(t *testing.T) {
	results := makeCDFResults(10, 5*time.Millisecond)
	var buf bytes.Buffer
	WriteCDF(&buf, results, 4)
	out := buf.String()
	if !strings.Contains(out, "Latency") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "Cumulative") {
		t.Errorf("expected cumulative column in output, got: %s", out)
	}
}

func TestWriteCDF_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteCDF(&buf, nil, 5)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' message, got: %s", buf.String())
	}
}
