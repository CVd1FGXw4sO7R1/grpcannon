package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeRequestSizeResults() []Result {
	return []Result{
		{Duration: 10 * time.Millisecond, PayloadBytes: 100, Err: nil},
		{Duration: 20 * time.Millisecond, PayloadBytes: 200, Err: nil},
		{Duration: 30 * time.Millisecond, PayloadBytes: 300, Err: nil},
		{Duration: 40 * time.Millisecond, PayloadBytes: 400, Err: errors.New("timeout")},
		{Duration: 50 * time.Millisecond, PayloadBytes: 500, Err: nil},
	}
}

func TestBuildRequestSizeReport_Empty(t *testing.T) {
	got := BuildRequestSizeReport(nil, 4)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestBuildRequestSizeReport_ZeroBuckets(t *testing.T) {
	got := BuildRequestSizeReport(makeRequestSizeResults(), 0)
	if got != nil {
		t.Errorf("expected nil for zero buckets")
	}
}

func TestBuildRequestSizeReport_BucketCount(t *testing.T) {
	res := makeRequestSizeResults()
	got := BuildRequestSizeReport(res, 4)
	if len(got) != 4 {
		t.Errorf("expected 4 buckets, got %d", len(got))
	}
}

func TestBuildRequestSizeReport_TotalCountSums(t *testing.T) {
	res := makeRequestSizeResults()
	buckets := BuildRequestSizeReport(res, 4)
	total := 0
	for _, b := range buckets {
		total += b.Count
	}
	if total != len(res) {
		t.Errorf("expected total count %d, got %d", len(res), total)
	}
}

func TestBuildRequestSizeReport_SuccessCountLteTotal(t *testing.T) {
	res := makeRequestSizeResults()
	buckets := BuildRequestSizeReport(res, 4)
	for _, b := range buckets {
		if b.Successes > b.Count {
			t.Errorf("bucket %s: successes %d > count %d", b.Label, b.Successes, b.Count)
		}
	}
}

func TestBuildRequestSizeReport_AvgPositive(t *testing.T) {
	res := makeRequestSizeResults()
	buckets := BuildRequestSizeReport(res, 2)
	for _, b := range buckets {
		if b.Successes > 0 && b.AvgLatency <= 0 {
			t.Errorf("expected positive avg latency for bucket %s", b.Label)
		}
	}
}

func TestWriteRequestSizeReport_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteRequestSizeReport(&buf, nil)
	if !strings.Contains(buf.String(), "no request-size data") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteRequestSizeReport_ValidOutput(t *testing.T) {
	res := makeRequestSizeResults()
	buckets := BuildRequestSizeReport(res, 3)
	var buf bytes.Buffer
	WriteRequestSizeReport(&buf, buckets)
	out := buf.String()
	if !strings.Contains(out, "Size Range") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "B") {
		t.Errorf("expected byte unit in output, got: %s", out)
	}
}
