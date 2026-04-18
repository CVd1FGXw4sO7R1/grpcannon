package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeTrendResults(n int, base time.Time, step time.Duration) []Result {
	res := make([]Result, n)
	for i := range res {
		res[i] = Result{
			StartedAt: base.Add(time.Duration(i) * step),
			Duration:  ms(10 + i),
		}
	}
	return res
}

func TestBuildTrendline_Empty(t *testing.T) {
	if BuildTrendline(nil, 5) != nil {
		t.Fatal("expected nil for empty results")
	}
}

func TestBuildTrendline_ZeroBuckets(t *testing.T) {
	base := time.Now()
	res := makeTrendResults(10, base, 100*time.Millisecond)
	if BuildTrendline(res, 0) != nil {
		t.Fatal("expected nil for zero buckets")
	}
}

func TestBuildTrendline_BucketCount(t *testing.T) {
	base := time.Now()
	res := makeTrendResults(20, base, 50*time.Millisecond)
	points := BuildTrendline(res, 4)
	if len(points) != 4 {
		t.Fatalf("expected 4 points, got %d", len(points))
	}
}

func TestBuildTrendline_RPSPositive(t *testing.T) {
	base := time.Now()
	res := makeTrendResults(20, base, 50*time.Millisecond)
	points := BuildTrendline(res, 4)
	for _, p := range points {
		if p.RPS < 0 {
			t.Fatalf("negative RPS: %f", p.RPS)
		}
	}
}

func TestBuildTrendline_WithErrors(t *testing.T) {
	base := time.Now()
	res := makeTrendResults(10, base, 100*time.Millisecond)
	res[0].Err = errTest
	res[1].Err = errTest
	points := BuildTrendline(res, 2)
	if len(points) == 0 {
		t.Fatal("expected points")
	}
}

func TestWriteTrendline_Output(t *testing.T) {
	base := time.Now()
	res := makeTrendResults(10, base, 100*time.Millisecond)
	var buf bytes.Buffer
	WriteTrendline(&buf, res, 2)
	out := buf.String()
	if !strings.Contains(out, "rps") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestWriteTrendline_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteTrendline(&buf, nil, 5)
	if !strings.Contains(buf.String(), "no trendline") {
		t.Errorf("expected empty message")
	}
}
