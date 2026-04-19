package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeAnomalyResults() []Result {
	base := []time.Duration{
		10, 11, 10, 12, 10, 11, 10, 10, 11, 500, // 500ms is the outlier
	}
	results := make([]Result, len(base))
	for i, d := range base {
		results[i] = Result{Duration: d * time.Millisecond}
	}
	return results
}

func TestFindAnomalies_Empty(t *testing.T) {
	if FindAnomalies(nil, 2.0) != nil {
		t.Fatal("expected nil for empty input")
	}
}

func TestFindAnomalies_AllErrors(t *testing.T) {
	results := []Result{
		{Err: errors.New("fail")},
		{Err: errors.New("fail")},
	}
	if FindAnomalies(results, 2.0) != nil {
		t.Fatal("expected nil when all results are errors")
	}
}

func TestFindAnomalies_DetectsOutlier(t *testing.T) {
	results := makeAnomalyResults()
	anomalies := FindAnomalies(results, 2.0)
	if len(anomalies) == 0 {
		t.Fatal("expected at least one anomaly")
	}
	found := false
	for _, a := range anomalies {
		if a.Duration >= 400*time.Millisecond {
			found = true
		}
	}
	if !found {
		t.Fatal("expected the 500ms outlier to be detected")
	}
}

func TestFindAnomalies_UniformNoAnomalies(t *testing.T) {
	var results []Result
	for i := 0; i < 20; i++ {
		results = append(results, Result{Duration: 10 * time.Millisecond})
	}
	if FindAnomalies(results, 2.0) != nil {
		t.Fatal("expected no anomalies for uniform latency")
	}
}

func TestWriteAnomalies_NoneDetected(t *testing.T) {
	var results []Result
	for i := 0; i < 10; i++ {
		results = append(results, Result{Duration: 10 * time.Millisecond})
	}
	var buf bytes.Buffer
	WriteAnomalies(&buf, results, 2.0)
	if !strings.Contains(buf.String(), "none detected") {
		t.Fatal("expected 'none detected'")
	}
}

func TestWriteAnomalies_ValidOutput(t *testing.T) {
	results := makeAnomalyResults()
	var buf bytes.Buffer
	WriteAnomalies(&buf, results, 2.0)
	out := buf.String()
	if !strings.Contains(out, "Anomalies") {
		t.Fatal("expected header")
	}
	if !strings.Contains(out, "Z-Score") {
		t.Fatal("expected Z-Score column")
	}
}
