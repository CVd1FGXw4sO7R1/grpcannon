package report

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func makeFanoutResults(concurrency int, n int, withErr bool) []Result {
	base := time.Now()
	results := make([]Result, n)
	for i := 0; i < n; i++ {
		var err error
		if withErr && i%2 == 0 {
			err = errors.New("rpc error")
		}
		results[i] = Result{
			Duration:    time.Duration(10+i) * time.Millisecond,
			Timestamp:   base.Add(time.Duration(i) * time.Millisecond),
			Error:       err,
			Concurrency: concurrency,
		}
	}
	return results
}

func TestBuildFanout_Empty(t *testing.T) {
	r := BuildFanout(nil, []int{1, 2})
	if len(r.Buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(r.Buckets))
	}
}

func TestBuildFanout_ZeroLevels(t *testing.T) {
	results := makeFanoutResults(4, 10, false)
	r := BuildFanout(results, nil)
	if len(r.Buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(r.Buckets))
	}
}

func TestBuildFanout_BucketCount(t *testing.T) {
	var all []Result
	all = append(all, makeFanoutResults(1, 5, false)...)
	all = append(all, makeFanoutResults(4, 5, false)...)
	all = append(all, makeFanoutResults(8, 5, false)...)

	r := BuildFanout(all, []int{1, 4, 8})
	if len(r.Buckets) != 3 {
		t.Errorf("expected 3 buckets, got %d", len(r.Buckets))
	}
}

func TestBuildFanout_TotalCountSums(t *testing.T) {
	var all []Result
	all = append(all, makeFanoutResults(2, 6, false)...)
	all = append(all, makeFanoutResults(8, 4, false)...)

	r := BuildFanout(all, []int{2, 8})
	total := 0
	for _, b := range r.Buckets {
		total += b.Count
	}
	if total != 10 {
		t.Errorf("expected total 10, got %d", total)
	}
}

func TestBuildFanout_P99Positive(t *testing.T) {
	all := makeFanoutResults(4, 20, false)
	r := BuildFanout(all, []int{4})
	if len(r.Buckets) == 0 {
		t.Fatal("expected at least one bucket")
	}
	if r.Buckets[0].P99Ms <= 0 {
		t.Errorf("expected positive P99, got %f", r.Buckets[0].P99Ms)
	}
}

func TestBuildFanout_OptimalSelected(t *testing.T) {
	var all []Result
	all = append(all, makeFanoutResults(1, 5, false)...)
	all = append(all, makeFanoutResults(8, 20, false)...)

	r := BuildFanout(all, []int{1, 8})
	// concurrency=8 has more results so higher RPS
	if r.Optimal != 8 {
		t.Errorf("expected optimal=8, got %d", r.Optimal)
	}
}

func TestWriteFanout_ValidOutput(t *testing.T) {
	all := makeFanoutResults(4, 10, false)
	r := BuildFanout(all, []int{4})
	var buf bytes.Buffer
	WriteFanout(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Concurrency") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "peak RPS") {
		t.Errorf("expected peak RPS summary, got: %s", out)
	}
}

func TestWriteFanout_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteFanout(&buf, &FanoutReport{})
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected no data message")
	}
}

func TestWriteFanout_NilReport(t *testing.T) {
	var buf bytes.Buffer
	WriteFanout(&buf, nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected no data message for nil report")
	}
}
