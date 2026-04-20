package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeHDRResults(durations []time.Duration) []Result {
	results := make([]Result, len(durations))
	for i, d := range durations {
		results[i] = Result{Duration: d}
	}
	return results
}

func TestBuildHDRHistogram_Empty(t *testing.T) {
	buckets := BuildHDRHistogram(nil, 10)
	if len(buckets) != 0 {
		t.Errorf("expected empty, got %d buckets", len(buckets))
	}
}

func TestBuildHDRHistogram_ZeroBuckets(t *testing.T) {
	results := makeHDRResults([]time.Duration{ms(10), ms(20)})
	buckets := BuildHDRHistogram(results, 0)
	if len(buckets) != 0 {
		t.Errorf("expected empty for zero buckets, got %d", len(buckets))
	}
}

func TestBuildHDRHistogram_BucketCount(t *testing.T) {
	durations := []time.Duration{
		ms(1), ms(5), ms(10), ms(50), ms(100), ms(500), ms(1000),
	}
	results := makeHDRResults(durations)
	buckets := BuildHDRHistogram(results, 8)
	if len(buckets) != 8 {
		t.Errorf("expected 8 buckets, got %d", len(buckets))
	}
}

func TestBuildHDRHistogram_CountsSum(t *testing.T) {
	durations := []time.Duration{
		ms(1), ms(5), ms(10), ms(50), ms(100),
	}
	results := makeHDRResults(durations)
	buckets := BuildHDRHistogram(results, 5)
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != len(durations) {
		t.Errorf("expected total count %d, got %d", len(durations), total)
	}
}

func TestBuildHDRHistogram_SkipsErrors(t *testing.T) {
	results := []Result{
		{Duration: ms(10)},
		{Duration: ms(20), Err: errors.New("fail")},
		{Duration: ms(30)},
	}
	buckets := BuildHDRHistogram(results, 4)
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != 2 {
		t.Errorf("expected 2 successful results counted, got %d", total)
	}
}

func TestBuildHDRHistogram_CumPctFinal100(t *testing.T) {
	durations := []time.Duration{ms(10), ms(20), ms(30), ms(40)}
	results := makeHDRResults(durations)
	buckets := BuildHDRHistogram(results, 4)
	last := buckets[len(buckets)-1]
	if last.CumPct < 99.9 {
		t.Errorf("expected final cumulative pct ~100, got %.2f", last.CumPct)
	}
}

func TestWriteHDRHistogram_ValidOutput(t *testing.T) {
	durations := []time.Duration{ms(10), ms(50), ms(200)}
	results := makeHDRResults(durations)
	var buf bytes.Buffer
	WriteHDRHistogram(&buf, results, 4)
	out := buf.String()
	if !strings.Contains(out, "HDR Histogram") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "Lower(ms)") {
		t.Errorf("expected column headers in output")
	}
}

func TestWriteHDRHistogram_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteHDRHistogram(&buf, nil, 5)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' message for empty results")
	}
}
