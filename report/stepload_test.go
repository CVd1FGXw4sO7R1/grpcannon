package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeStepResults(concurrencies []int, dur time.Duration, withErr bool) []Result {
	results := make([]Result, len(concurrencies))
	for i, c := range concurrencies {
		r := Result{Duration: dur, Concurrency: c}
		if withErr && i%2 == 0 {
			r.Err = errSentinel
		}
		results[i] = r
	}
	return results
}

func TestBuildStepLoad_Empty(t *testing.T) {
	got := BuildStepLoad(nil, 4)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestBuildStepLoad_ZeroSteps(t *testing.T) {
	results := makeStepResults([]int{1, 2, 3}, 10*time.Millisecond, false)
	got := BuildStepLoad(results, 0)
	if got != nil {
		t.Errorf("expected nil for zero steps, got %v", got)
	}
}

func TestBuildStepLoad_BucketCount(t *testing.T) {
	var results []Result
	for i := 1; i <= 20; i++ {
		results = append(results, Result{Duration: 5 * time.Millisecond, Concurrency: i})
	}
	buckets := BuildStepLoad(results, 4)
	if len(buckets) != 4 {
		t.Errorf("expected 4 buckets, got %d", len(buckets))
	}
}

func TestBuildStepLoad_TotalCountSums(t *testing.T) {
	var results []Result
	for i := 1; i <= 10; i++ {
		results = append(results, Result{Duration: 5 * time.Millisecond, Concurrency: i})
	}
	buckets := BuildStepLoad(results, 2)
	total := 0
	for _, b := range buckets {
		total += b.Total
	}
	if total != len(results) {
		t.Errorf("expected total %d, got %d", len(results), total)
	}
}

func TestBuildStepLoad_SuccessCount(t *testing.T) {
	results := []Result{
		{Duration: 10 * time.Millisecond, Concurrency: 1},
		{Duration: 0, Concurrency: 1, Err: errSentinel},
		{Duration: 15 * time.Millisecond, Concurrency: 1},
	}
	buckets := BuildStepLoad(results, 1)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket")
	}
	if buckets[0].Successes != 2 {
		t.Errorf("expected 2 successes, got %d", buckets[0].Successes)
	}
	if buckets[0].Failures != 1 {
		t.Errorf("expected 1 failure, got %d", buckets[0].Failures)
	}
}

func TestBuildStepLoad_P99Positive(t *testing.T) {
	var results []Result
	for i := 0; i < 100; i++ {
		results = append(results, Result{Duration: time.Duration(i+1) * time.Millisecond, Concurrency: 1})
	}
	buckets := BuildStepLoad(results, 1)
	if buckets[0].P99Ms <= 0 {
		t.Errorf("expected positive P99, got %f", buckets[0].P99Ms)
	}
}

func TestWriteStepLoad_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteStepLoad(&buf, nil)
	if !strings.Contains(buf.String(), "no step-load data") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestWriteStepLoad_ValidOutput(t *testing.T) {
	var results []Result
	for i := 1; i <= 8; i++ {
		results = append(results, Result{Duration: time.Duration(i) * time.Millisecond, Concurrency: i})
	}
	buckets := BuildStepLoad(results, 2)
	var buf bytes.Buffer
	WriteStepLoad(&buf, buckets)
	out := buf.String()
	if !strings.Contains(out, "Concurrency") {
		t.Errorf("expected header in output, got %q", out)
	}
	if !strings.Contains(out, "Avg(ms)") {
		t.Errorf("expected Avg(ms) column, got %q", out)
	}
}
