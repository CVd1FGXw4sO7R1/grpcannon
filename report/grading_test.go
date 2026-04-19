package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeGradingResults(n int, errEvery int, latency time.Duration) []Result {
	results := make([]Result, n)
	for i := range results {
		r := Result{Duration: latency}
		if errEvery > 0 && (i+1)%errEvery == 0 {
			r.Err = fmt.Errorf("error")
		}
		results[i] = r
	}
	return results
}

func TestAssignGrade_NilReport(t *testing.T) {
	gr := AssignGrade(nil)
	if gr.Grade != GradeF {
		t.Fatalf("expected F, got %s", gr.Grade)
	}
}

func TestAssignGrade_EmptyResults(t *testing.T) {
	r := &Report{}
	gr := AssignGrade(r)
	if gr.Grade != GradeF {
		t.Fatalf("expected F, got %s", gr.Grade)
	}
}

func TestAssignGrade_AllSuccessFastLatency(t *testing.T) {
	results := make([]Result, 100)
	for i := range results {
		results[i] = Result{Duration: 10 * time.Millisecond}
	}
	r := New(results)
	gr := AssignGrade(r)
	if gr.Grade != GradeA {
		t.Fatalf("expected A, got %s (score %d)", gr.Grade, gr.Score)
	}
	if gr.Score < 90 {
		t.Fatalf("expected score >= 90, got %d", gr.Score)
	}
}

func TestAssignGrade_HighErrorRate(t *testing.T) {
	results := make([]Result, 100)
	for i := range results {
		results[i] = Result{Duration: 10 * time.Millisecond, Err: fmt.Errorf("fail")}
	}
	r := New(results)
	gr := AssignGrade(r)
	if gr.Grade != GradeF {
		t.Fatalf("expected F, got %s", gr.Grade)
	}
}

func TestAssignGrade_SlowP99(t *testing.T) {
	results := make([]Result, 100)
	for i := range results {
		latency := 50 * time.Millisecond
		if i >= 99 {
			latency = 2000 * time.Millisecond
		}
		results[i] = Result{Duration: latency}
	}
	r := New(results)
	gr := AssignGrade(r)
	// p99 >= 1000ms → 0 pts for p99; success 100% → 40; p50 fast → 30 = 70 → B
	if gr.Grade != GradeB {
		t.Fatalf("expected B, got %s (score %d)", gr.Grade, gr.Score)
	}
}

func TestAssignGrade_ReasonsNotEmpty(t *testing.T) {
	results := make([]Result, 10)
	for i := range results {
		results[i] = Result{Duration: 30 * time.Millisecond}
	}
	r := New(results)
	gr := AssignGrade(r)
	if len(gr.Reasons) == 0 {
		t.Fatal("expected reasons to be populated")
	}
}

func TestWriteGrade_ContainsGrade(t *testing.T) {
	results := make([]Result, 50)
	for i := range results {
		results[i] = Result{Duration: 20 * time.Millisecond}
	}
	r := New(results)
	var buf bytes.Buffer
	WriteGrade(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Grade:") {
		t.Fatalf("expected 'Grade:' in output, got: %s", out)
	}
}
